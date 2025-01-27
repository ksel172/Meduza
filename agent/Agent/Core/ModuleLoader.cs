using Agent.Core;
using Meduza.Agent.ModuleBase;
using System.Reflection;

internal static class ModuleLoader
{
    internal static IModule Load(ModuleBytesModel moduleData)
    {
        var loadContext = new ModuleLoadContext(
            moduleData.ModuleBytes,
            moduleData.DependencyBytes ?? new Dictionary<string, byte[]>());

        Assembly assembly = loadContext.LoadMainModule();
        return LoadModule(assembly);
    }

    public static IModule LoadModule(Assembly assembly)
    {
        foreach (Type type in assembly.GetTypes())
        {
            if (typeof(IModule).IsAssignableFrom(type))
            {
                if (Activator.CreateInstance(type) is IModule result)
                {
                    result.Commands = LoadCommands(assembly).ToList();
                    return result;
                }
            }
        }
        string availableTypes = string.Join(",", assembly.GetTypes().Select(t => t.FullName));
        throw new ApplicationException(
            $"Can't find any type which implements IModule in {assembly}.\n" +
            $"Available types: {availableTypes}");
    }

    public static IEnumerable<ICommand> LoadCommands(Assembly assembly)
    {
        var count = 0;
        foreach (Type type in assembly.GetTypes())
        {
            if (typeof(ICommand).IsAssignableFrom(type))
            {
                if (Activator.CreateInstance(type) is ICommand result)
                {
                    count++;
                    yield return result;
                }
            }
        }
        if (count == 0)
        {
            string availableTypes = string.Join(",", assembly.GetTypes().Select(t => t.FullName));
            throw new ApplicationException(
                $"Can't find any type which implements ICommand in {assembly}.\n" +
                $"Available types: {availableTypes}");
        }
    }
}