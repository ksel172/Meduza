using System.Text.Json.Serialization;

namespace Agent.Models
{
    public enum AgentTaskType
    {
        LoadAssembly,
        UnloadAssembly,
        AgentCommand,
        ShellCommand,
        HelpCommand,
        SetDelay,
        SetJitter,
        GetTasks,
        KillTasks,
        Exit,
        Unknown
    }

    public enum AgentTaskStatus
    {
        Uninitialized,
        Queued,
        Sent,
        Running,
        Complete,
        Failed,
        Aborted
    }

    public class AgentTask
    {
        [JsonPropertyName("task_id")]
        public string Id { get; set; } = string.Empty;

        [JsonPropertyName("agent_id")]
        public string AgentId { get; set; } = string.Empty;

        [JsonPropertyName("type")]
        public AgentTaskType Type { get; set; } = AgentTaskType.AgentCommand;

        [JsonPropertyName("status")]
        public AgentTaskStatus Status { get; set; } = AgentTaskStatus.Uninitialized;

        [JsonPropertyName("module")]
        public string? Module { get; set; } = string.Empty;

        [JsonPropertyName("command")]
        public AgentCommand Command { get; set; } = new AgentCommand();

        [JsonPropertyName("created")]
        public DateTime TaskCreated { get; set; } = DateTime.MinValue;

        [JsonPropertyName("started")]
        public DateTime TaskStarted { get; set; } = DateTime.MinValue;

        [JsonPropertyName("finished")]
        public DateTime TaskCompleted { get; set; } = DateTime.MinValue;

        public bool IsCancellationTokenSourceSet { get; set; } = false;
    }
}
