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
        public Guid Id { get; set; }

        public Dictionary<string, string> Headers { get; set; } = [];

        public List<string> CallbackUrls { get; set; } = [];

        public CallbackRotationType RotationType { get; set; } = CallbackRotationType.Fallback;

        public int RotationRetries { get; set; } = 500;

        public int Sleep { get; set; } = 10;

        // Jitter is a percentage here
        public int Jitter { get; set; } = 20;

        public DateTime? KillDate { get; set; } = DateTime.MaxValue;

        public DateTime WorkingHoursStart { get; set; } = DateTime.MinValue;

        public DateTime WorkingHoursEnd { get; set; } = DateTime.MinValue;
    }
}
