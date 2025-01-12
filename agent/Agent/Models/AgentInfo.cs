using System.Text.Json.Serialization;

namespace Agent.Models
{
    public class AgentInfo
    {
        [JsonPropertyName("agent_id")]
        public string? AgentId { get; set; }

        [JsonPropertyName("host_name")]
        public string? HostName { get; set; }

        [JsonPropertyName("ip_address")]
        public string? IpAddress { get; set; }

        [JsonPropertyName("username")]
        public string? UserName { get; set; }

        [JsonPropertyName("system_info")]
        public string? SystemInfo { get; set; }

        [JsonPropertyName("os_info")]
        public string? OsInfo { get; set; }
    }
}
