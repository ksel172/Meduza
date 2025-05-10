package listener

import (
	"context"
	"errors"
	"sync"

	"github.com/ksel172/Meduza/teamserver/models"
)

// TODO
// Recreate listener implementations after a cycle of running-stopping is finished
// Internal server cannot be reused after Shutdown is called

// Listener is a representation of a listener of any kind
type Listener struct {
	models.Listener // embed all of the data fields in the listener data model
	listener        ListenerImplementation
	mux             sync.RWMutex
	statusUpdatesCh chan string
}

func (l *Listener) Start(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	// Initialize channel for the duration of the operation
	l.statusUpdatesCh = make(chan string, 3)
	defer close(l.statusUpdatesCh)

	if l.Status != StatusReady {
		return errors.New("listener is not ready to start")
	}

	l.updateStatus(StatusStarting)
	if err := l.listener.Start(ctx); err != nil {
		l.updateStatus(StatusFailed)
		return err
	}
	l.updateStatus(StatusRunning)

	return nil
}

func (l *Listener) Stop(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	// Initialize channel for the duration of the operation
	l.statusUpdatesCh = make(chan string, 3)
	defer close(l.statusUpdatesCh)

	if l.Status != StatusRunning {
		return errors.New("listener is not running")
	}

	l.updateStatus(StatusStopping)
	if err := l.listener.Stop(ctx); err != nil {
		l.updateStatus(StatusFailed)
		return err
	}
	l.updateStatus(StatusReady)

	return nil
}

func (l *Listener) Terminate(ctx context.Context) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	// Initialize channel for the duration of the operation
	l.statusUpdatesCh = make(chan string, 3)
	defer close(l.statusUpdatesCh)

	l.updateStatus(StatusTerminating)
	if err := l.listener.Terminate(ctx); err != nil {
		l.updateStatus(StatusFailed)
		return err
	}
	l.updateStatus(StatusPending)

	return nil
}

// This operation should not affect running listeners depending on the update
func (l *Listener) UpdateConfig(ctx context.Context, listenerUpdate models.Listener) error {
	l.mux.Lock()
	defer l.mux.Unlock()

	l.Name = listenerUpdate.Name
	l.Description = listenerUpdate.Description
	l.Heartbeat = listenerUpdate.Heartbeat
	l.RawConfig = listenerUpdate.RawConfig

	return nil
}

// This function assumes the lock is already held
// Updates the listener status and sends status update to channel
func (l *Listener) updateStatus(status string) {
	prevStatus := l.Status
	l.Status = status

	// Notify if channel is provided and status changed
	if l.statusUpdatesCh != nil && prevStatus != status {
		l.statusUpdatesCh <- status
	}
}
