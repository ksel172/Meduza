using Agent.Core.Utils.MessageTransformer;
using System.IO.Compression;

namespace Agent.Core.Utils.Compression
{
    public class CompressionDecorator : BaseTransformerDecorator
    {
        public CompressionDecorator(IMessageTransformer wrappedTransformer)
            :base(wrappedTransformer)
        {
        }
        public override string Transform(string input, string key = null)
        {
            string compressedData = Compress(input);
            return base.Transform(compressedData);
        }
        private string Compress(string data)
        {
            byte[] dataBytes = System.Text.Encoding.UTF8.GetBytes(data);
            byte[] compressedBytes;
            using (MemoryStream memoryStream = new MemoryStream())
            {
                using (DeflateStream deflateStream = new DeflateStream(memoryStream, CompressionMode.Compress))
                {
                    deflateStream.Write(dataBytes, 0, dataBytes.Length);
                }
                compressedBytes = memoryStream.ToArray();
            }
            return System.Text.Encoding.UTF8.GetString(compressedBytes);
        }
    }
    public class DecompressionDecorator : BaseTransformerDecorator
    {
        public DecompressionDecorator(IMessageTransformer wrappedTransformer)
            :base(wrappedTransformer)
        {
        }
        public override string Transform(string input, string key = null)
        {
            string compressedData = Decompress(input);
            return base.Transform(compressedData);
        }
        public string Decompress(string compressed)
        {
            using (MemoryStream inputStream = new MemoryStream(compressed.Length))
            {
                inputStream.Write(System.Text.Encoding.UTF8.GetBytes(compressed), 0, compressed.Length);
                inputStream.Seek(0, SeekOrigin.Begin);
                using (MemoryStream outputStream = new MemoryStream())
                {
                    using (DeflateStream deflateStream = new DeflateStream(inputStream, CompressionMode.Decompress))
                    {
                        byte[] buffer = new byte[4096];
                        int bytesRead;
                        while ((bytesRead = deflateStream.Read(buffer, 0, buffer.Length)) != 0)
                        {
                            outputStream.Write(buffer, 0, bytesRead);
                        }
                    }
                    return System.Text.Encoding.UTF8.GetString(outputStream.ToArray());
                }
            }
        }
    }
   
}