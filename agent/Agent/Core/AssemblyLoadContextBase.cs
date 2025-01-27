using System;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Runtime.Loader;
using Meduza.Agent.ModuleBase;
using Meduza.Agent;
using Agent.Core;

public class AssemblyLoadContextBase : AssemblyLoadContext
{
    private Dictionary<string, Assembly> _loadedAssemblies;
    private ModuleBytesModel? _moduleBytes;

    public AssemblyLoadContextBase() : base(isCollectible: true)
    {
        _loadedAssemblies = new Dictionary<string, Assembly>();
        this.Resolving += OnResolving;
        this.Unloading += OnUnloading;
    }

    private void OnUnloading(AssemblyLoadContext obj)
    {
        _loadedAssemblies.Clear();
        _moduleBytes = null;
    }

    public void RegisterModule(ModuleBytesModel moduleBytes)
    {
        if (moduleBytes.ModuleBytes == null)
            throw new ArgumentNullException(nameof(moduleBytes.ModuleBytes));

        _moduleBytes = moduleBytes;
    }

    public IModule LoadModule(ModuleBytesModel moduleBytes)
    {
        RegisterModule(moduleBytes);

        using (var stream = new MemoryStream(moduleBytes.ModuleBytes))
        {
            var assembly = LoadFromStream(stream);

            // Store the main assembly with its full name
            if (assembly.FullName != null)
            {
                _loadedAssemblies[assembly.FullName] = assembly;
            }

            return LoadModuleFromAssembly(assembly);
        }
    }

    public Assembly LoadFromStream(MemoryStream stream)
    {
        return LoadFromStream(stream, null);
    }

    private IModule LoadModuleFromAssembly(Assembly assembly)
    {
        var moduleType = assembly.GetTypes()
            .FirstOrDefault(type =>
                typeof(IModule).IsAssignableFrom(type) &&
                !type.IsInterface &&
                !type.IsAbstract);

        if (moduleType == null)
        {
            string availableTypes = string.Join(",", assembly.GetTypes().Select(t => t.FullName));
            throw new ApplicationException(
                $"No type implementing IModule found in {assembly.FullName}.\n" +
                $"Available types: {availableTypes}");
        }

        var moduleInstance = Activator.CreateInstance(moduleType);
        if (moduleInstance is not IModule module)
        {
            throw new InvalidCastException(
                $"Failed to cast type '{moduleType.FullName}' to IModule.");
        }

        module.Commands = LoadCommands(assembly).ToList();
        return module;
    }

    private IEnumerable<ICommand> LoadCommands(Assembly assembly)
    {
        var commandTypes = assembly.GetTypes()
            .Where(type =>
                typeof(ICommand).IsAssignableFrom(type) &&
                !type.IsInterface &&
                !type.IsAbstract)
            .ToList();

        if (!commandTypes.Any())
        {
            string availableTypes = string.Join(",", assembly.GetTypes().Select(t => t.FullName));
            throw new ApplicationException(
                $"No types implementing ICommand found in {assembly.FullName}.\n" +
                $"Available types: {availableTypes}");
        }

        return commandTypes.Select(type => (ICommand)Activator.CreateInstance(type)!);
    }

    protected override Assembly? Load(AssemblyName assemblyName)
    {
        // Check if already loaded
        foreach (var assembly in _loadedAssemblies)
        {
            if (assembly.Value.GetName().Name == assemblyName.FullName)
            {
                return assembly.Value;
            }
        }


        if (_moduleBytes?.DependencyBytes != null)
        {
            foreach (var dependencyName in _moduleBytes.DependencyBytes.Keys)
            {
                if (_moduleBytes.DependencyBytes.TryGetValue(dependencyName, out byte[] dependencyBytes))
                {
                    using (var stream = new MemoryStream(dependencyBytes))
                    {
                        var assembly = LoadFromStream(stream);
                        _loadedAssemblies[assembly.FullName] = assembly;
                        return assembly;
                    }
                }
            }
        }
        return null;
    }

    private Assembly? OnResolving(AssemblyLoadContext context, AssemblyName assemblyName)
    {
        return Load(assemblyName);
    }

    public void UnloadContext()
    {
        this.Unload();
    }
}
