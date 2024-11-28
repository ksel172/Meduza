using Agent.Models;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.Linq;
using System.Net.Sockets;
using System.Net;
using System.Text;
using System.Text.RegularExpressions;
using System.Threading.Tasks;

namespace Agent.Services
{
    public class AgentInformationService
    {
        public async Task<AgentInfo> GetAgentInfoAsync()
        {
            return new AgentInfo
            {
                AgentId = await GetUUID(),
                HostName = await GetMachineName(),
                IpAddress = GetLocalIpAddress(),
                UserName = await GetName(),
            };
        }
        public async Task<string> GetUUID()
        {
            string uuidGet = await ExecuteCommand("wmic csproduct get UUID");

            string uuid = uuidGet.Replace("UUID", "");
            if (uuid != null)
            {
                string retUuid = Regex.Replace(uuid, @"\s+", string.Empty);
                return retUuid;
            }
            return "ERROR: Failed to get UUID";
        }
        public async Task<string> GetMOBO()
        {
            string moboIdGet = await ExecuteCommand("wmic baseboard get SerialNumber");

            string moboId = moboIdGet.Replace("SerialNumber", "");
            if (moboId != null)
            {
                string retMobo = Regex.Replace(moboId, @"\s+", string.Empty);
                return retMobo;
            }
            return "ERROR: Failed to get Serial Number";
        }
        public async Task<string> GetName()
        {
            string name = await ExecuteCommand("whoami");
            if (name == null) return "Failed to get username";
            return name;
        }
        public async Task<string> GetMachineName()
        {
            string name = await ExecuteCommand("hostname");
            if (name == null) return "Failed to get hostname";
            return name;
        }

        public async Task<string> GetUsers()
        {
            string users = await ExecuteCommand("net user");
            if (users == null) return "Failed to get machine users";
            return users;
        }
        public async Task<string> GetSystemInfo()
        {
            string sysInfo = await ExecuteCommand("systeminfo");
            if (sysInfo == null) return "Failed to get system info";
            return sysInfo;
        }
        public async Task<string> GetDiskInfo()
        {
            string diskInfo = await ExecuteCommand("wmic logicaldisk get caption,description,drivetype,filesystem");
            if (diskInfo == null) return "Failed to get disk info";
            return diskInfo;
        }
        public async Task<string> GetDiskUsage()
        {
            string usageInfo = await ExecuteCommand("wmic logicaldisk get caption,description,size,freespace\r\n");
            if (usageInfo == null) return "Failed to get disk usage info";
            return usageInfo;
        }
        public async Task<string> GetLoggedUsers()
        {
            string loggedUsers = await ExecuteCommand("quser");
            if (loggedUsers == null) return "Failed to get active users";
            return loggedUsers;
        }

        private string GetLocalIpAddress()
        {
            var host = Dns.GetHostEntry(Dns.GetHostName());
            foreach (var ip in host.AddressList)
            {
                if (ip.AddressFamily is AddressFamily.InterNetwork)
                {
                    return ip.ToString();
                }
            }
            return "169.69.69.69";
        }

        private async Task<string> ExecuteCommand(string command)
        {
            if (string.IsNullOrEmpty(command)) throw new ArgumentNullException(nameof(command));

            Console.WriteLine(AppContext.BaseDirectory);

            var processStartInfo = new ProcessStartInfo
            {
                FileName = "cmd",
                Arguments = $"/c {command}",
                RedirectStandardOutput = true,
                RedirectStandardError = true,
                CreateNoWindow = true,
                UseShellExecute = false,
                WorkingDirectory = AppContext.BaseDirectory
            };

            var process = new Process { StartInfo = processStartInfo };

            process.Start();
            string output = await process.StandardOutput.ReadToEndAsync(); // Capture the output asynchronously
            string errorOutput = await process.StandardError.ReadToEndAsync(); // Capture error output asynchronously
            process.WaitForExit();

            if (!string.IsNullOrEmpty(errorOutput))
            {
                Console.WriteLine("Error Output:");
                Console.WriteLine(errorOutput);
                return errorOutput;
            }
            return output;
        }
    }
}
