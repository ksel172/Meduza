using Agent.Models;
using System.Collections.Concurrent;

internal static class AgentTaskExtensions
{
    private static readonly SemaphoreSlim semaphore = new SemaphoreSlim(1, 1);

    internal static void QueueRunningStatus(this AgentTask agentTask, ConcurrentQueue<AgentTask> messageQueue, object messageQueueLock)
    {
        Console.WriteLine($"Attempting to lock for task {agentTask.Id}");
        semaphore.Wait();
        try
        {
            agentTask.Status = AgentTaskStatus.Running;
            agentTask.TaskStarted = DateTime.UtcNow;
            Console.WriteLine($"Task {agentTask.Id} status set to running");
            lock (messageQueueLock)
            {
                messageQueue.Enqueue(agentTask);
                Console.WriteLine($"Task {agentTask.Id} enqueued");
            }
        }
        finally
        {
            semaphore.Release();
            Console.WriteLine($"Lock released for task {agentTask.Id}");
        }
    }

    internal static void QueueQueuedStatus(this AgentTask agentTask, ConcurrentQueue<AgentTask> messageQueue, object messageQueueLock)
    {
        agentTask.Status = AgentTaskStatus.Queued;
        lock (messageQueueLock)
        {
            messageQueue.Enqueue(agentTask);
        }
    }

    internal static void QueueCompletedStatus(this AgentTask agentTask, ConcurrentQueue<AgentTask> messageQueue, object messageQueueLock)
    {
        agentTask.Status = AgentTaskStatus.Complete;
        agentTask.TaskCompleted = DateTime.UtcNow;
        lock (messageQueueLock)
        {
            messageQueue.Enqueue(agentTask);
        }
    }
}
