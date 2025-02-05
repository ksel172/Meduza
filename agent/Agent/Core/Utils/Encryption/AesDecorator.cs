using System;
using System.Security.Cryptography;
using System.Text;
using Agent.Core.Utils.MessageTransformer;

namespace Agent.Core.Utils.Encryption
{
    public class AesEncryptionDecorator : BaseTransformerDecorator
    {
        public AesEncryptionDecorator(IMessageTransformer wrappedTransformer)
            : base(wrappedTransformer) { }

        public override string Transform(string input, string key)
        {
            if (string.IsNullOrEmpty(key)) return "Key is required";

            string encryptedData = Encrypt(input, key);
            return base.Transform(encryptedData);
        }

        private string Encrypt(string data, string key)
        {
            using var aes = Aes.Create();
            aes.Key = DeriveAesKey(key);
            aes.GenerateIV(); // Generates a new IV
            aes.Mode = CipherMode.CBC;
            aes.Padding = PaddingMode.PKCS7;

            using var encryptor = aes.CreateEncryptor();
            byte[] plainBytes = System.Text.Encoding.UTF8.GetBytes(data);
            byte[] encryptedBytes = encryptor.TransformFinalBlock(plainBytes, 0, plainBytes.Length);

            // Prepend IV to encrypted data and encode in Base64
            byte[] result = new byte[aes.IV.Length + encryptedBytes.Length];
            Buffer.BlockCopy(aes.IV, 0, result, 0, aes.IV.Length);
            Buffer.BlockCopy(encryptedBytes, 0, result, aes.IV.Length, encryptedBytes.Length);

            return Convert.ToBase64String(result);
        }

        private byte[] DeriveAesKey(string key)
        {
            using var sha256 = SHA256.Create();
            return sha256.ComputeHash(System.Text.Encoding.UTF8.GetBytes(key)).AsSpan(0, 32).ToArray(); // 256-bit AES key
        }
    }

    public class AesDecryptionDecorator : BaseTransformerDecorator
    {
        public AesDecryptionDecorator(IMessageTransformer wrappedTransformer)
            : base(wrappedTransformer) { }

        public override string Transform(string input, string key)
        {
            if (string.IsNullOrEmpty(key)) return "Key is required";

            string decryptedData = Decrypt(input, key);
            return base.Transform(decryptedData);
        }

        private string Decrypt(string data, string key)
        {
            using var aes = Aes.Create();
            aes.Key = DeriveAesKey(key);
            aes.Mode = CipherMode.CBC;
            aes.Padding = PaddingMode.PKCS7;

            byte[] inputBytes = Convert.FromBase64String(data);

            // Extract IV
            byte[] iv = new byte[aes.BlockSize / 8];
            byte[] encryptedBytes = new byte[inputBytes.Length - iv.Length];

            Buffer.BlockCopy(inputBytes, 0, iv, 0, iv.Length);
            Buffer.BlockCopy(inputBytes, iv.Length, encryptedBytes, 0, encryptedBytes.Length);

            aes.IV = iv;

            using var decryptor = aes.CreateDecryptor();
            byte[] decryptedBytes = decryptor.TransformFinalBlock(encryptedBytes, 0, encryptedBytes.Length);

            return System.Text.Encoding.UTF8.GetString(decryptedBytes);
        }

        private byte[] DeriveAesKey(string key)
        {
            using var sha256 = SHA256.Create();
            return sha256.ComputeHash(System.Text.Encoding.UTF8.GetBytes(key)).AsSpan(0, 32).ToArray(); // 256-bit AES key
        }
    }
}
