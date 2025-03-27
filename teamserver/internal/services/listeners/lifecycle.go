package services

import (
	"context"

	"github.com/ksel172/Meduza/teamserver/utils"
)

type ListenerLifecycleManager interface {
	Start(ctx context.Context, listener *Listener) error
	Stop(ctx context.Context, listener *Listener) error
	Terminate(ctx context.Context, listener *Listener) error
	UpdateConfig(ctx context.Context, listener *Listener, config any) error
}

// Implementation for managed listeners
type ManagedLifecycleManager struct{}

func NewManagedLifecycleManager() *ManagedLifecycleManager {
	return &ManagedLifecycleManager{}
}

// Implementation for scheduled listeners
type ScheduledLifecycleManager struct{}

func NewScheduledLifecycleManager() *ScheduledLifecycleManager {
	return &ScheduledLifecycleManager{}
}

func (m *ManagedLifecycleManager) Start(ctx context.Context, l *Listener) error {
	utils.AssertNotNil(l.listener)

	l.Status = StatusStarting
	if err := l.listener.Start(ctx); err != nil {
		return err
	}
	l.Status = StatusRunning

	return nil
}

func (m *ManagedLifecycleManager) Stop(ctx context.Context, l *Listener) error {
	utils.AssertNotNil(l.listener)

	l.Status = StatusStopping
	if err := l.listener.Stop(ctx); err != nil {
		return err
	}
	l.Status = StatusReady

	return nil
}

func (m *ManagedLifecycleManager) Terminate(ctx context.Context, l *Listener) error {
	utils.AssertNotNil(l.listener)
	l.Status = StatusTerminating
	return l.listener.Terminate(ctx)
}

func (m *ManagedLifecycleManager) UpdateConfig(ctx context.Context, l *Listener, config any) error {
	utils.AssertNotNil(l.listener)

	// This could introduce some bugs that we need to handle later
	// For example, what if the listener actually updates its config but returns an error anyway for some reason?
	// Or what if the return from this request is missed and it times out, the listener config will not be synced
	// with how the controller. Some polling mechanism or reconciliation loop could fix this.
	if err := l.listener.UpdateConfig(ctx); err != nil {
		return err
	}
	l.Config = config

	return nil
}

func (m *ScheduledLifecycleManager) Start(ctx context.Context, l *Listener) error {
	l.Status = StatusStarting
	return nil
}

func (m *ScheduledLifecycleManager) Stop(ctx context.Context, l *Listener) error {
	l.Status = StatusStopping
	return nil
}

func (m *ScheduledLifecycleManager) Terminate(ctx context.Context, l *Listener) error {
	l.Status = StatusTerminating
	return nil
}

// We just need to update the config and the listener will poll the controller for updates to its config
func (m *ScheduledLifecycleManager) UpdateConfig(ctx context.Context, l *Listener, config any) error {
	l.Config = config
	return nil
}
