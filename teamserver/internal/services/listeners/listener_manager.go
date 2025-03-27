package services

import (
	"context"
	"errors"
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

func (m *ListenerManager) GetListeners(ctx context.Context) ([]Listener, error) {
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
	for i := range listeners {
		listener := &listeners[i] // Use pointer to avoid copying mutex
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
	for i := range listeners {
		listener := &listeners[i] // Use pointer to avoid copying mutex
		listenerStatuses[listener.Config.ID] = listener.Config.Status
	}
	return listenerStatuses, nil
}

// Adding a listener to a listenerController does not start it.
func (m *ListenerManager) addListener(ctx context.Context, listenerConfig ListenerConfig) error {
	// Check if listener with same ID already exists
	_, err := m.listenerDal.GetListenerById(ctx, listenerConfig.ID)
	if err == nil {
		return fmt.Errorf("listener with ID %s already exists", listenerConfig.ID)
	}

	// Create listener object
	listener, err := NewListenerFromConfig(listenerConfig)
	if err != nil {
		return fmt.Errorf("failed to create listener config: %w", err)
	}

	// Add to DAL
	if err := m.listenerDal.CreateListener(ctx, listener); err != nil {
		return fmt.Errorf("failed to store listener: %w", err)
	}

	// Initialize synchronization record for this listener
	m.syncMux.Lock()
	m.synchronizationLog[listener.Config.ID] = time.Now()
	m.syncMux.Unlock()

	return nil
}

// Starts an already existing listener
func (m *ListenerManager) startListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listener, err := m.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(m.startTimeout)*time.Second)
		defer cancel()
		if err := listener.Start(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start listener: %w", err)
			return
		}

		// Update listener status in DAL
		updates := map[string]any{"status": "running"}
		if err := m.listenerDal.UpdateListener(ctx, listenerID, updates); err != nil {
			errChan <- fmt.Errorf("failed to update listener status: %w", err)
			return
		}
		close(errChan)
	}()

	return nil
}

// Stops a running listener
func (m *ListenerManager) stopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listener, err := m.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(m.stopTimeout)*time.Second)
		defer cancel()
		if err := listener.Stop(ctx); err != nil {
			errChan <- fmt.Errorf("failed to stop listener: %w", err)
			return
		}

		// Update listener status in DAL
		updates := map[string]any{"status": "stopped"}
		if err := m.listenerDal.UpdateListener(ctx, listenerID, updates); err != nil {
			errChan <- fmt.Errorf("failed to update listener status: %w", err)
			return
		}
		close(errChan)
	}()

	return nil
}

// Terminate a listener and remove from registry
func (m *ListenerManager) terminateListener(ctx context.Context, listenerID string) error {
	listener, err := m.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID '%s' not found: %w", listenerID, err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(m.stopTimeout)*time.Second)
	defer cancel()
	if err := listener.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	// Remove from DAL
	if err := m.listenerDal.DeleteListener(ctx, listenerID); err != nil {
		return fmt.Errorf("failed to delete listener from storage: %w", err)
	}

	// Clean up synchronization record
	m.syncMux.Lock()
	delete(m.synchronizationLog, listenerID)
	m.syncMux.Unlock()

	return nil
}

// Updates a listener's config
func (m *ListenerManager) updateListener(ctx context.Context, listenerConfig ListenerConfig) error {
	listener, err := m.listenerDal.GetListenerById(ctx, listenerConfig.ID)
	if err != nil {
		return fmt.Errorf("listener with ID '%s' not found: %w", listenerConfig.ID, err)
	}

	// Update listener config
	if err := listener.UpdateConfig(ctx, listenerConfig); err != nil {
		return fmt.Errorf("failed to update listener config: %w", err)
	}

	// Save updated listener to DAL
	updates := map[string]any{
		"config": listenerConfig,
	}
	return m.listenerDal.UpdateListener(ctx, listenerConfig.ID, updates)
}

// External listener only
// Update a listener status
func (m *ListenerManager) updateListenerStatus(ctx context.Context, listenerID, status string) error {
	listener, err := m.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener not found: %w", err)
	}

	if listener.Config.Deployment != DeploymentExternal {
		return errors.New("operation not allowed for local listeners")
	}

	listener.UpdateStatus(ctx, status)

	// Update in DAL
	updates := map[string]any{"status": status}
	return m.listenerDal.UpdateListener(ctx, listenerID, updates)
}

// External listener only
// Synchronizes configurations with external listeners
func (m *ListenerManager) synchronize(ctx context.Context, listenerID string) (ListenerConfig, error) {
	listener, err := m.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return ListenerConfig{}, fmt.Errorf("listener not found: %w", err)
	}

	// Update last synchronization time
	m.syncMux.Lock()
	m.synchronizationLog[listenerID] = time.Now()
	m.syncMux.Unlock()

	return listener.Config, nil
}

func (m *ListenerManager) AutoStart(ctx context.Context) error {
	listeners, err := m.listenerDal.GetActiveListeners(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active listeners: %w", err)
	}

	for i := range listeners {
		listener := &listeners[i]
		if listener.Config.Status == StatusRunning {
			continue
		}

		if err := m.startListener(ctx, listener.Config.ID, make(chan<- error)); err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}
	}

	return nil
}
