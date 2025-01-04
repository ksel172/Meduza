package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type (
	Id       = uuid.UUID
	status   = int32
	LogLevel = string
)

// Listener Statuses (Enum)
const (
	StatusStopped    status = 0
	StatusRunning    status = 1
	StatusPaused     status = 2
	StatusProcessing status = 3
	StatusError      status = 4
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
	ID          Id     `json:"id"`                    // UUID
	Type        string `json:"type"`                  // Listener Type (http, tcp, etc.)
	Name        string `json:"name"`                  // Name
	Status      int    `json:"status"`                // 0 = stopped, 1 = running, 2 = paused, 3 = processing
	Description string `json:"description,omitempty"` // description
	Config      any    `json:"config"`                // Configuration of the Listener

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

// Validate ensures the configuration is valid before use.
func (config *HTTPListenerConfig) Validate() error {
	if config.HostBind == "" {
		return fmt.Errorf("HostBind is required")
	}
	if config.PortBind == "" {
		return fmt.Errorf("PortBind is required")
	}
	/*
		if len(config.Uris) == 0 {
			return fmt.Errorf("At least one URI is required")
		}
	*/
	if config.Secure {
		if config.Certificate.CertPath == "" || config.Certificate.KeyPath == "" {
			return fmt.Errorf("Certificate paths are required for secure mode")
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
