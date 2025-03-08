package models

import (
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
)

const (
	// URL parameter constants
	ParamListenerID string = "listener_id"
)

// Listener Statuses (Enum)
type status uint8

const (
	StatusStopped status = iota
	StatusRunning
	StatusPaused
	StatusProcessing
	StatusError
)

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

// Log levels
type LogLevel string

const (
	Silly LogLevel = "silly" // logs everthing, including verbose information
	Debug LogLevel = "debug" // logs detailed debugging information.
	Info  LogLevel = "info"  // logs general informational messages.
	Error LogLevel = "error" // logs error messages about critical failures.
	Fatal LogLevel = "fatal" // logs critical failures.
	All   LogLevel = "all"   // logs all levels
)

// ValidLogLevel contains all valid logging levels for runtime validation.
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
	ID          uuid.UUID `json:"id"`                    // UUID
	Type        string    `json:"type"`                  // Listener Type (http, tcp, etc.)
	Name        string    `json:"name"`                  // Name
	Status      int       `json:"status"`                // 0 = stopped, 1 = running, 2 = paused, 3 = processing
	Description string    `json:"description,omitempty"` // description
	Config      any       `json:"config"`                // Configuration of the Listener

	// Logging
	LoggingEnabled bool    `json:"logging_enabled"` // Toggle for enabling Logs
	Logging        Logging `json:"logging"`         // logging structure

	// Time related fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	StoppedAt *time.Time `json:"stopped_at,omitempty"`
}

// Logging defines the configuration for logging, including the log path and log level.
type Logging struct {
	LogPath  string   `json:"log_path,omitempty"` // Log path
	LogLevel LogLevel `json:"log_level"`          // Log Level (Example - Silly , debug , info)
}

type ListenerRequest struct {
	Type        string `json:"type"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Config      any    `json:"config"`
}

type HTTPListenerConfig struct {
	WorkingHours string           `json:"working_hours"`
	Hosts        []string         `json:"hosts"`
	HostBind     string           `json:"host_bind"`
	HostRotation HostRotationType `json:"host_rotation"`
	PortBind     string           `json:"port_bind"`
	//PortConn         string           `json:"port_conn"`
	Secure           bool           `json:"secure"`
	UserAgent        string         `json:"user_agent"`
	Headers          []Header       `json:"headers"`
	Uris             []string       `json:"uris"`
	Certificate      TLSCertificate `json:"certificate"`
	WhitelistEnabled bool           `json:"whitelist_enabled"`
	Whitelist        []string       `json:"whitelist"`
	BlacklistEnabled bool           `json:"blacklist_enabled"`
	Blacklist        []string       `json:"blacklist"`
	ProxySettings    ProxySettings  `json:"proxy_settings"`
}

func (config *HTTPListenerConfig) Validate() error {

	portRangeStart := conf.GetListenerPortRangeStart()
	portRangeEnd := conf.GetListenerPortRangeEnd()

	if config.HostBind == "" {
		return fmt.Errorf("HostBind is required")
	}
	if config.PortBind == "" {
		return fmt.Errorf("PortBind is required")
	}

	portBindInt, err := strconv.Atoi(config.PortBind)
	if err != nil {
		return fmt.Errorf("PortBind must be a valid integer")
	}

	if portBindInt < portRangeStart || portBindInt > portRangeEnd {
		return fmt.Errorf("PortBind must be within the range %d-%d", portRangeStart, portRangeEnd)
	}

	if len(config.Hosts) == 0 {
		return fmt.Errorf("at least one host is required")
	}

	for _, host := range config.Hosts {
		if host == "" {
			return fmt.Errorf("hosts cannot contain empty values")
		}
	}

	if config.Secure {
		if config.Certificate.CertPath == "" || config.Certificate.KeyPath == "" {
			return fmt.Errorf("certificate paths are required for secure mode")
		}
	}

	if config.WhitelistEnabled && len(config.Whitelist) == 0 {
		return fmt.Errorf("whitelist is enabled but no whitelist entries are provided")
	}

	if config.BlacklistEnabled && len(config.Blacklist) == 0 {
		return fmt.Errorf("blacklist is enabled but no blacklist entries are provided")
	}

	for _, header := range config.Headers {
		if header.Key == "" || header.Value == "" {
			return fmt.Errorf("headers must have both key and value")
		}
	}

	if config.ProxySettings.Enabled {
		if config.ProxySettings.Type == "" {
			return fmt.Errorf("proxy type is required when proxy is enabled")
		}
		if config.ProxySettings.Port == "" {
			return fmt.Errorf("proxy port is required when proxy is enabled")
		}
	}

	return nil
}

type TCPListenerConfig struct {
	PortBind   string `json:"port_bind"`
	HostBind   string `json:"host_bind"`
	BufferSize int    `json:"buffer_size"`
	Timeout    int    `json:"timeout"`
}

type SMBListenerConfig struct {
	PipeName     string `json:"pipe_name"`
	MaxInstances int    `json:"max_instances"`
	KillDate     int64  `json:"kill_date"`
}

type ForeignListenerConfig struct {
	Endpoint       string         `json:"endpoint"`
	Authentication Authentication `json:"authentication"`
}

// Used by foreign listeners
type Authentication struct {
	Enabled  bool   `json:"enabled"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type TLSCertificate struct {
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

type HostRotationType string

const (
	Fallback   HostRotationType = "Fallback"
	Sequential HostRotationType = "Sequential"
	RoundRobin HostRotationType = "RoundRobin"
	Random     HostRotationType = "Random"
)

var ValidCallbackRotationTypes = []HostRotationType{
	Fallback,
	Sequential,
	RoundRobin,
	Random,
}
