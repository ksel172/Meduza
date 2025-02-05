using System.Text.Json.Serialization;

namespace Agent.Models.C2Request
{
    public class C2Request
    {
        [JsonPropertyName("reason")]
        public C2RequestReason Reason { get; set; }

        [JsonPropertyName("agent_id")]
        public string AgentId { get; set; } = string.Empty;

        [JsonPropertyName("config_id")]
        public string ConfigId { get; set; } = string.Empty;

        [JsonPropertyName("agent_status")]
        public AgentStatus AgentStatus { get; set; } = AgentStatus.Uninitialized;

        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;

        // public string Hmac { get; set; } = string.Empty;
    }

    public class AuthenticationResponse
    {

        [JsonPropertyName("public_key")]
        public string PublicKey { get; set; } = string.Empty;

        [JsonPropertyName("session_token")]
        public string SessionToken { get; set; } = string.Empty;
    }

    public class EncryptedRequest
    {
        [JsonPropertyName("message")]
        public string Message { get; set; } = string.Empty;

    }
}