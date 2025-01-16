namespace Agent.Core.Utils.MessageTransformer
{
    public class BaseTransformer : IMessageTransformer
    {
        public string Transform(string input, string key = null)
        {
            return input;
        }
    }

}
