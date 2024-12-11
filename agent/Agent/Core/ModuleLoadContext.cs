using System.Reflection;
using System.Runtime.Loader;

namespace Agent.Core
{
    public class ModuleLoadContext : AssemblyLoadContext
    {
        public Stream? AssemblyBytes { get; set; }

        private AssemblyDependencyResolver resolver;


        protected override Assembly? Load(AssemblyName assemblyName)
        {
            if (AssemblyBytes is not null)
            {
                return LoadFromStream(AssemblyBytes);
            }

            return base.Load(assemblyName);
        }
    }
}
