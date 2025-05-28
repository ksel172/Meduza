package listener

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	http_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/http"
	smb_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/smb"
	tcp_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/tcp"
	"github.com/ksel172/Meduza/teamserver/models"
)

type ListenerActionFunc func(context.Context) error

type ListenerImplementation interface {
	Start(context.Context) error        // Start starts a listener with Ready status
	Stop(context.Context) error         // Stop simply stops a listener from listening. It will still be active and sending heartbeats.
	Terminate(context.Context) error    // Close kills a listener process.
	UpdateConfig(context.Context) error // Listener updates its own configuration

	Validate() error // Listener validates its configuration
}

// This function is called after a listener is retrieved from storage
// The lifecycleManager and ListenerImplementation fields will be nil
// we must check how the listener is setup to run and prepare the fields
// for usage
func createListenerFromModel(listenerModel models.Listener, checkinController checkin.ICheckInController) (*Listener, error) {
	listener := Listener{
		Listener: listenerModel,
	}

	// Create the concrete listener implementation
	listenerImplementation, err := createListenerImplementation(listener.Kind, listener.Host, listener.Port, listener.RawConfig, checkinController)
	if err != nil {
		return nil, err
	}
	listener.listener = listenerImplementation

	// Initialize statusUpdates channel
	// Allow buffered updates
	listener.statusUpdatesCh = make(chan string, 3)

	// Validate configuration
	if err := listener.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &listener, nil
}

// Creates the local listener implementation based on the provided config byte array
func createListenerImplementation(kind string, host string, port int, config json.RawMessage,
	checkinController checkin.ICheckInController) (ListenerImplementation, error) {
	switch kind {

	case models.HTTPListenerKind:
		var httpConfig http_listener.HTTPListenerConfig
		if err := json.Unmarshal(config, &httpConfig); err != nil {
			return nil, fmt.Errorf("failed to unmarshal HTTP config: %w", err)
		}

		// Inject host and port into http listener config
		httpConfig.Host = host
		httpConfig.Port = port

		implementation, err := http_listener.NewHTTPListener(httpConfig, checkinController)
		if err != nil {
			return nil, fmt.Errorf("failed to create http implementation: %w", err)
		}
		return implementation, nil

	case models.TCPListenerKind:
		return &tcp_listener.TCPListener{}, nil

	case models.SMBListenerKind:
		return &smb_listener.SMBListener{}, nil

	default:
		return nil, fmt.Errorf("unsupported listener kind: %s", kind)
	}
}

func ValidateListenerConfig(kind string, config json.RawMessage) error {
	switch kind {
	case models.HTTPListenerKind:
		var httpConfig http_listener.HTTPListenerConfig
		if err := json.Unmarshal(config, &httpConfig); err != nil {
			return fmt.Errorf("failed to unmarshal HTTP config: %w", err)
		}

	case models.TCPListenerKind:
		var tcpConfig tcp_listener.TCPListenerConfig
		if err := json.Unmarshal(config, &tcpConfig); err != nil {
			return fmt.Errorf("failed to unmarshal TPC config: %w", err)
		}

	case models.SMBListenerKind:
		var smbConfig smb_listener.SMBListenerConfig
		if err := json.Unmarshal(config, &smbConfig); err != nil {
			return fmt.Errorf("failed to unmarshal SMB config: %w", err)
		}

	default:
		return fmt.Errorf("unsupported listener kind: %s", kind)
	}

	return nil
}
