package listener

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	http_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/http"
	smb_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/smb"
	tcp_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/tcp"
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

	Validate() error // Listener validates its configuration
}

// Creates the local listener implementation based on the provided config byte array
func CreateListenerImplementation(kind string, config json.RawMessage) (ListenerImplementation, error) {
	switch kind {

	case HTTPListenerKind:
		var httpConfig http_listener.HTTPListenerConfig
		if err := json.Unmarshal(config, &httpConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal HTTP config: %w", err)
		}

		implementation, err := http_listener.NewHTTPListener(httpConfig, &checkin.CheckInController{})
		if err != nil {
			return nil, fmt.Errorf("failed to create http implementation: %w", err)
		}
		return implementation, nil

	case TCPListenerKind:
		return &tcp_listener.TCPListener{}, nil

	case SMBListenerKind:
		return &smb_listener.SMBListener{}, nil

	default:
		return nil, fmt.Errorf("unsupported listener kind: %s", kind)
	}
}
