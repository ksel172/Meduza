using System.Reflection;
using System.Runtime.Loader;

public class ModuleLoadContext : AssemblyLoadContext
{
    private readonly Dictionary<string, byte[]> dependencyBytes;
    private readonly byte[] moduleBytes;

    public ModuleLoadContext(byte[] moduleBytes, Dictionary<string, byte[]> dependencyBytes)
        : base(isCollectible: true)
    {
        this.moduleBytes = moduleBytes;
        this.dependencyBytes = dependencyBytes;

        this.Resolving += ModuleLoadContext_Resolving;
    }

    private Assembly ModuleLoadContext_Resolving(AssemblyLoadContext context, AssemblyName name)
    {
        Console.WriteLine($"Resolving {name.Name}");
        if (dependencyBytes.TryGetValue(name.Name + ".dll", out byte[] assemblyBytes))
        {
            using var stream = new MemoryStream(assemblyBytes);
            return LoadFromStream(stream);
        }
        return null;
    }

    public Assembly LoadMainModule()
    {
        using var stream = new MemoryStream(moduleBytes);
        return LoadFromStream(stream);
    }
}