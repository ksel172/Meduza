using Agent.Models;
using Agent.Services;

AgentInformationService agentInformationService = new AgentInformationService();
AgentInfo agentInformation = await agentInformationService.GetAgentInfoAsync();
try
{
    var baseConfig = ConfigLoader.LoadEmbeddedConfig();
    
    if (agentInformation is not null)
    {
        baseConfig.AgentId = agentInformation.AgentId;
    }

    while (true)
    {

        await Task.Delay(baseConfig.Sleep * 1000);
    }
}
catch (Exception ex)
{
    Console.WriteLine($"Agent initialization failed: {ex.Message}");
}