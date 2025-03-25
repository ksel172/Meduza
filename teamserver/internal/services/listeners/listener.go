package services

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/ksel172/Meduza/teamserver/utils"
)

const (
	StatusPending     = "pending"     // Starting up resources
	StatusReady       = "ready"       // Idle, waiting for initialization/start
	StatusStarting    = "starting"    // Listener is being started
	StatusRunning     = "running"     // Running, server listening
	StatusStopping    = "stopping"    // Listener is stopping
	StatusTerminating = "terminating" // Listener is terminating

	LifecycleManaged   = "managed"   // Listener is managed by the manager and listen for changes
	LifecycleScheduled = "scheduled" // Listener is scheduled by the manager and polls for changes

	DeploymentLocal    = "local"    // Listener is deployed and managed locally within the same process
	DeploymentExternal = "external" // Listener is deployed anywhere else and communicates over the network
)

// Listener provides an abstraction over a listener of any kind
// Listeners MIGHT have an implementation or not
type Listener struct {

	// Metadata
	ID string

	// Listener operation configuration
	Config ListenerConfig
	mux    sync.RWMutex // Any writes to the listener should lock it, it can be read concurrently though

	// Lifecycle manager, differ based on the listener lifecycle
	lifecycleManager ListenerLifecycleManager

	// Listener concrete implementation
	listener ListenerImplementation
}

func NewListenerFromConfig(config ListenerConfig) (*Listener, error) {
	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	var lifecycleManager ListenerLifecycleManager
	if config.Lifecycle == LifecycleManaged {
		lifecycleManager = NewManagedLifecycleManager()
	} else if config.Lifecycle == LifecycleScheduled {
		lifecycleManager = NewScheduledLifecycleManager()
	} else {
		return nil, fmt.Errorf("unknown lifecycle: %s", config.Lifecycle)
	}

	// Create listener implementation from defaults if listener is local
	var impl ListenerImplementation
	if config.Deployment == DeploymentLocal {
		var err error
		impl, err = CreateImplementation(config.Kind)
		if err != nil {
			return nil, fmt.Errorf("failed to create listener implementation: %w", err)
		}
	}

	return &Listener{
		ID:               config.ID,
		Config:           config,
		lifecycleManager: lifecycleManager,
		listener:         impl,
	}, nil
}

func (l *Listener) Start(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.Config.Status != StatusReady {
		return errors.New("listener is not ready to start")
	}

	return l.lifecycleManager.Start(ctx, l)
}

func (l *Listener) Stop(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	if l.Config.Status != StatusRunning {
		return errors.New("listener is not running")
	}

	return l.lifecycleManager.Stop(ctx, l)
}

func (l *Listener) Terminate(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	return l.lifecycleManager.Terminate(ctx, l)
}

func (l *Listener) UpdateConfig(ctx context.Context, config ListenerConfig) error {
	errs := []error{errors.New("cannot update fields: ")}

	// Validate fields that cannot be updated remain the same
	if l.Config.Kind != config.Kind {
		errs = append(errs, errors.New("kind"))
	}
	if l.Config.Lifecycle != config.Lifecycle {
		errs = append(errs, errors.New("lifecycle"))
	}
	if l.Config.Deployment != config.Deployment {
		errs = append(errs, errors.New(DeploymentLocal))
	}

	if len(errs) > 1 {
		return errors.Join(errs...)
	}

	l.lifecycleManager.UpdateConfig(ctx, l, config)

	return nil
}

// External listeners should use this to update listener status
// The listener is sendign a response back to confirm it received and performed
// the requested operation asynchronously
func (l *Listener) UpdateStatus(ctx context.Context, status string) {
	utils.AssertEquals(l.Config.Kind, DeploymentExternal)

	l.mux.Lock()
	defer l.mux.Unlock()

	l.Config.Status = status
}
