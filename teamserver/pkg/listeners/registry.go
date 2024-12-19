package listeners

import (
	"sync"
	"time"
)

type ListenerRegistry interface {
	Start() error
	Stop(timeout time.Duration) error
}

// Registry holds active listeners by ID
type Registry struct {
	mu        sync.Mutex
	listeners map[string]ListenerRegistry
}

func NewRegistry() *Registry {
	return &Registry{
		listeners: make(map[string]ListenerRegistry),
	}
}

func (r *Registry) AddListener(id string, listener ListenerRegistry) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.listeners[id] = listener
}

func (r *Registry) GetListener(id string) (ListenerRegistry, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	listener, exists := r.listeners[id]
	return listener, exists
}

func (r *Registry) RemoveListener(id string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.listeners, id)
}
