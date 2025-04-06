package listener

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// Listener is a representation of a listener of any kind
type Listener struct {
	models.Listener // embed all of the data fields in the listener data model

	mux sync.RWMutex

	// Lifecycle manager, differ based on the listener lifecycle
	lifecycleManager ListenerLifecycleManager

	// Listener concrete implementation
	listener ListenerImplementation
}

// This function is called after a listener is retrieved from storage
// The lifecycleManager and ListenerImplementation fields will be nil
// we must check how the listener is setup to run and prepare the fields
// for usage
func CreateListenerFromModel(listenerModel models.Listener) (*Listener, error) {
	listener := Listener{}
	listener.Listener = listenerModel

	switch listener.Lifecycle {
	case LifecycleManaged:
		listener.lifecycleManager = NewManagedLifecycleManager()
	case LifecycleScheduled:
		listener.lifecycleManager = NewScheduledLifecycleManager()
	default:
		return nil, fmt.Errorf("invalid lifecycle: %s", listener.Lifecycle)
	}

	// Create the concrete listener implementation
	listenerImplementation, err := CreateListenerImplementation(listener.Kind, listener.RawConfig)
	if err != nil {
		return nil, err
	}
	listener.listener = listenerImplementation

	// Validate configuration
	if err := listener.ValidateConfig(); err != nil {
		return nil, fmt.Errorf("failed to validate config: %w", err)
	}

	return &listener, nil
}

func (l *Listener) Start(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.Status != StatusReady {
		return errors.New("listener is not ready to start")
	}

	return l.lifecycleManager.Start(ctx, l)
}

func (l *Listener) Stop(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.Status != StatusRunning {
		return errors.New("listener is not running")
	}

	return l.lifecycleManager.Stop(ctx, l)
}

func (l *Listener) Terminate(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	return l.lifecycleManager.Terminate(ctx, l)
}

// This operation should not affect running listeners depending on the update
func (l *Listener) UpdateConfig(ctx context.Context, newConfig *Listener) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	// Update the listener configuration
	l.Kind = newConfig.Kind
	l.Name = newConfig.Name
	l.Description = newConfig.Description
	l.Status = newConfig.Status
	l.Heartbeat = newConfig.Heartbeat
	l.RawConfig = newConfig.RawConfig
	l.Lifecycle = newConfig.Lifecycle
	l.UpdatedAt = time.Now()

	return nil
}

// External listeners should use this to update listener status
// The listener is sendign a response back to confirm it received and performed
// the requested operation asynchronously
func (l *Listener) UpdateStatus(ctx context.Context, status string) {
	utils.AssertEquals(l.Kind, ExternalListenerKind)

	l.mux.Lock()
	defer l.mux.Unlock()

	l.Status = status
}
