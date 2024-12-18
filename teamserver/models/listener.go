package models

import "time"

// Listener Types (Enum)
const (
	ListenerTypeHTTP    = "http"
	ListenerTypeTCP     = "tcp"
	ListenerTypeSMB     = "smb"
	ListenerTypeForeign = "foreign"
)

var AllowedListenerTypes = []string{
	ListenerTypeHTTP,
	ListenerTypeTCP,
	ListenerTypeSMB,
	ListenerTypeForeign,
}

const (
	StatusStopped    = 0
	StatusRunning    = 1
	StatusPaused     = 2
	StatusProcessing = 3
	StatusError      = 4
)

type ListenerRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Config      any    `json:"config"`
}

type Listener struct {
	ID          string     `json:"id"`
	Type        string     `json:"type"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Status      int        `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	StartedAt   *time.Time `json:"started_at"`
	StoppedAt   *time.Time `json:"stopped_at"`
	Config      any        `json:"config"`
}

type HTTPConfig struct {
	KillDate         int64            `json:"kill_date"`
	WorkingHours     string           `json:"working_hours"`
	Hosts            []string         `json:"hosts"`
	HostBind         string           `json:"host_bind"`
	HostRotation     string           `json:"host_rotation"`
	PortBind         string           `json:"port_bind"`
	PortConn         string           `json:"port_conn"`
	Secure           bool             `json:"secure"`
	HostHeader       string           `json:"host_header"`
	Headers          []Header         `json:"headers"`
	Uris             []string         `json:"uris"`
	Certificate      Certificate      `json:"certificate"`
	WhitelistEnabled bool             `json:"whitelist_enabled"`
	Whitelist        []string         `json:"whitelist,omitempty"`
	BlacklistEnabled bool             `json:"blacklist_enabled"`
	Blacklist        []string         `json:"blacklist,omitempty"`
	ProxySettings    ProxyConfig      `json:"proxy_settings"`
	ResponseRules    ResponseSettings `json:"response_rules"`
}

type TCPConfig struct {
	PortBind   string `json:"port_bind"`
	HostBind   string `json:"host_bind"`
	BufferSize int    `json:"buffer_size"`
	Timeout    int    `json:"timeout"`
}

type SMBConfig struct {
	PipeName     string `json:"pipe_name"`
	MaxInstances int    `json:"max_instances"`
	KillDate     int64  `json:"kill_date"`
}

type ForeignConfig struct {
	Endpoint       string         `json:"endpoint"`
	Authentication Authentication `json:"authentication"`
}

type Authentication struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Certificate struct {
	CertPath string `json:"cert_path"`
	KeyPath  string `json:"key_path"`
}

type ProxyConfig struct {
	Enabled  bool   `json:"enabled"`
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseSettings struct {
	Headers []Header `json:"headers"`
}
