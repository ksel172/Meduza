package listener

import _ "github.com/go-playground/validator/v10"

const (
	// Possible listener statuses
	StatusPending     = "pending"     // Starting up resources
	StatusReady       = "ready"       // Idle, waiting for initialization/start
	StatusStarting    = "starting"    // Listener is being started
	StatusRunning     = "running"     // Running, server listening
	StatusStopping    = "stopping"    // Listener is stopping
	StatusTerminating = "terminating" // Listener is terminating

	// Listener lifecycle modes
	LifecycleManaged   = "managed"   // Listener is managed by the manager and listen for changes
	LifecycleScheduled = "scheduled" // Listener is scheduled by the manager and polls for changes

	// Supported listener kinds
	HTTPListenerKind     string = "http"
	TCPListenerKind      string = "tcp"
	SMBListenerKind      string = "smb"
	ExternalListenerKind string = "external"
)
