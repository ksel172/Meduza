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
        public Guid Id { get; set; }

        public Guid AgentId { get; set; }

        public AgentTaskType Type { get; set; } = AgentTaskType.AgentCommand;

        public AgentTaskStatus Status { get; set; } = AgentTaskStatus.Uninitialized;

        public string? Module { get; set; } = string.Empty;

        public AgentCommand Command { get; set; } = new AgentCommand();

        public DateTime TaskCreated { get; set; } = DateTime.MinValue;

        public DateTime TaskStarted { get; set; } = DateTime.MinValue;

        public DateTime TaskCompleted { get; set; } = DateTime.MinValue;

        public bool IsCancellationTokenSourceSet { get; set; } = false;
    }
}
