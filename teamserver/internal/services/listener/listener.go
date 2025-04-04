package listener

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	http_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/http"
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
	l := Listener{}
	l.Listener = listenerModel

	switch l.Lifecycle {
	case LifecycleManaged:
		l.lifecycleManager = NewManagedLifecycleManager()
	case LifecycleScheduled:
		l.lifecycleManager = NewScheduledLifecycleManager()
	default:
		return nil, fmt.Errorf("invalid lifecycle: %s", l.Lifecycle)
	}

	switch l.Deployment {
	case DeploymentExternal:
		// TODO: Add external deployment support
	case DeploymentLocal:
		switch l.Kind {
		case HTTPListenerKind:
			var httpConfig http_listener.HTTPListenerConfig
			if err := utils.MapToStruct(l.Config, &httpConfig); err != nil {
				return nil, fmt.Errorf("failed to convert config to HTTPListenerConfig: %w", err)
			}

			if err := httpConfig.ValidateConfig(); err != nil {
				return nil, fmt.Errorf("invalid HTTP listener config: %w", err)
			}

			listenerImplementation, err := http_listener.NewHTTPListener(httpConfig, &checkin.CheckInController{})
			if err != nil {
				return nil, fmt.Errorf("failed to create HTTP listener: %w", err)
			}
			l.listener = listenerImplementation
		default:
			return nil, fmt.Errorf("unsupported listener kind: %s", l.Kind)
		}
		// TODO: Add other cases here for different listener kinds
	default:
		return nil, fmt.Errorf("invalid deployment: %s", l.Deployment)
	}

	return &l, nil
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
