using Agent.Models.Agent;

namespace Agent.Models
{
    public class C2Request
    {
        public Guid AgentId { get; set; }

        public C2RequestReason Reason { get; set; } = C2RequestReason.Task;

        public AgentStatus AgentStatus { get; set; } = AgentStatus.Uninitialized;

        public string Message { get; set; } = string.Empty;

        public string Hmac { get; set; } = string.Empty;
    }
}
