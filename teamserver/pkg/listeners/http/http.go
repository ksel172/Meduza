package http

// Config represents the configuration settings for the service, including network, security, and proxy settings.
type Config struct {
	WorkingHours     string        `json:"working_hours"`
	Hosts            []string      `json:"hosts"`
	HostBind         string        `json:"host_bind"`
	HostRotation     string        `json:"host_rotation"`
	PortBind         string        `json:"port_bind"`
	PortConn         string        `json:"port_conn"`
	Secure           bool          `json:"secure"`
	HostHeader       string        `json:"host_header"`
	Headers          []Header      `json:"headers"`
	Uris             []string      `json:"uris"`
	Certificate      Certificate   `json:"certificate"`
	WhitelistEnabled bool          `json:"whitelist_enabled"`
	Whitelist        []string      `json:"whitelist"`
	BlacklistEnabled bool          `json:"blacklist_enabled"`
	ProxySettings    ProxySettings `json:"proxy_settings"`
}

// Certificate holds the paths to the SSL certificate and its corresponding private key.
type Certificate struct {
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}

// ProxySettings represents the configuration for a proxy server.
type ProxySettings struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Header represents a key-value pair used in HTTP headers.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
