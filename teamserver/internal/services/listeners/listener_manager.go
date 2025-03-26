package services

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// ListenerManager handles the lifecycle of listeners
type ListenerManager struct {
	// Manager start and stop timeout config
	startTimeout int
	stopTimeout  int

	listenerDal IListenerDAL

	// Track synchronization timestamps for external listeners
	synchronizationLog map[string]time.Time
	syncMux            sync.Mutex
}

func NewListenerManager(listenerDAL IListenerDAL) *ListenerManager {
	return &ListenerManager{
		startTimeout:       15,
		stopTimeout:        15,
		listenerDal:        listenerDAL,
		synchronizationLog: make(map[string]time.Time),
	}
}

// Changed to match DAL return type
func (m *ListenerManager) GetListeners(ctx context.Context) ([]*Listener, error) {
	return m.listenerDal.GetAllListeners(ctx)
}

// GetListenerStatuses returns a map of all listener statuses
func (m *ListenerManager) GetListenerStatuses(ctx context.Context) (map[string]string, error) {
	return m.getListenerStatuses(ctx)
}

func (m *ListenerManager) AddListener(ctx context.Context, config ListenerConfig) error {
	return m.addListener(ctx, config)
}

// StartListener starts a listener (public wrapper for testing)
func (m *ListenerManager) StartListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	return m.startListener(ctx, listenerID, errChan)
}

// StopListener stops a running listener (public wrapper for testing)
func (m *ListenerManager) StopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	return m.stopListener(ctx, listenerID, errChan)
}

// TerminateListener terminates a listener (public wrapper for testing)
func (m *ListenerManager) TerminateListener(ctx context.Context, listenerID string) error {
	return m.terminateListener(ctx, listenerID)
}

// Loop to watch if any of the listeners failed to comply
func (m *ListenerManager) watchListeners(ctx context.Context) {
	listeners, err := m.listenerDal.GetAllListeners(ctx)
	if err != nil {
		// Log error and return
		return
	}

	terminateQueue := []string{}
	for _, listener := range listeners { // Use direct pointer, no need for &
		if listener.Config.Deployment != DeploymentExternal {
			continue
		}

		m.syncMux.Lock()
		lastSync, ok := m.synchronizationLog[listener.Config.ID]
		m.syncMux.Unlock()

		if !ok {
			// Initialize sync record if not exists
			m.syncMux.Lock()
			m.synchronizationLog[listener.Config.ID] = time.Now()
			m.syncMux.Unlock()
			continue
		}

		// If a listener has exceeded the time allowed it can miss communications
		if time.Since(lastSync) > time.Duration(listener.Config.Heartbeat*5) {
			terminateQueue = append(terminateQueue, listener.Config.ID)
		}
	}

	// Rest of the function remains the same
}

// Also update the getListenerStatuses method
func (m *ListenerManager) getListenerStatuses(ctx context.Context) (map[string]string, error) {
	listeners, err := m.listenerDal.GetAllListeners(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get listeners: %w", err)
	}

	listenerStatuses := make(map[string]string)
	for _, listener := range listeners { // Use direct pointer, no need for &
		listenerStatuses[listener.Config.ID] = listener.Config.Status
	}
	return listenerStatuses, nil
}

// Rest of the methods can remain the same

// Update AutoStart method as well
func (m *ListenerManager) AutoStart(ctx context.Context) error {
	listeners, err := m.listenerDal.GetActiveListeners(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active listeners: %w", err)
	}

	for _, listener := range listeners { // Use direct pointer, no need for &
		if listener.Config.Status == StatusRunning {
			continue
		}

		if err := m.startListener(ctx, listener.Config.ID, make(chan<- error)); err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}
	}

	return nil
}
