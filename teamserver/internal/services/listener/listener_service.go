package listener

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

type statusUpdate struct {
	listenerID string
	status     string
}

// ListenerService is the entrypoint for listener operations exposed to clients
type ListenerService struct {
	startTimeout int
	stopTimeout  int

	listenerDal dal.IListenerDAL

	// Keep track of the runtime listener representations
	activeListeners map[string]*Listener
	statusUpdates   chan statusUpdate
	mux             sync.RWMutex
}

func NewListenerService(listenerDAL dal.IListenerDAL) *ListenerService {
	ls := &ListenerService{
		startTimeout:    15,
		stopTimeout:     15,
		listenerDal:     listenerDAL,
		activeListeners: make(map[string]*Listener),
		statusUpdates:   make(chan statusUpdate, 100),
	}

	// Start the status update processor
	go ls.processStatusUpdates()

	return ls
}

func (ls *ListenerService) processStatusUpdates() {
	for update := range ls.statusUpdates {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		updates := map[string]any{"status": update.status}

		if err := ls.listenerDal.UpdateListener(ctx, update.listenerID, updates); err != nil {
			logger.Error(fmt.Sprintf("Failed to update listener status: %v", err))
		}

		cancel()
	}
}

func (ls *ListenerService) monitorListenerStatus(listener *Listener) {
	for status := range listener.statusUpdatesCh {
		ls.statusUpdates <- statusUpdate{
			listenerID: listener.ID,
			status:     status,
		}
	}
}

func (ls *ListenerService) StartListener(ctx context.Context, listenerID string) error {
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	// if exists {
	// if listener.Status == StatusStarting {
	// 	return fmt.Errorf("listener with ID %s is starting", listenerID)
	// } else if listener.Status == StatusRunning {
	// 	return fmt.Errorf("listener with ID %s is already running", listenerID)
	// }
	// } else {
	if !exists {
		listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
		if err != nil {
			return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
		}

		listener, err = createListenerFromModel(listenerModel)
		if err != nil {
			return fmt.Errorf("failed to create listener from model: %w", err)
		}

		go ls.monitorListenerStatus(listener)
	}

	return ls.startListener(ctx, listener)
}

