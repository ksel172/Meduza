package listeners

import (
	"context"
	"errors"

	http_listener "github.com/ksel172/Meduza/teamserver/pkg/listeners/http"
	smb_listener "github.com/ksel172/Meduza/teamserver/pkg/listeners/smb"
	tcp_listener "github.com/ksel172/Meduza/teamserver/pkg/listeners/tcp"
)

/*
This is only used for Local listener deployments

Listener implementation specs

 1. Custom protocol.
    HTTP, TCP, SMB ...
 2. Custom programming language. Use whatever language to code listeners.
 3. Listeners should declare a version in their communications.
    That is the version of the API the listenerController will be using.
 4. Listeners should handle the following tasks
    a. Agent communication
    Decryption/encryption
    b. Request forwarding.
    Should parse agent request and forward to the listener controller
*/
type ListenerImplementation interface {
	Start(context.Context) error        // Start starts a listener with Ready status
	Stop(context.Context) error         // Stop simply stops a listener from listening. It will still be active and sending heartbeats.
	Terminate(context.Context) error    // Close kills a listener process.
	UpdateConfig(context.Context) error // Listener updates its own configuration
}

func CreateImplementation(kind string) (ListenerImplementation, error) {
	switch kind {
	case "http":
		return &http_listener.HTTPListener{}, nil
	case "tcp":
		return &tcp_listener.TCPListener{}, nil
	case "smb":
		return &smb_listener.SMBListener{}, nil
	default:
		return nil, errors.New("unsupported listener kind")
	}
}
