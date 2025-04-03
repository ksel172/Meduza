package listener

import (
	"context"
	"errors"
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

	// Track synchronization timestamps for external listeners
	synchronizationLog map[string]time.Time
	syncMux            sync.Mutex
}

func NewListenerService(listenerDAL dal.IListenerDAL) *ListenerService {
	return &ListenerService{
		startTimeout:       15,
		stopTimeout:        15,
		listenerDal:        listenerDAL,
		synchronizationLog: make(map[string]time.Time),
	}
}

func (ls *ListenerService) StartListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	// Listener needs to be setup once data fields are read from storage
	listener, err := CreateListenerFromModel(listenerModel)
	if err != nil {
		return fmt.Errorf("failed to create listener from model: %w", err)
	}

	// Run listener in another goroutine
	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.startTimeout)*time.Second)
		defer cancel()
		if err := listener.Start(ctx); err != nil {
			errChan <- fmt.Errorf("failed to start listener: %w", err)
			return
		}

		// Update listener status in DAL
		updates := map[string]any{"status": "running"}
		if err := ls.listenerDal.UpdateListener(ctx, listenerID, updates); err != nil {
			errChan <- fmt.Errorf("failed to update listener status: %w", err)
			return
		}
		close(errChan)
	}()

	// Initialize synchronization record for this listener
	ls.syncMux.Lock()
	ls.synchronizationLog[listener.ID] = time.Now()
	ls.syncMux.Unlock()

	return nil
}

func (ls *ListenerService) StopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	// Listener needs to be setup once data fields are read from storage
	l := CreateListenerFromModel(listenerModel)

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

	l, err := CreateListenerFromModel(listenerModel)
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

	// Clean up synchronization record
	ls.syncMux.Lock()
	delete(ls.synchronizationLog, listenerID)
	ls.syncMux.Unlock()

	return nil
}

func (ls *ListenerService) UpdateListener(ctx context.Context, listenerModel models.Listener) error {
	listener, err := CreateListenerFromModel(listenerModel)
	if err != nil {
		return err
	}

	// Update listener config
	if err := listener.UpdateConfig(ctx, listener); err != nil {
		return fmt.Errorf("failed to update listener config: %w", err)
	}

	// Save updated listener to DAL
	updates := map[string]any{
		"config": listener.Config,
	}
	return ls.listenerDal.UpdateListener(ctx, listener.ID, updates)
}

func (ls *ListenerService) UpdateListenerStatus(ctx context.Context, listenerID, status string) error {
	listenerModel, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener not found: %w", err)
	}

	l, err := CreateListenerFromModel(listenerModel)
	if err != nil {
		return err
	}

	if l.Deployment != DeploymentExternal {
		return errors.New("operation not allowed for local listeners")
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

	l, err := CreateListenerFromModel(listenerModel)
	if err != nil {
		return err
	}

	// Update last synchronization time
	ls.syncMux.Lock()
	ls.synchronizationLog[listenerID] = time.Now()
	ls.syncMux.Unlock()

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
