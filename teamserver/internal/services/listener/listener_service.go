package listener

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
)

/*
TODO: restructure listener service
    1. Think about the operation kinds
        a. Action operations
        b. Data operations
    2. Should a controller control action operations while the service provides data operations?
    3. Synchronization is done through database, previously was a listener map
    4. Maybe move all data operations to handler and rename to service to controller?
*/

// ListenerService is the entrypoint for listener operations exposed to clients
type ListenerService struct {
	startTimeout int
	stopTimeout  int

	listenerDal dal.IListenerDAL

	// Keep track of the runtime listener representations
	activeListeners map[string]*Listener
	mux             sync.RWMutex
}

func NewListenerService(listenerDAL dal.IListenerDAL) *ListenerService {
	return &ListenerService{
		startTimeout:    15,
		stopTimeout:     15,
		listenerDal:     listenerDAL,
		activeListeners: make(map[string]*Listener),
	}
}

func (ls *ListenerService) StartListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	// Check if the listener is already active
	ls.mux.RLock()
	_, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	if exists {
		return fmt.Errorf("listener with ID %s is already running", listenerID)
	}

	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	// Listener needs to be setup once data fields are read from storage
	listener, err := createListenerFromModel(listenerModel)
	if err != nil {
		return fmt.Errorf("failed to create listener from model: %w", err)
	}

	// Start listener in goroutine
	go ls.doListenerAction(ctx, listener.Start, listenerID, errChan, listener)

	return nil
}

// goroutine helper function
func (ls *ListenerService) doListenerAction(ctx context.Context,
	fn ListenerActionFunc,
	listenerID string,
	errChan chan<- error,
	listener *Listener) {

	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.startTimeout)*time.Second)
	defer cancel()

	if err := fn(ctx); err != nil {
		errChan <- fmt.Errorf("failed listener operation: %w", err)
		return
	}

	// Add to active listeners map
	ls.mux.Lock()
	ls.activeListeners[listenerID] = listener
	ls.mux.Unlock()

	// Update listener status in DAL
	updates := map[string]any{"status": "running"}
	if err := ls.listenerDal.UpdateListener(ctx, listenerID, updates); err != nil {
		errChan <- fmt.Errorf("failed to update listener status: %w", err)
		return
	}

	close(errChan)
}

func (ls *ListenerService) StopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	// First check if we have it in our active map
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	// If not found in active listeners, try loading from database
	if !exists {
		listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
		if err != nil {
			return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
		}

		// Listener needs to be setup once data fields are read from storage
		listener, err = createListenerFromModel(listenerModel)
		if err != nil {
			return fmt.Errorf("failed to create listener from model: %w", err)
		}
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
		defer cancel()
		if err := listener.Stop(ctx); err != nil {
			errChan <- fmt.Errorf("failed to stop listener: %w", err)
			return
		}

		// Remove from active listeners map
		ls.mux.Lock()
		delete(ls.activeListeners, listenerID)
		ls.mux.Unlock()

		// Update listener status in DAL
		updates := map[string]any{"status": "stopped"}
		if err := ls.listenerDal.UpdateListener(ctx, listenerID, updates); err != nil {
			errChan <- fmt.Errorf("failed to update listener status: %w", err)
			return
		}
		close(errChan)
	}()

	return nil
}

func (ls *ListenerService) TerminateListener(ctx context.Context, listenerID string) error {
	// First check if we have it in our active map
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	// If not found in active listeners, try loading from database
	if !exists {
		listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
		if err != nil {
			return fmt.Errorf("listener with ID '%s' not found: %w", listenerID, err)
		}

		listener, err = createListenerFromModel(listenerModel)
		if err != nil {
			return err
		}
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
	defer cancel()
	if err := listener.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	// Remove from active listeners map
	ls.mux.Lock()
	delete(ls.activeListeners, listenerID)
	ls.mux.Unlock()

	// Remove from DAL
	if err := ls.listenerDal.DeleteListener(ctx, listenerID); err != nil {
		return fmt.Errorf("failed to delete listener from storage: %w", err)
	}

	return nil
}

func (ls *ListenerService) UpdateListener(ctx context.Context, listenerModel models.Listener) error {
	// First check if we have it in our active map
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerModel.ID]
	ls.mux.RUnlock()

	// If not found in active listeners, create from model
	if !exists {
		var err error
		listener, err = createListenerFromModel(listenerModel)
		if err != nil {
			return err
		}
	}

	if err := listener.UpdateConfig(ctx, listener); err != nil {
		return fmt.Errorf("failed to update listener config: %w", err)
	}

	if exists {
		ls.mux.Lock()
		ls.activeListeners[listenerModel.ID] = listener
		ls.mux.Unlock()
	}

	// Save updated listener to DAL
	updates := map[string]any{
		"config": listener.RawConfig,
	}
	return ls.listenerDal.UpdateListener(ctx, listener.ID, updates)
}

func (ls *ListenerService) UpdateListenerStatus(ctx context.Context, listenerID, status string) error {
	// First check if we have it in our active map
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	// If not found in active listeners, try loading from database
	if !exists {
		listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
		if err != nil {
			return fmt.Errorf("listener not found: %w", err)
		}

		listener, err = createListenerFromModel(listenerModel)
		if err != nil {
			return err
		}
	}

	listener.UpdateStatus(ctx, status)

	// Update the in-memory copy if it exists
	if exists {
		ls.mux.Lock()
		ls.activeListeners[listenerID] = listener
		ls.mux.Unlock()
	}

	// Update in DAL
	updates := map[string]any{"status": status}
	return ls.listenerDal.UpdateListener(ctx, listenerID, updates)
}

func (ls *ListenerService) synchronize(ctx context.Context, listenerID string) (*Listener, error) {
	// First check if we have it in our active map
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	if exists {
		return listener, nil
	}

	// If not in active map, retrieve from database
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return nil, fmt.Errorf("listener not found: %w", err)
	}

	listener, err = createListenerFromModel(listenerModel)
	if err != nil {
		return nil, err
	}

	return listener, nil
}

func (ls *ListenerService) AutoStart(ctx context.Context) error {
	listeners, err := ls.listenerDal.GetActiveListeners(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active listeners: %w", err)
	}

	for _, listener := range listeners {
		if listener.Status == StatusRunning {
			// Create instance and add to active listeners map
			listenerInstance, err := createListenerFromModel(listener)
			if err != nil {
				return fmt.Errorf("failed to create listener instance: %w", err)
			}

			ls.mux.Lock()
			ls.activeListeners[listener.ID] = listenerInstance
			ls.mux.Unlock()
			continue
		}

		errChan := make(chan error, 1)
		if err := ls.StartListener(ctx, listener.ID, errChan); err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}

		// Optionally wait for errors from the start operation
		// err = <-errChan
		// if err != nil {
		//     return fmt.Errorf("error during listener start: %w", err)
		// }
	}

	return nil
}

// GetActiveListeners returns a copy of the current active listeners map
func (ls *ListenerService) GetActiveListeners() map[string]*Listener {
	ls.mux.RLock()
	defer ls.mux.RUnlock()

	// Create a copy to avoid access issues
	result := make(map[string]*Listener, len(ls.activeListeners))
	for id, listener := range ls.activeListeners {
		result[id] = listener
	}

	return result
}
