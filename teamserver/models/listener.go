package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ksel172/Meduza/teamserver/pkg/conf"
)

const ParamListenerID string = "listener_id"

// This exists to prevent a circular import in the database layer
// Database returns only the data fields of a listener
type Listener struct {
	ID          string `json:"id"`
	Kind        string `json:"kind" validate:"required"` // http, tcp, smb, custom, etc
	Status      string `json:"status"`                   // running, stopped etc
	Name        string `json:"name"`
	Description string `json:"description"`

	// External only fields
	External bool   `json:"external"` // true if the listener is external, false if it is local
	Host     string `json:"host"`
	Port     int    `json:"port"`

	// eventually add tags, tags can be created and are stored in another table
	// reference from tags table, many to many relationship
	// Tags        []string `json:"tags"`

	// These configurations are only allowed for external listeners
	// Not controlled within the same process as the C2 server
	Heartbeat int `json:"heartbeat"` //

	// Config holds implementation specific configs for external listeners, otherwise they are accessed from the listener field
	RawConfig json.RawMessage `json:"config" validate:"required"`

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
