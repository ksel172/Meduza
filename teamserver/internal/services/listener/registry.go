package listener

import (
	"sync"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
)

type ListenerRegistry struct {
	listeners map[string]*Listener
	mux       sync.RWMutex
}

func NewListenerRegistry(listenerDal dal.IListenerDAL) *ListenerRegistry {
	return &ListenerRegistry{listeners: make(map[string]*Listener)}
}

func (lr *ListenerRegistry) Get(listenerID string) (*Listener, bool) {
	lr.mux.RLock()
	listener, exists := lr.listeners[listenerID]
	lr.mux.RUnlock()

	if !exists {
		return nil, false
	}

	return listener, true
}

func (lr *ListenerRegistry) Register(listener *Listener) {
	lr.mux.Lock()
	lr.listeners[listener.ID] = listener
	lr.mux.Unlock()
}

func (lr *ListenerRegistry) Delete(listenerID string) {
	lr.mux.Lock()
	delete(lr.listeners, listenerID)
	lr.mux.Unlock()
}
