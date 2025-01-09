namespace Agent.Models.C2Request
{
    public class C2Request
    {
        public string AgentId { get; set; } = string.Empty;

        public C2RequestReason Reason { get; set; } = C2RequestReason.Task;

        public AgentStatus AgentStatus { get; set; } = AgentStatus.Uninitialized;

        public string Message { get; set; } = string.Empty;

        public string Hmac { get; set; } = string.Empty;
    }
}
