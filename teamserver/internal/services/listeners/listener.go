package services

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ksel172/Meduza/teamserver/utils"
)

const (
	ParamListenerID string = "listener_id"
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

type Listener struct {
	ID   string `json:"id"`
	Type string `json:"type"` // http, tcp, smb, custom, etc

	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"` // running, stopped etc

	Heartbeat int `json:"heartbeat"`
	Config    any `json:"config"`

	// fixed specs, cannot be modified after set unless listener is restarted
	Lifecycle  string `json:"lifecycle" validate:"oneof:scheduled managed"` // scheduled || managed
	Deployment string `json:"deployment" validate:"oneof:local external"`

	// Time related fields
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	StartedAt *time.Time `json:"started_at,omitempty"`
	StoppedAt *time.Time `json:"stopped_at,omitempty"`

	mux sync.RWMutex

	// Lifecycle manager, differ based on the listener lifecycle
	lifecycleManager ListenerLifecycleManager

	// Listener concrete implementation
	listener ListenerImplementation
}

func NewListenerFromBase(base *Listener) (*Listener, error) {
	// Initialize new listener based on the base configuration
	newListener := &Listener{
		ID:               base.ID,
		Type:             base.Type,
		Name:             base.Name,
		Description:      base.Description,
		Status:           base.Status,
		Heartbeat:        base.Heartbeat,
		Config:           base.Config,
		Lifecycle:        base.Lifecycle,
		Deployment:       base.Deployment,
		CreatedAt:        base.CreatedAt,
		UpdatedAt:        base.UpdatedAt,
		StartedAt:        base.StartedAt,
		StoppedAt:        base.StoppedAt,
		lifecycleManager: base.lifecycleManager,
		listener:         base.listener,
	}

	return newListener, nil
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

func (l *Listener) UpdateConfig(ctx context.Context, newConfig *Listener) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	// Update the listener configuration
	l.Type = newConfig.Type
	l.Name = newConfig.Name
	l.Description = newConfig.Description
	l.Status = newConfig.Status
	l.Heartbeat = newConfig.Heartbeat
	l.Config = newConfig.Config
	l.Lifecycle = newConfig.Lifecycle
	l.Deployment = newConfig.Deployment
	l.UpdatedAt = time.Now()

	return nil
}

// External listeners should use this to update listener status
// The listener is sendign a response back to confirm it received and performed
// the requested operation asynchronously
func (l *Listener) UpdateStatus(ctx context.Context, status string) {
	// Fixed: Check Deployment field instead of Type
	utils.AssertEquals(l.Deployment, DeploymentExternal)

	l.mux.Lock()
	defer l.mux.Unlock()

	l.Status = status
}
