using System;
using System.Collections.Generic;
using System.Linq;
using System.Security.Cryptography;
using System.Text;
using System.Threading.Tasks;

namespace Agent.Core.Utils
{
    internal class Hmac
    {
        private static string GenerateHMAC(byte[] message, byte[] key)
        {
            using (HMACSHA256 hmac = new HMACSHA256(key))
            {
                byte[] hash = hmac.ComputeHash(message);
                return Convert.ToBase64String(hash);
            }
        }

        private static bool VerifyHMAC(byte[] message, string receivedHMAC, byte[] key)
        {
            string computedHMAC = GenerateHMAC(message, key);
            return computedHMAC == receivedHMAC;
        }
    }
}
