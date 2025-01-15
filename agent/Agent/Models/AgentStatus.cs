using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Agent.Models
{
    public enum AgentStatus
    {
        Uninitialized,
        Stage0,
        Stage1,
        Stage2,
        Active,
        Lost,
        Exited,
        Disconnected,
        Hidden
    }
}
