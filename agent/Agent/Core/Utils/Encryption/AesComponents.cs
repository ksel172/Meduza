using System.Security.Cryptography;

namespace Agent.Core.Utils.Encryption
{
    public class AesComponents
    {
        public CipherMode CipherMode { get; set; } = CipherMode.CBC;
        public PaddingMode PaddingMode { get; set; } = PaddingMode.PKCS7;
        public int IvLength { get; set; } = 16;
    }
}
