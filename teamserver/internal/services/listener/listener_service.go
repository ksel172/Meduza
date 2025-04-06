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
	activeListeners map[string]Listener
	mux             sync.RWMutex
}

func NewListenerService(listenerDAL dal.IListenerDAL) *ListenerService {
	return &ListenerService{
		startTimeout:    15,
		stopTimeout:     15,
		listenerDal:     listenerDAL,
		activeListeners: make(map[string]Listener),
	}
}

func (ls *ListenerService) StartListener(ctx context.Context, listenerID string, errChan chan<- error) error {
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
	go ls.doListenerAction(ctx, listener.Start, listenerID, errChan)

	return nil
}

// goroutine helper function
func (ls *ListenerService) doListenerAction(ctx context.Context, fn ListenerActionFunc, listenerID string, errChan chan<- error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.startTimeout)*time.Second)
	defer cancel()

	if err := fn(ctx); err != nil {
		errChan <- fmt.Errorf("failed listener operation: %w", err)
		return
	}

	// Update listener status in DAL
	updates := map[string]any{"status": "running"}
	if err := ls.listenerDal.UpdateListener(ctx, listenerID, updates); err != nil {
		errChan <- fmt.Errorf("failed to update listener status: %w", err)
		return
	}

	close(errChan)
}

func (ls *ListenerService) StopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	// Listener needs to be setup once data fields are read from storage
	l, err := createListenerFromModel(listenerModel)
	if err != nil {
		return fmt.Errorf("failed to create listener from model: %w", err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
		defer cancel()
		if err := l.Stop(ctx); err != nil {
			errChan <- fmt.Errorf("failed to stop listener: %w", err)
			return
		}

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
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID '%s' not found: %w", listenerID, err)
	}

	l, err := createListenerFromModel(listenerModel)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
	defer cancel()
	if err := l.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to close listener: %w", err)
	}

	// Remove from DAL
	if err := ls.listenerDal.DeleteListener(ctx, listenerID); err != nil {
		return fmt.Errorf("failed to delete listener from storage: %w", err)
	}

	return nil
}

func (ls *ListenerService) UpdateListener(ctx context.Context, listenerModel models.Listener) error {
	listener, err := createListenerFromModel(listenerModel)
	if err != nil {
		return err
	}

	// Update listener config
	if err := listener.UpdateConfig(ctx, listener); err != nil {
		return fmt.Errorf("failed to update listener config: %w", err)
	}

	// Save updated listener to DAL
	updates := map[string]any{
		"config": listener.RawConfig,
	}
	return ls.listenerDal.UpdateListener(ctx, listener.ID, updates)
}

func (ls *ListenerService) UpdateListenerStatus(ctx context.Context, listenerID, status string) error {
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener not found: %w", err)
	}

	l, err := createListenerFromModel(listenerModel)
	if err != nil {
		return err
	}

	l.UpdateStatus(ctx, status)

	// Update in DAL
	updates := map[string]any{"status": status}
	return ls.listenerDal.UpdateListener(ctx, listenerID, updates)
}

func (ls *ListenerService) synchronize(ctx context.Context, listenerID string) (*Listener, error) {
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return nil, fmt.Errorf("listener not found: %w", err)
	}

	l, err := createListenerFromModel(listenerModel)
	if err != nil {
		return nil, err
	}

	return l, nil
}

// This will not work as there is no reference to the actual listener implementations currently
// The service must keep track of all listener internally for functionality to exist
func (ls *ListenerService) AutoStart(ctx context.Context) error {
	listeners, err := ls.listenerDal.GetActiveListeners(ctx)
	if err != nil {
		return fmt.Errorf("failed to get active listeners: %w", err)
	}

	for _, listener := range listeners {
		if listener.Status == StatusRunning {
			continue
		}

		if err := ls.StartListener(ctx, listener.ID, make(chan<- error)); err != nil {
			return fmt.Errorf("failed to start listener: %w", err)
		}
	}

	return nil
}
