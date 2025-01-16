using System.Net.Http.Headers;
using System.Text;
using Agent.Core.Utils.MessageTransformer;

namespace Agent.Core.Utils.Encryption
{
    public class XorEncryptionDecorator : BaseTransformerDecorator
    {
        private readonly XorComponents xorComponents;
        public XorEncryptionDecorator(IMessageTransformer wrappedTranformer)
            : base(wrappedTranformer)
        {
            xorComponents = new XorComponents();
        }
        public override string Transform(string input, string key = null)
        {
            string encryptedData = Encrypt(input);
            return base.Transform(encryptedData);
        }
        private string Encrypt(string data)
        {
            var key = xorComponents.GenerateRandomKey();
            var bytes = System.Text.Encoding.UTF8.GetBytes(data);
            var crypted = xorComponents.XorCrypt(bytes, key);
            var magic = System.Text.Encoding.UTF8.GetBytes(xorComponents.Delimiter);
            var combinedKey = new byte[magic.Length + key.Length + magic.Length];
            Buffer.BlockCopy(magic, 0, combinedKey, 0, magic.Length);                           // Prepend magic
            Buffer.BlockCopy(key, 0, combinedKey, magic.Length, key.Length);                    // Copy key
            Buffer.BlockCopy(magic, 0, combinedKey, magic.Length + key.Length, magic.Length);   // Append magic

            var random = new Random();
            var startPosition = random.Next(1, crypted.Length);
            var endPosition = startPosition + combinedKey.Length;
            var combined = new byte[combinedKey.Length + crypted.Length];
            Buffer.BlockCopy(crypted, 0, combined, 0, startPosition);
            Buffer.BlockCopy(combinedKey, 0, combined, startPosition, combinedKey.Length);
            Buffer.BlockCopy(crypted, startPosition, combined, endPosition, crypted.Length - startPosition);

            return System.Text.Encoding.UTF8.GetString(combined);
        }
    }
    public class XorDecryptionDecorator : BaseTransformerDecorator
    {
        private readonly XorComponents xorComponents;
        public XorDecryptionDecorator(IMessageTransformer wrappedTransformer)
            : base(wrappedTransformer)
        {
            xorComponents = new XorComponents();
        }
        public override string Transform(string input, string key = null)
        {
            string decryptedData = Decrypt(input);
            return base.Transform(decryptedData);
        }
        private string Decrypt(string data)
        {
            var decoded = System.Text.Encoding.UTF8.GetBytes(data);
            var magic = System.Text.Encoding.UTF8.GetBytes(xorComponents.Delimiter);
            var magicPositions = xorComponents.Search(decoded, magic);
            if (magicPositions.Count == 0) return string.Empty;
            var keyLength = (magicPositions[1] - magicPositions[0]) - magic.Length;
            var key = new byte[keyLength];

            Buffer.BlockCopy(decoded, magicPositions[0] + magic.Length, key, 0, keyLength);

            var totalDiscardLength = magic.Length + keyLength + magic.Length;
            var encrypted = new byte[decoded.Length - totalDiscardLength];

            Buffer.BlockCopy(decoded, 0, encrypted, 0, magicPositions[0]);
            Buffer.BlockCopy(decoded, magicPositions[0] + totalDiscardLength, encrypted, magicPositions[0], encrypted.Length - magicPositions[0]);

            var decrypted = xorComponents.XorCrypt(encrypted, key);

            return System.Text.Encoding.UTF8.GetString(decrypted);
        }
    }
}