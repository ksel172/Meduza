package models

import "time"

const ParamListenerID string = "listener_id"

// This exists to prevent a circular import in the database layer
// Database returns only the data fields of a listener
type Listener struct {
	ID          string `json:"id"`
	Kind        string `json:"kind" validate:"required"` // http, tcp, smb, custom, etc
	Status      string `json:"status"`                   // running, stopped etc
	Name        string `json:"name"`
	Description string `json:"description"`

	// eventually add tags, tags can be created and are stored in another table
	// reference from tags table, many to many relationship
	// Tags        []string `json:"tags"`

	// These configurations are only allowed for external listeners
	// Not controlled within the same process as the C2 server
	Heartbeat int `json:"heartbeat"` //

	// Config holds implementation specific configs for external listeners, otherwise they are accessed from the listener field
	Config any `json:"config" validate:"required"`

	// Listener operation method
	Lifecycle  string `json:"lifecycle" validate:"oneof:scheduled managed"`
	Deployment string `json:"deployment" validate:"oneof:local external"`

	// Auditability fields
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	StartedAt time.Time `json:"started_at,omitempty"`
	StoppedAt time.Time `json:"stopped_at,omitempty"`
}
