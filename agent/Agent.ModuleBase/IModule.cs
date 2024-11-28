using System.Collections.Generic;
using System.Dynamic;

namespace Agent.ModuleBase
{
    public interface IModule
    {
        string Name { get; }

        string Author { get; }

        string Description { get; }

        string FileName { get; }

        string Usage { get; }

        List<ICommand> Commands { get; set; }
    }
}
