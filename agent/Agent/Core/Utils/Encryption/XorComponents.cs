using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Agent.Core.Utils.Encryption
{
    public class XorComponents
    {
        public string Delimiter { get; set; } = "jRlOPs";

        public byte[] XorCrypt(byte[] data, byte[] keyBytes)
        {
            for (int i = 0, j = 0; i < data.Length; i++, j++)
            {
                if (j == keyBytes.Length) j = 0;
                data[i] = (byte)(data[i] ^ keyBytes[j]);
            }
            return data;
        }

        // CHANGED: Changed this to return a byte[] since the rest of the code is written to use them
        public byte[] GenerateRandomKey()
        {
            var random = new Random();
            var length = random.Next(24, 33);

            var randomKey = new byte[length];

            for (int i = 0; i < length; i++)
            {
                char randomChar;
                do
                {
                    randomChar = (char)random.Next(33, 127);
                } while (randomChar == '\\' || randomChar == '/' || randomChar == ' ');

                randomKey[i] = (byte)randomChar;
            }

            return randomKey;
        }
        public List<int> Search(byte[] src, byte[] pattern)
        {
            var results = new List<int>();
            int maxFirstCharSlot = src.Length - pattern.Length + 1;
            for (int i = 0; i < maxFirstCharSlot; i++)
            {
                if (src[i] != pattern[0])
                    continue;

                for (int j = pattern.Length - 1; j >= 1; j--)
                {
                    if (src[i + j] != pattern[j]) break;
                    if (j == 1) results.Add(i);
                }
            }
            return results;
        }
    }
}
