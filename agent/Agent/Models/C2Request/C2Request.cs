using System.Text.Json.Serialization;

namespace Agent.Models.C2Request
{
    public class RegisterC2Request
    {
        [JsonPropertyName("agent_id")]
        public string AgentId { get; set; } = string.Empty;

        [JsonPropertyName("config_id")]
        public string ConfigId { get; set; } = string.Empty;

        [JsonPropertyName("agent_status")]
        public string AgentStatus { get; set; } = "uninitialized";

        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;

        // public string Hmac { get; set; } = string.Empty;
    }

    public class CheckinC2Request
    {
        [JsonPropertyName("agent_id")]
        public string AgentId { get; set; } = string.Empty;

        [JsonPropertyName("reason")]
        public string Reason { get; set; } = string.Empty;

        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;
    }
}
