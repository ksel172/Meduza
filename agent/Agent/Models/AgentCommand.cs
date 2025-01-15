using System.Text.Json.Serialization;

namespace Agent.Models
{
    public class AgentCommand
    {

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("started")]
        public DateTime CommandStarted { get; set; } = DateTime.MinValue;

        [JsonPropertyName("completed")]
        public DateTime CommandCompleted { get; set; } = DateTime.MinValue;

        [JsonPropertyName("parameters")]
        public string[] Parameters { get; set; } = [];

        [JsonPropertyName("output")]
        public string Output { get; set; } = string.Empty;
    }
}