func (ls *ListenerService) startListener(ctx context.Context, listener *Listener) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.startTimeout)*time.Second)
	defer cancel()

	startCh := make(chan error, 1)
	go func() {
		startCh <- listener.Start(ctx)
		close(startCh)
	}()

	select {
	case err := <-startCh:
		if err != nil {
			return err
		}
		ls.mux.Lock()
		ls.activeListeners[listener.ID] = listener
		ls.mux.Unlock()
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

// The expected behavior, and since all listeners that are running should be kept track of, return an error if the listener is not mapped
// Do not try retrieving from the database as the runtime struct of the listener is required to stop it, otherwise, you are simply creating a listener just to stop it
func (ls *ListenerService) StopListener(ctx context.Context, listenerID string) error {
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	// if listener.Status == StatusStopping {
	// 	return fmt.Errorf("listener with ID %s is already stopping", listenerID)
	// } else if listener.Status == StatusReady {
	// 	return fmt.Errorf("listener with ID %s is already stopped", listenerID)
	// }

	if !exists {
		return fmt.Errorf("trying to stop listener that is not mapped")
	}

	return ls.stopListener(ctx, listener)
}

func (ls *ListenerService) stopListener(ctx context.Context, listener *Listener) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
	defer cancel()

	stopCh := make(chan error, 1)
	go func() {
		stopCh <- listener.Stop(ctx)
		close(stopCh)
	}()

	select {
	case err := <-stopCh:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Terminate is an operation that fully stops a listener, killing processes and removing from active listeners map
func (ls *ListenerService) TerminateListener(ctx context.Context, listenerID string) error {
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerID]
	ls.mux.RUnlock()

	if !exists {
		return fmt.Errorf("trying to terminate listener that is not mapped")
	}

	return ls.terminateListener(ctx, listener)
}

func (ls *ListenerService) terminateListener(ctx context.Context, listener *Listener) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
	defer cancel()

	terminateCh := make(chan error, 1)
	go func() {
		terminateCh <- listener.Terminate(ctx)
		close(terminateCh)
	}()

	select {
	case err := <-terminateCh:
		if err != nil {
			return err
		}
		ls.mux.Lock()
		delete(ls.activeListeners, listener.ID)
		ls.mux.Unlock()
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (ls *ListenerService) UpdateListener(ctx context.Context, listenerModel models.Listener) error {
	// First check if we have it in our active map
	ls.mux.RLock()
	listener, exists := ls.activeListeners[listenerModel.ID]
	ls.mux.RUnlock()

	// If not found in active listeners, create from model
	if exists {
		if err := listener.UpdateConfig(ctx, listenerModel); err != nil {
			return fmt.Errorf("failed to update listener config: %w", err)
		}
	}

	// Save updated listener to DAL
	updates := map[string]any{
		"config": listener.RawConfig,
	}
	return ls.listenerDal.UpdateListener(ctx, listener.ID, updates)
}

func (ls *ListenerService) AutoStart(ctx context.Context) error {
	listeners, err := ls.listenerDal.GetActiveListeners(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active listeners: %w", err)
	}

	for _, listener := range listeners {
		// Create listener instance first

		listenerInstance, err := createListenerFromModel(listener)
		if err != nil {
			return fmt.Errorf("failed to create listener instance: %w", err)
		}

		// Start monitoring status updates for all listeners
		go ls.monitorListenerStatus(listenerInstance)

		// Add to active listeners map regardless of status
		ls.mux.Lock()
		ls.activeListeners[listener.ID] = listenerInstance
		ls.mux.Unlock()

		if listener.Status == StatusRunning {

			// Set status as ready, because otherwise listener won't be ready to start
			// another option could be to clean up the listeners when the server is force shutdown
			// but what if the server dies. Listeners will cease to work simply because there is no
			// full-proof feature to get them up and running. We know listeners die on shutdown, so let
			// us just set as ready.
			listenerInstance.Status = StatusReady
			if err := ls.startListener(ctx, listenerInstance); err != nil {
				logger.Error(fmt.Sprintf("Failed to start listener %s during AutoStart: %v", listener.ID, err))

				updates := map[string]any{"status": StatusFailed}
				if updateErr := ls.listenerDal.UpdateListener(ctx, listener.ID, updates); updateErr != nil {
					logger.Error(fmt.Sprintf("Failed to update listener status: %v", updateErr))
				}
			}
		}
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

// TODO: I think this will not happen through here but from the external package
// it should start up a server to receive requests from the external listeners, whenever there are any
// func (ls *ListenerService) UpdateListenerStatus(ctx context.Context, listenerID, status string) error {
// 	// First check if we have it in our active map
// 	ls.mux.RLock()
// 	listener, exists := ls.activeListeners[listenerID]
// 	ls.mux.RUnlock()

// 	// If not found in active listeners, try loading from database
// 	if !exists {
// 		listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
// 		if err != nil {
// 			return fmt.Errorf("listener not found: %w", err)
// 		}

// 		listener, err = createListenerFromModel(listenerModel)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	listener.UpdateStatus(ctx, status)

// 	// Update the in-memory copy if it exists
// 	if exists {
// 		ls.mux.Lock()
// 		ls.activeListeners[listenerID] = listener
// 		ls.mux.Unlock()
// 	}

// 	// Update in DAL
// 	updates := map[string]any{"status": status}
// 	return ls.listenerDal.UpdateListener(ctx, listenerID, updates)
// }

// func (ls *ListenerService) synchronize(ctx context.Context, listenerID string) (*Listener, error) {
// 	// First check if we have it in our active map
// 	ls.mux.RLock()
// 	listener, exists := ls.activeListeners[listenerID]
// 	ls.mux.RUnlock()

// 	if exists {
// 		return listener, nil
// 	}

// 	// If not in active map, retrieve from database
// 	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
// 	if err != nil {
// 		return nil, fmt.Errorf("listener not found: %w", err)
// 	}

// 	listener, err = createListenerFromModel(listenerModel)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return listener, nil
// }
