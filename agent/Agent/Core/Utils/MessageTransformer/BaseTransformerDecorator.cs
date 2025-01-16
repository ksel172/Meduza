namespace Agent.Core.Utils.MessageTransformer
{
    public abstract class BaseTransformerDecorator : IMessageTransformer
    {
        private readonly IMessageTransformer _wrappedTransformer;

        protected BaseTransformerDecorator(IMessageTransformer wrappedTransformer)
        {
            _wrappedTransformer = wrappedTransformer;
        }

        public virtual string Transform(string input, string key = null)
        {
            return _wrappedTransformer.Transform(input, key);
        }
    }
}
