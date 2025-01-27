using Meduza.Agent;
using Meduza.Agent.ModuleBase;
using System.Reflection;
using System.Runtime.Loader;

namespace Agent.Core
{
    internal static class ModuleLoader
    {
        internal static IModule Load(ModuleBytesModel moduleBytesModel)
        {
            if (moduleBytesModel.ModuleBytes == null)
            {
                throw new ArgumentNullException(nameof(moduleBytesModel.ModuleBytes));
            }

            var loader = new AssemblyLoadContextBase();
            return loader.LoadModule(moduleBytesModel);
        }
    }
}