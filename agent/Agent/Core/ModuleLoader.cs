using Agent.ModuleBase;
using System.Reflection;

namespace Agent.Core
{
    internal static class ModuleLoader
    {
        internal static IModule Load(byte[] data)
        {
            return LoadModule(Assembly.Load(data));
        }

        internal static IModule LoadModule(Assembly assembly)
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
                $"Can't find any type which implements IModule in {assembly} from {assembly.Location}.\n" +
                $"Available types: {availableTypes}");
        }

        private static IEnumerable<ICommand> LoadCommands(Assembly assembly)
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
                    $"Can't find any type which implements ICommand in {assembly} from {assembly.Location}.\n" +
                    $"Available types: {availableTypes}");
            }
        }
    }
}
