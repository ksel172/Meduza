using Agent.Models;
using System.Reflection;
using System.Text.Json;

public class ConfigLoader
{
    public static AgentConfig LoadEmbeddedConfig()
    {
        // Get the current assembly
        var assembly = Assembly.GetExecutingAssembly();

        // Find the resource name (typically namespace.filename)
        string[] resourceNames = assembly.GetManifestResourceNames();
        string configResourceName = resourceNames.FirstOrDefault(r =>
            r.EndsWith("baseconf.json", StringComparison.OrdinalIgnoreCase))
            ?? throw new FileNotFoundException("Embedded configuration file not found");

        // Read the embedded resource
        using var stream = assembly.GetManifestResourceStream(configResourceName);
        if (stream == null)
            throw new InvalidOperationException("Could not load embedded configuration stream");

        using var reader = new StreamReader(stream);
        string jsonContent = reader.ReadToEnd();

        // Create custom deserialization options
        var options = new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true,
            ReadCommentHandling = JsonCommentHandling.Skip
        };

        // Deserialize the nested configuration
        var rootConfig = JsonSerializer.Deserialize<RootConfig>(jsonContent, options)
            ?? throw new InvalidOperationException("Failed to deserialize configuration");

        return rootConfig.AgentConfig;
    }
    private class RootConfig
    {
        public AgentConfig AgentConfig { get; set; } = new AgentConfig();
    }
}