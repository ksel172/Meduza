package listeners

import (
	"sync"

	"github.com/ksel172/Meduza/teamserver/models"
)

// Registry holds active listeners by ID
type Registry struct {
	mu        sync.Mutex
	listeners map[string]models.Listener
}

// Create a New Instance of ListenerRegistry.
// Returns map of ListenerRegistry :- map[string]ListenerRegistry
func NewRegistry() *Registry {
	return &Registry{
		listeners: make(map[string]models.Listener),
	}
}

// AddListener function adds a Listener to registry.
func (r *Registry) AddListener(listener models.Listener) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners[listener.ID.String()] = listener
}

// GetListener function help to get the added listeners
// Help in checking if the listener is currently running or not.
func (r *Registry) GetListener(id string) (models.Listener, bool) {
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
