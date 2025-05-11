package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ksel172/Meduza/teamserver/pkg/conf"
)

/*
Listener possible status:

Local:
Ready - runtime NOT mapped, only stored in database
Starting - runtime mapped, not listening to requests
Running - runtime mapped, listening to requests
Stopping - runtime mapped
Terminating - runtime in the process of being unmapped
Failed - operation failed, requires cleanup and restart in some cases

External:
Ready - runtime mapped, ready to receive start request
Starting - runtime mapped, not listening to requests
Running - runtime mapped, listening to requests
Stopping - runtime mapped
Terminating - runtime in the process of being unmapped, resources cleaned up, program should exit
Failed - operation failed, requires cleanup and restart in some cases
*/

const (
	// Possible listener statuses
	StatusPending     = "pending"     // No resource created, waiting for any signal to start up
	StatusReady       = "ready"       // Idle, waiting for initialization/start
	StatusStarting    = "starting"    // Listener is being started
	StatusRunning     = "running"     // Running, server listening
	StatusStopping    = "stopping"    // Listener is stopping
	StatusTerminating = "terminating" // Listener is terminating
	StatusFailed      = "failed"      // Listener failed to start / crashed

	// Listener lifecycle modes
	LifecycleManaged   = "managed"   // Listener is managed by the manager and listen for changes
	LifecycleScheduled = "scheduled" // Listener is scheduled by the manager and polls for changes

	// Supported listener kinds
	HTTPListenerKind     string = "http"
	TCPListenerKind      string = "tcp"
	SMBListenerKind      string = "smb"
	ExternalListenerKind string = "external"

	// Parameter names
	ParamListenerID string = "listener_id"
)

// Database returns only the data fields of a listener
type Listener struct {
	ID          string `json:"id"`
	Kind        string `json:"kind" validate:"required"` // http, tcp, smb, custom, etc
	Status      string `json:"status"`                   // running, stopped etc
	Name        string `json:"name"`
	Description string `json:"description"`

	// External only fields
	External  bool   `json:"external"`  // true if the listener is external, false if it is local
	Host      string `json:"host"`      // if local, localhost
	Port      int    `json:"port"`      // if local, port assigned automatically by port manager
	Heartbeat int    `json:"heartbeat"` // To check if external listener is alive

	// Config holds implementation specific configs for external listeners, otherwise they are accessed from the listener field
	RawConfig json.RawMessage `json:"config" validate:"required"`

	// eventually add tags, tags can be created and are stored in another table
	// reference from tags table, many to many relationship
	// Tags        []string `json:"tags"`

	// Auditability fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StartedAt time.Time `json:"started_at,omitempty"`
	StoppedAt time.Time `json:"stopped_at,omitempty"`
}

func (l *Listener) Validate() error {

	portRangeStart := conf.GetListenerPortRangeStart()
	portRangeEnd := conf.GetListenerPortRangeEnd()

	if l.External {
		if l.Host == "" {
			return fmt.Errorf("host is required")
		}
		if l.Port == 0 {
			return fmt.Errorf("port is required")
		}
	}

	if l.Kind == "" {
		return fmt.Errorf("kind is required")
	}

	if l.Port < portRangeStart || l.Port > portRangeEnd {
		return fmt.Errorf("port must be between %d and %d", portRangeStart, portRangeEnd)
	}

	return nil
}

type CreateLocalListenerRequest struct {
	Kind        string `json:"kind" validate:"required"` // http, tcp, smb, custom, etc
	Name        string `json:"name"`
	Description string `json:"description"`
	Heartbeat   int    `json:"heartbeat"` //
}

func (clr CreateLocalListenerRequest) IntoListener() Listener {

	heartbeat := clr.Heartbeat
	if clr.Heartbeat < 30 {
		heartbeat = 30
	}

	return Listener{
		Kind:        clr.Kind,
		Status:      StatusReady,
		Name:        clr.Name,
		Description: clr.Description,

		// External fields
		External:  false,
		Host:      "localhost",
		Port:      8010, // Must replace with port manager service implementation later
		RawConfig: json.RawMessage(`{}`),

		Heartbeat: heartbeat,
	}
}
