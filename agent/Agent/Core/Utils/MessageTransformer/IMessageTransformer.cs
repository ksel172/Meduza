namespace Agent.Core.Utils.MessageTransformer
{
    public interface IMessageTransformer
    {
        string Transform(string input, string key = null);
    }
}
