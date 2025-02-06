using System;
using System.Collections.Generic;
using System.Text.Json.Serialization;

namespace Agent.Models
{

    public class BaseConfig
    {
        [JsonPropertyName("agent_id")]
        public string AgentId { get; set; } = string.Empty;

        [JsonPropertyName("config_id")]
        public string AgentConfigId { get; set; } = string.Empty;

        [JsonPropertyName("listener_id")]
        public string ListenerId { get; set; } = string.Empty;

        [JsonPropertyName("token")]
        public string Token { get; set; } = string.Empty;

        [JsonPropertyName("sleep")]
        public int Sleep { get; set; } = 5;

        [JsonPropertyName("jitter")]
        public int Jitter { get; set; } = 3;

        [JsonPropertyName("kill_date")]
        public DateTime? KillDate { get; set; }

        [JsonPropertyName("working_hours_start")]
        public int WorkingHoursStart { get; set; }

        [JsonPropertyName("working_hours_end")]
        public int WorkingHoursEnd { get; set; }

        [JsonPropertyName("config")]
        public CommunicationConfig Config { get; set; }

        public class CommunicationConfig
        {
            [JsonPropertyName("hosts")]
            public List<string> Hosts { get; set; }

            [JsonPropertyName("headers")]
            public List<Header> Headers { get; set; }

            [JsonPropertyName("host_rotation")]
            public string HostRotation { get; set; }

            [JsonPropertyName("proxy_settings")]
            public ProxySettings ProxySettings { get; set; }
        }

        public class Header
        {
            [JsonPropertyName("key")]
            public string Key { get; set; }

            [JsonPropertyName("value")]
            public string Value { get; set; }
        }

        public class ProxySettings
        {
            [JsonPropertyName("enabled")]
            public bool Enabled { get; set; } = false;

            [JsonPropertyName("password")]
            public string Password { get; set; }

            [JsonPropertyName("port")]
            public string Port { get; set; }

            [JsonPropertyName("type")]
            public string Type { get; set; }

            [JsonPropertyName("username")]
            public string Username { get; set; }
        }
    }
}
