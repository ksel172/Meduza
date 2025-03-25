package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"
)

// Managers can manage both LocalListeners:
// deployed within the same process, supported only by default listener kinds
// and ExternalListeners:
// both default + custom listeners, communication through the network
type Manager struct {
	listeners    map[string]*Listener // Managed listeners
	listenersMux sync.RWMutex         // Listeners RWmutex

	// Implement some mechanism to log listener synchronization calls (heartbeat + sync in one)
	// to remove listeners from jurisdiction if they miss 10 heartbeats in a row
	synchronizationLog map[string]time.Time
	syncMux            sync.Mutex

	// Local listeners
	// usedPorts []int      // List of currently used ports
	// portsMux  sync.Mutex // Mutex for ports availability

	// Manager start and stop timeout config
	startTimeout int
	stopTimeout  int
	// terminateTimeout int

	// Move this to the Controller struct
	// External listeners
	// allowedIPs []string // What IPs are allowed to connect with the controller
}

func (m *Manager) GetListeners() {
	panic("unimplemented")
}

// GetListenerStatuses returns a map of all listener statuses (public method for testing)
func (m *Manager) GetListenerStatuses() map[string]string {
	return m.getListenerStatuses()
}

func (m *Manager) AddListener(config ListenerConfig) error {
	return m.addListener(config)
}

// StartListener starts a listener (public wrapper for testing)
func (m *Manager) StartListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	return m.startListener(ctx, listenerID, errChan)
}

// StopListener stops a running listener (public wrapper for testing)
func (m *Manager) StopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	return m.stopListener(ctx, listenerID, errChan)
}

// TerminateListener terminates a listener (public wrapper for testing)
func (m *Manager) TerminateListener(ctx context.Context, listenerID string) error {
	return m.terminateListener(ctx, listenerID)
}

// Add synchronization loop to remove inactive listeners, 10 * listener.Heartbeat
// Add start up sequence after creating the manager -> starting up local listeners etc

func NewManager(listeners map[string]*Listener) (*Manager, error) {
	return &Manager{
		listeners:          listeners,
		synchronizationLog: make(map[string]time.Time),
		startTimeout:       15,
		stopTimeout:        15,
	}, nil
}

// Loop to watch if any of the listeners failed to comply
func (m *Manager) watchListeners() {
	// Read a consistent state of the listeners, push all listeners that need to be removed to a queue
	m.listenersMux.RLock()
	terminateQueue := []string{}
	for listenerID, listener := range m.listeners {
		if listener.Config.Deployment != DeploymentExternal {
			continue
		}

		// This should never happen, so we panic in case it happens
		lastSync, ok := m.synchronizationLog[listenerID]
		if !ok {
			panic(fmt.Sprintf("no synchronization log for listener: %s", listenerID))
		}

		// If a listener has exceeded the time allowed it can miss communications with the controller
		// we will try to send a terminate command anyway
		// This serves one purpose for both deployment kinds: to remove the listener from the controller's jurisdiction
		// But there is one added benefit for managed listeners: one last attempt to terminate and prevent resource leakage
		if time.Since(lastSync) > time.Duration(listener.Config.Heartbeat*5) {
			terminateQueue = append(terminateQueue, listenerID)
		}
	}
	m.listenersMux.RUnlock()

	if len(terminateQueue) == 0 {
		return
	}

	// This operation might be kind of costly as each termination must be done sequentially
	// if too many listeners are terminated at once it might take a while
	for _, listenerID := range terminateQueue {
		if err := m.terminateListener(context.Background(), listenerID); err != nil {
			log.Printf("failed to terminate listener, might've been already terminated: %s", listenerID)
			return
		}
	}
}

