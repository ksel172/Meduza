using Agent.Models;
using System;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text.Json;

public class ConfigLoader
{
    public static BaseConfig LoadEmbeddedConfig()
    {
        // Get the current assembly
        var assembly = Assembly.GetExecutingAssembly();

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

        var options = new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true,
            ReadCommentHandling = JsonCommentHandling.Skip
        };

        var baseConfig = JsonSerializer.Deserialize<BaseConfig>(jsonContent, options)
            ?? throw new InvalidOperationException("Failed to deserialize BaseConfig");

        return baseConfig;
    }
}
