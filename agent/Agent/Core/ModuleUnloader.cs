using Meduza.Agent.ModuleBase;
using System.Reflection;
using System.Runtime.Loader;

namespace Agent.Core
{
    internal static class ModuleUnloader
    {
        internal static void Unload(CustomModuleLoadContext context, IModule module)
        {
            // Cleanup module commands
            if (module.Commands != null)
            {
                foreach (var command in module.Commands)
                {
                    if (command.OutputStream != null)
                    {
                        command.OutputStream.Dispose();
                    }
                }
                module.Commands.Clear();
            }

            // Unload the assembly context
            context.Unload();
            GC.Collect();
            GC.WaitForPendingFinalizers();
        }

        internal static bool IsModuleLoaded(CustomModuleLoadContext context)
        {
            return !context.IsCollectible;
        }
    }
}