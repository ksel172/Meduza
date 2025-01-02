package listeners

import (
	"sync"
	"time"
)

// Interface for ListenerRegistry
type ListenerRegistry interface {
	Start() error
	Stop(timeout time.Duration) error
}

// Registry holds active listeners by ID
type Registry struct {
	mu        sync.Mutex
	listeners map[string]ListenerRegistry
}

// Create a New Instance of ListenerRegistry.
// Returns map of ListenerRegistry :- map[string]ListenerRegistry
func NewRegistry() *Registry {
	return &Registry{
		listeners: make(map[string]ListenerRegistry),
	}
}

// AddListener function adds a Listener to registry.
func (r *Registry) AddListener(id string, listener ListenerRegistry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners[id] = listener
}

// GetListener function help to get the added listeners
// Help in checking if the listener is currently running or not.
func (r *Registry) GetListener(id string) (ListenerRegistry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	listener, exists := r.listeners[id]
	return listener, exists
}

// RemoveListener function remove the current running listeners
func (r *Registry) RemoveListener(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.listeners, id)
}