// Adding a listener to a listenerController does not start it.
func (m *Manager) addListener(listenerConfig ListenerConfig) error {
	// Create listener object before acquiring any locks
	listener, err := NewListenerFromConfig(listenerConfig)
	if err != nil {
		return fmt.Errorf("failed to create listener config: %w", err)
	}

	m.listenersMux.Lock()
	defer m.listenersMux.Unlock()

	_, ok := m.listeners[listenerConfig.ID]
	if ok {
		return fmt.Errorf("listener with ID %s already exists", listenerConfig.ID)
	}

	// Local listener specific configurations
	// if listener.Config.Deployment == DeploymentLocal {
	// 	// Validate port selection
	// 	m.portsMux.Lock()
	// 	for _, port := range m.usedPorts {
	// 		if listenerConfig.Port == port {
	// 			return fmt.Errorf("selected port already in use")
	// 		}
	// 	}
	// 	m.usedPorts = append(m.usedPorts, listenerConfig.Port)
	// 	m.portsMux.Unlock()
	// }

	m.listeners[listener.Config.ID] = listener

	m.syncMux.Lock()
	m.synchronizationLog[listener.Config.ID] = time.Now()
	m.syncMux.Unlock()

	return nil
}

// Starts an already existing listener
func (m *Manager) startListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	m.listenersMux.RLock()
	defer m.listenersMux.RUnlock()

	listener, ok := m.listeners[listenerID]
	if !ok {
		return fmt.Errorf("listener with ID %s not found", listenerID)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(m.startTimeout)*time.Second)
		defer cancel()
		if err := listener.Start(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start listener: %w", err)
		}
		close(errChan)
	}()

	return nil
}

// Stops a running listener
func (m *Manager) stopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	m.listenersMux.RLock()

	listener, ok := m.listeners[listenerID]
	if !ok {
		m.listenersMux.RUnlock()
		return fmt.Errorf("listener with ID %s not found", listenerID)
	}
	m.listenersMux.RUnlock()

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(m.stopTimeout)*time.Second)
		defer cancel()
		if err := listener.Stop(ctx); err != nil {
			errChan <- fmt.Errorf("failed to stop listener: %w", err)
		}
		close(errChan)
	}()

	return nil
}

// Terminate a listener and remove from registry
func (m *Manager) terminateListener(ctx context.Context, listenerID string) error {
	m.listenersMux.RLock()
	listener, ok := m.listeners[listenerID]
	if !ok {
		return fmt.Errorf("listener with ID '%s' not found", listenerID)
	}
	m.listenersMux.RUnlock()

	ctx, cancel := context.WithTimeout(ctx, time.Duration(m.stopTimeout)*time.Second)
	defer cancel()
	if err := listener.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	m.listenersMux.Lock()
	delete(m.listeners, listenerID)
	m.listenersMux.Unlock()

	return nil
}

// Updates a listener's config
func (m *Manager) updateListener(ctx context.Context, listenerConfig ListenerConfig) error {
	m.listenersMux.RLock()
	defer m.listenersMux.RUnlock()

	// Search the listener from the provided ID in the config
	listener, ok := m.listeners[listenerConfig.ID]
	if !ok {
		return fmt.Errorf("listener with ID '%s' not found", listenerConfig.ID)
	}

	// Update listener config, the listener itself controls its own relaunch, if necessary
	if err := listener.UpdateConfig(ctx, listenerConfig); err != nil {
		return fmt.Errorf("failed to update listener config: %w", err)
	}

	return nil
}

// External listener only
// Update a listener status
func (m *Manager) updateListenerStatus(ctx context.Context, listenerID, status string) error {
	m.listenersMux.RLock()
	defer m.listenersMux.RUnlock()

	listener, ok := m.listeners[listenerID]
	if !ok {
		return fmt.Errorf("listener not found: %s", listenerID)
	}

	if listener.Config.Kind != DeploymentExternal {
		return errors.New("operation not allowed for local listeners")
	}

	listener.UpdateStatus(ctx, status)

	return nil
}

// External listener only
// Synchronizes configurations with external listeners
func (m *Manager) synchronize(listenerID string) (ListenerConfig, error) {
	m.listenersMux.RLock()
	defer m.listenersMux.RUnlock()

	listener, ok := m.listeners[listenerID]
	if !ok {
		return ListenerConfig{}, fmt.Errorf("listener not found: %s", listenerID)
	}

	// Update last synchronization time
	m.syncMux.Lock()
	m.synchronizationLog[listenerID] = time.Now()
	m.syncMux.Unlock()

	return listener.Config, nil
}

func (m *Manager) getListenerStatuses() map[string]string {
	m.listenersMux.RLock()
	listenerStatuses := make(map[string]string)
	for id, list := range m.listeners {
		listenerStatuses[id] = list.Config.Status
	}
	m.listenersMux.RUnlock()
	return listenerStatuses
}
