using System.Security.Cryptography;
using System.Text;
using Agent.Core.Utils.MessageTransformer;

namespace Agent.Core.Utils.Encryption
{
    public class AesEncryptionDecorator : BaseTransformerDecorator
    {
        private readonly AesComponents aesComponents;

        public AesEncryptionDecorator(IMessageTransformer wrappedTranformer)
            : base(wrappedTranformer)
        {
            aesComponents = new AesComponents();
        }
        public override string Transform(string input, string key)
        {
            if (key == null) return "Key is required";

            string encryptedData = Encrypt(input, key);
            return base.Transform(encryptedData);
        }
        // Returns IV (16 bytes) + EncryptedData byte array
        private string Encrypt(string data, string key)
        {
            var sessionKey = System.Security.Cryptography.Aes.Create();
            sessionKey.Mode = aesComponents.CipherMode;
            sessionKey.Padding = aesComponents.PaddingMode;
            sessionKey.GenerateIV();
            sessionKey.Key = System.Text.Encoding.UTF8.GetBytes(key);

            byte[] encrypted = sessionKey.CreateEncryptor().TransformFinalBlock(System.Text.Encoding.UTF8.GetBytes(data), 0, data.Length);
            byte[] result = new byte[sessionKey.IV.Length + encrypted.Length];
            Buffer.BlockCopy(sessionKey.IV, 0, result, 0, sessionKey.IV.Length);
            Buffer.BlockCopy(encrypted, 0, result, sessionKey.IV.Length, encrypted.Length);
            return System.Text.Encoding.UTF8.GetString(result);
        }
    }
    public class AesDecryptionDecorator : BaseTransformerDecorator
    {
        private readonly AesComponents aesComponents;

        public AesDecryptionDecorator(IMessageTransformer wrappedTranformer)
            : base(wrappedTranformer)
        {
            aesComponents = new AesComponents();
        }
        public override string Transform(string input, string key)
        {
            if (key == null) return "Key is required";

            string encryptedData = Decrypt(input, key);
            return base.Transform(encryptedData);
        }

        // Data should be of format: IV (16 bytes) + EncryptedBytes
        private string Decrypt(string data, string key)
        {
            var sessionKey = System.Security.Cryptography.Aes.Create();
            byte[] iv = new byte[aesComponents.IvLength];
            Buffer.BlockCopy(System.Text.Encoding.UTF8.GetBytes(data), 0, iv, 0, aesComponents.IvLength);
            sessionKey.IV = iv;
            sessionKey.Key = System.Text.Encoding.UTF8.GetBytes(key);
            byte[] encryptedData = new byte[data.Length - aesComponents.IvLength];
            Buffer.BlockCopy(System.Text.Encoding.UTF8.GetBytes(data), aesComponents.IvLength, encryptedData, 0, data.Length - aesComponents.IvLength);
            byte[] decrypted = sessionKey.CreateDecryptor().TransformFinalBlock(encryptedData, 0, encryptedData.Length);

            return System.Text.Encoding.UTF8.GetString(decrypted);
        }

        // Convenience method for decrypting an EncryptedMessagePacket
        //public static byte[] Decrypt(AgentEncryptedMessage encryptedMessage, byte[] key)
        //{
        //    byte[] iv = Convert.FromBase64String(encryptedMessage.IV);
        //    byte[] encrypted = Convert.FromBase64String(encryptedMessage.EncryptedMessage);
        //    byte[] combined = new byte[iv.Length + encrypted.Length];
        //    Buffer.BlockCopy(iv, 0, combined, 0, iv.Length);
        //    Buffer.BlockCopy(encrypted, 0, combined, iv.Length, encrypted.Length);

        //    return Decrypt(combined, key);
        //}
    }
}
