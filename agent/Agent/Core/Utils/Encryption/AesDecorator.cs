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
            byte[] derivedKey = DeriveAesKey(key);
            byte[] plaintext = System.Text.Encoding.UTF8.GetBytes(data);

            using var aes = new AesGcm(derivedKey);
            byte[] nonce = new byte[AesGcm.NonceByteSizes.MaxSize];
            using (var rng = RandomNumberGenerator.Create())
            {
                rng.GetBytes(nonce);
            }

            byte[] ciphertext = new byte[plaintext.Length];
            byte[] tag = new byte[AesGcm.TagByteSizes.MaxSize];

            aes.Encrypt(nonce, plaintext, ciphertext, tag);

            // Combine nonce + ciphertext + tag
            byte[] result = new byte[nonce.Length + ciphertext.Length + tag.Length];
            Buffer.BlockCopy(nonce, 0, result, 0, nonce.Length);
            Buffer.BlockCopy(ciphertext, 0, result, nonce.Length, ciphertext.Length);
            Buffer.BlockCopy(tag, 0, result, nonce.Length + ciphertext.Length, tag.Length);

            return Convert.ToBase64String(result);
        }

        private byte[] DeriveAesKey(string key)
        {
            using var sha256 = SHA256.Create();
            return sha256.ComputeHash(System.Text.Encoding.UTF8.GetBytes(key));
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
            byte[] derivedKey = DeriveAesKey(key);
            byte[] combined = Convert.FromBase64String(data);

            int nonceSize = AesGcm.NonceByteSizes.MaxSize;
            int tagSize = AesGcm.TagByteSizes.MaxSize;

            byte[] nonce = new byte[nonceSize];
            byte[] ciphertext = new byte[combined.Length - nonceSize - tagSize];
            byte[] tag = new byte[tagSize];

            Buffer.BlockCopy(combined, 0, nonce, 0, nonceSize);
            Buffer.BlockCopy(combined, nonceSize, ciphertext, 0, ciphertext.Length);
            Buffer.BlockCopy(combined, nonceSize + ciphertext.Length, tag, 0, tagSize);

            byte[] plaintext = new byte[ciphertext.Length];

            using var aes = new AesGcm(derivedKey);
            aes.Decrypt(nonce, ciphertext, tag, plaintext);

            return System.Text.Encoding.UTF8.GetString(plaintext);
        }

        private byte[] DeriveAesKey(string key)
        {
            using var sha256 = SHA256.Create();
            return sha256.ComputeHash(System.Text.Encoding.UTF8.GetBytes(key));
        }
    }
}