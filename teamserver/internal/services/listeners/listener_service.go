package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

type ListenersService struct {
	startTimeout int
	stopTimeout  int

	listenerDal IListenerDAL

	// Track synchronization timestamps for external listeners
	synchronizationLog map[string]time.Time
	syncMux            sync.Mutex
}

func NewListenerService(listenerDAL IListenerDAL) *ListenersService {
	return &ListenersService{
		startTimeout:       15,
		stopTimeout:        15,
		listenerDal:        listenerDAL,
		synchronizationLog: make(map[string]time.Time),
	}
}

func (ls *ListenersService) GetListeners(ctx context.Context) ([]*Listener, error) {
	listeners, err := ls.listenerDal.GetAllListeners(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get listeners: %w", err)
	}

	return listeners, nil
}

func (ls *ListenersService) GetListener(ctx context.Context, listenerID string) (*Listener, error) {
	listener, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return nil, fmt.Errorf("failed to get listener: %w", err)
	}

	return listener, nil
}

func (ls *ListenersService) GetListenerStatuses(ctx context.Context) (map[string]string, error) {
	listeners, err := ls.listenerDal.GetAllListeners(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get listeners: %w", err)
	}

	listenerStatuses := make(map[string]string)
	for _, listener := range listeners {
		listenerStatuses[listener.ID] = listener.Status
	}
	return listenerStatuses, nil
}

func (ls *ListenersService) AddListener(ctx context.Context, listener *Listener) error {
	// Check if listener with same ID already exists
	// _, err := ls.listenerDal.GetListenerByID(ctx, listener.ID)
	// if err == nil {
	// 	return fmt.Errorf("listener with ID %s already exists", listener.ID)
	// }

	// Create listener object
	newListener, err := NewListenerFromBase(listener)
	if err != nil {
		return fmt.Errorf("failed to create listener config: %w", err)
	}

	// Add to DAL
	if err := ls.listenerDal.CreateListener(ctx, newListener); err != nil {
		return fmt.Errorf("failed to store listener: %w", err)
	}

	// Initialize synchronization record for this listener
	ls.syncMux.Lock()
	ls.synchronizationLog[listener.ID] = time.Now()
	ls.syncMux.Unlock()

	return nil
}

func (ls *ListenersService) StartListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listener, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

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

	return nil
}

func (ls *ListenersService) StopListener(ctx context.Context, listenerID string, errChan chan<- error) error {
	listener, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID %s not found: %w", listenerID, err)
	}

	go func() {
		ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
		defer cancel()
		if err := listener.Stop(ctx); err != nil {
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

func (ls *ListenersService) TerminateListener(ctx context.Context, listenerID string) error {
	listener, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener with ID '%s' not found: %w", listenerID, err)
	}

	ctx, cancel := context.WithTimeout(ctx, time.Duration(ls.stopTimeout)*time.Second)
	defer cancel()
	if err := listener.Terminate(ctx); err != nil {
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

func (ls *ListenersService) UpdateListener(ctx context.Context, listener *Listener) error {
	existingListener, err := ls.listenerDal.GetListenerById(ctx, listener.ID)
	if err != nil {
		return fmt.Errorf("listener with ID '%s' not found: %w", listener.ID, err)
	}

	// Update listener config
	if err := existingListener.UpdateConfig(ctx, listener); err != nil {
		return fmt.Errorf("failed to update listener config: %w", err)
	}

	// Save updated listener to DAL
	updates := map[string]any{
		"config": listener.Config,
	}
	return ls.listenerDal.UpdateListener(ctx, listener.ID, updates)
}

func (ls *ListenersService) UpdateListenerStatus(ctx context.Context, listenerID, status string) error {
	listener, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return fmt.Errorf("listener not found: %w", err)
	}

	if listener.Deployment != DeploymentExternal {
		return errors.New("operation not allowed for local listeners")
	}

	listener.UpdateStatus(ctx, status)

	// Update in DAL
	updates := map[string]any{"status": status}
	return ls.listenerDal.UpdateListener(ctx, listenerID, updates)
}

func (ls *ListenersService) synchronize(ctx context.Context, listenerID string) (*Listener, error) {
	listener, err := ls.listenerDal.GetListenerById(ctx, listenerID)
	if err != nil {
		return nil, fmt.Errorf("listener not found: %w", err)
	}

	// Update last synchronization time
	ls.syncMux.Lock()
	ls.synchronizationLog[listenerID] = time.Now()
	ls.syncMux.Unlock()

	return listener, nil
}

func (ls *ListenersService) GetListenerByName(ctx context.Context, name string) (*Listener, error) {
	listener, err := ls.listenerDal.GetListenerByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to get listener: %w", err)
	}

	return listener, nil
}

func (ls *ListenersService) AutoStart(ctx context.Context) error {
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
