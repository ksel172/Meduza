using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace Agent.Core
{
    public class ModuleBytesModel
    {
        [JsonPropertyName("module_bytes")]
        public byte[] ModuleBytes { get; set; }

        [JsonPropertyName("dependency_bytes")]
        public Dictionary<string, byte[]> DependencyBytes { get; set; }
    }
}
