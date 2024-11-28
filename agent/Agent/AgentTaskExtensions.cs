﻿using Agent.Models;
using System;
using System.Collections.Concurrent;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Agent
{
    internal static class AgentTaskExtensions
    {
        internal static void QueueRunningStatus(this AgentTask agentTask, ConcurrentQueue<AgentTask> messageQueue, object messageQueueLock)
        {
            agentTask.Status = AgentTaskStatus.Running;
            agentTask.TaskStarted = DateTime.UtcNow;
            lock (messageQueueLock)
            {
                messageQueue.Enqueue(agentTask);
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
}
