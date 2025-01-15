using System;
using System.IO;

namespace Agent.ModuleBase
{
    public interface ICommand
    {
        string Name { get; }

        string Description { get; }

        Stream OutputStream { get; set; }

        string Execute(string[] args = null);
    }
}
