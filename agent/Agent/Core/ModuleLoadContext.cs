
using System.Reflection;
using System.Runtime.Loader;

public class ModuleLoadContext : AssemblyLoadContext
{
    public byte[]? AssemblyBytes { get; set; }
    public Dictionary<string, byte[]>? DependencyBytes { get; set; }

    public ModuleLoadContext() : base(isCollectible: true)
    {
        Resolving += ResolveAssembly;
    }

    private Assembly? ResolveAssembly(AssemblyLoadContext context, AssemblyName assemblyName)
    {
        if (DependencyBytes != null && assemblyName.Name != null)
        {
            if (DependencyBytes.TryGetValue(assemblyName.Name, out var dependencyBytes))
            {
                return LoadFromStream(new MemoryStream(dependencyBytes));
            }
        }
        return null;
    }

    protected override Assembly? Load(AssemblyName assemblyName)
    {
        return null; // Let the custom resolver handle assembly loading
    }

    public Assembly LoadFromMemory()
    {
        if (AssemblyBytes == null)
        {
            throw new InvalidOperationException("No module assembly provided.");
        }

        LoadDependencies();
        return LoadFromStream(new MemoryStream(AssemblyBytes));
    }

    public void LoadDependencies()
    {
        if (DependencyBytes != null)
        {
            foreach (var dependency in DependencyBytes)
            {
                // Load each dependency into the context from memory
                var dependencyStream = new MemoryStream(dependency.Value);
                LoadFromStream(dependencyStream);
            }
        }
    }
}
