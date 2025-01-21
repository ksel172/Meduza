using Agent.Core.Utils.MessageTransformer;

namespace Agent.Core.Utils.Encoding
{
    public class UrlSafeBase64EncodingDecorator : BaseTransformerDecorator
    {
        // https://stackoverflow.com/questions/26353710/how-to-achieve-base64-url-safe-encoding-in-c
        public UrlSafeBase64EncodingDecorator(IMessageTransformer wrappedTransformer)
            : base(wrappedTransformer)
        {
        }

        public override string Transform(string input, string key = null)
        {
            string encodedData = Encode(System.Text.Encoding.UTF8.GetBytes(input));
            return base.Transform(encodedData);
        }
        private string Encode(byte[] data)
        {
            if (data is null) return string.Empty;

            return Convert.ToBase64String(data)
                .TrimEnd('=')
                .Replace('+', '-')
                .Replace('/', '_');
        }
    }
    public class UrlSafeBase64DecodingDecorator : BaseTransformerDecorator
    {
        public UrlSafeBase64DecodingDecorator(IMessageTransformer wrappedTransformer)
           : base(wrappedTransformer)
        {
        }
        public override string Transform(string input, string key = default)
        {
            string decodedData = Decode(input);
            return base.Transform(decodedData);
        }
        private string Decode(string data)
        {
            if (string.IsNullOrWhiteSpace(data)) return string.Empty;

            var converted = data.Replace('-', '+')
                                .Replace('_', '/');

            switch (converted.Length % 4)
            {
                case 2:
                    converted += "==";
                    break;
                case 3:
                    converted += "=";
                    break;
            }

            return System.Text.Encoding.UTF8.GetString(Convert.FromBase64String(converted));
        }
    }
}
