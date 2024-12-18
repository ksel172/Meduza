package listeners

import (
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
