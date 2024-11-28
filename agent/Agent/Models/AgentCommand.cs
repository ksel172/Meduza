namespace Agent.Models
{
    public class AgentCommand
    {
        public int Id { get; set; } = -1;

        public string Name { get; set; } = string.Empty;

        public DateTime CommandStarted { get; set; } = DateTime.MinValue;

        public DateTime CommandCompleted { get; set; } = DateTime.MinValue;

        public string[] Parameters { get; set; } = [];

        public string Output { get; set; } = string.Empty;
    }
}
