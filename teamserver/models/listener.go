package models

import (
	"time"
)

// Listener Types (Enum)
const (
	ListenerTypeHTTP    = "http"
	ListenerTypeHTTPS   = "https"
	ListenerTypeTCP     = "tcp"
	ListenerTypeSMB     = "smb"
	ListenerTypeForeign = "external"
)

// AllowedListenerTypes - Valid listener types
var AllowedListenerTypes = []string{
	ListenerTypeHTTP,
	ListenerTypeHTTPS,
	ListenerTypeTCP,
	ListenerTypeSMB,
	ListenerTypeForeign,
}

// Listener Statuses (Enum)
const (
	StatusStopped    = 0
	StatusRunning    = 1
	StatusPaused     = 2
	StatusProcessing = 4
	StatusError      = 5
)

// Listener - Model stored in Redis
type Listener struct {
	ID          string `json:"id"`     // UUID
	Type        string `json:"type"`   // Listener Type (http, tcp, etc.)
	Host        string `json:"host"`   // IP or hostname
	Port        int    `json:"port"`   // Port number
	Status      int    `json:"status"` // 0 = stopped, 1 = running, 2 = paused, 3 = processing
	Description string `json:"description,omitempty"`

	// Time related fields
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	StartedAt *time.Time `json:"startedAt,omitempty"`
	StoppedAt *time.Time `json:"stoppedAt,omitempty"`

	// SSL/TLS
	certPath string `json:"certPath,omitempty"` // Certificate (HTTPS)
	keyPath  string `json:"keyPath,omitempty"`  // Private Key (HTTPS)

	// Whitelisting and Blacklisting
	WhitelistEnabled bool     `json:"whitelistEnabled"`
	Whitelist        []string `json:"whitelist,omitempty"`
	BlacklistEnabled bool     `json:"blacklistEnabled"`
	Blacklist        []string `json:"blacklist,omitempty"`

	// Logging
	LoggingEnabled bool   `json:"loggingEnabled"`
	LogPath        string `json:"logPath,omitempty"`
	LogLevel       string `json:"logLevel"`
}
