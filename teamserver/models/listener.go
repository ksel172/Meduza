package models

import (
	"time"

	"github.com/google/uuid"
)

// Listener type represent String (16 bytes)
type ListenerType = string

// Id represents UUID (16 bytes)
type Id = uuid.UUID

// Status is typed Int (32 bit)
type status = int32

// LogLevel represent different levels of logging verbosity.
type LogLevel = string

// Listener Types (Enum)
const (
	ListenerTypeHTTP    ListenerType = "http"
	ListenerTypeHTTPS   ListenerType = "https"
	ListenerTypeTCP     ListenerType = "tcp"
	ListenerTypeSMB     ListenerType = "smb"
	ListenerTypeForeign ListenerType = "external"
)

// AllowedListenerTypes - Valid listener types
var AllowedListenerTypes = []ListenerType{
	ListenerTypeHTTP,
	ListenerTypeHTTPS,
	ListenerTypeTCP,
	ListenerTypeSMB,
	ListenerTypeForeign,
}

// Listener Statuses (Enum)
const (
	StatusStopped    status = 0
	StatusRunning    status = 1
	StatusPaused     status = 2
	StatusProcessing status = 3
	StatusError      status = 4
)

const (
	Silly LogLevel = "silly" // logs everthing, including verbose information
	Debug LogLevel = "debug" // logs detailed debugging information.
	Info  LogLevel = "info"  // logs general informational messages.
	Error LogLevel = "error" // logs error messages about critical failures.
	Fatal LogLevel = "fatal" // logs critical failures.
	All   LogLevel = "all"   // logs all levels
)

// ValidLogLevel - Valids Log Level
var ValidLogLevel = []LogLevel{
	Silly,
	Debug,
	Info,
	Error,
	Fatal,
	All,
}

// Listener represents a listener configuration , including settings for logging, response rules, and network configuration.
type Listener struct {
	ID            Id            `json:"id"`                    // UUID
	Type          string        `json:"type"`                  // Listener Type (http, tcp, etc.)
	Name          string        `json:"name"`                  // Name
	Status        int           `json:"status"`                // 0 = stopped, 1 = running, 2 = paused, 3 = processing
	Description   string        `json:"description,omitempty"` // description
	Config        Config        `json:"config"`                // Configuration of the Listener
	ResponseRules ResponseRules `json:"response_rules"`        // Response Rules consists headers

	// Logging
	LoggingEnabled bool    `json:"logging_enabled"` // Toggle for enabling Logs
	Logging        Logging `json:"logging"`         // logging structure

	// Time related fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	StoppedAt *time.Time `json:"stopped_at,omitempty"`
}

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
	BlacklistEnabled bool          `json:"blackedlist_enabled"`
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

// ResponseRules defines rules for response headers.
type ResponseRules struct {
	Headers []Header `json:"headers"`
}

// Header represents a key-value pair used in HTTP headers.
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Logging defines the configuration for logging, including the log path and log level.
type Logging struct {
	LogPath  string   `json:"logPath,omitempty"` // Log path
	LogLevel LogLevel `json:"logLevel"`          // Log Level (Example - Silly , debug , info)
}
