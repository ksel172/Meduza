using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Agent.Services
{
    internal class OperatingSystemInformationService
    {
        private readonly OperatingSystem operatingSystem;

        internal OperatingSystemInformationService()
        {
            operatingSystem = Environment.OSVersion;
        }

        internal string GetOsName()
        {
            try
            {
                return operatingSystem.VersionString.Replace("Microsoft Windows NT ", "Windows ")[0..^8];
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error getting Operating System name: {ex.Message}");
            }

            return string.Empty;
        }
    }
}
