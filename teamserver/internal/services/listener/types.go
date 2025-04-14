package listener

import _ "github.com/go-playground/validator/v10"

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
)
