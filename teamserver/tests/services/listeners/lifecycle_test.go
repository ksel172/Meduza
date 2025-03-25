package service_tests

import (
	"context"
	"errors"
	"testing"

	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Lifecycle manager mock
type MockLifecycleManager struct {
	mock.Mock
}

func (m *MockLifecycleManager) Start(ctx context.Context, listener *services.Listener) error {
	args := m.Called(ctx, listener)
	return args.Error(0)
}

func (m *MockLifecycleManager) Stop(ctx context.Context, listener *Listener) error {
	args := m.Called(ctx, listener)
	return args.Error(0)
}

func (m *MockLifecycleManager) Terminate(ctx context.Context, listener *Listener) error {
	args := m.Called(ctx, listener)
	return args.Error(0)
}

func (m *MockLifecycleManager) UpdateConfig(ctx context.Context, listener *Listener, config ListenerConfig) error {
	args := m.Called(ctx, listener)
	return args.Error(0)
}

// Listener implementation mock
type MockListenerImplementation struct {
	mock.Mock
}

func (m *MockListenerImplementation) Start(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockListenerImplementation) Stop(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockListenerImplementation) Terminate(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockListenerImplementation) UpdateConfig(context.Context) error {
	args := m.Called()
	return args.Error(0)
}

func TestManagedLifecycle_Start(t *testing.T) {
	tests := []struct {
		name       string
		config     ListenerConfig
		endStatus  string
		startError error // error for internal listener implementation
	}{
		{
			name:       "successful-start",
			config:     makeValidListenerConfig("test-listener-id", "http", StatusReady),
			endStatus:  StatusRunning,
			startError: nil,
		},
		{
			name:       "start-error",
			config:     makeValidListenerConfig("test-listener-id", "http", StatusReady),
			endStatus:  StatusStarting,
			startError: errors.New("listener failed to start"),
		},
		{
			name:       "nil-listener",
			config:     makeValidListenerConfig("test-listener-id", "http", StatusReady),
			endStatus:  StatusReady,
			startError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.config.Deployment = LifecycleManaged
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			// Specific test for assertion of listener, doesn;t really fit in with the rest
			// maybe move somewhere else later
			if tc.name == "nil-listener" {
				listener.listener = nil
				assert.Panics(t, func() { listener.Start(context.Background()) }, "the function did not panic on a nil listener")
				return
			}

			// Make mock internal listener implementation
			mockListenerImplementation := &MockListenerImplementation{}
			mockListenerImplementation.On("Start", mock.Anything).Return(tc.startError).Once()
			listener.listener = mockListenerImplementation

			err = listener.lifecycleManager.Start(context.Background(), listener)
			if tc.startError != nil {
				assert.Error(t, err, "Expected an error on lifecycle manager start")
			} else {
				assert.NoError(t, err, "Expected no error on lifecycle manager start")
			}

			mockListenerImplementation.AssertExpectations(t)

			if listener.Config.Status != tc.endStatus {
				t.Errorf("expected config to match, expected: %s, actual: %s", tc.endStatus, listener.Config.Status)
			}
		})
	}
}

func TestScheduledLifecycle_Start(t *testing.T) {
	tests := []struct {
		name      string
		config    ListenerConfig
		endStatus string
	}{
		{
			name:      "successful-scheduled-start",
			config:    makeValidListenerConfig("scheduled-listener-id", "http", StatusReady),
			endStatus: StatusStarting,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create scheduled listener
			tc.config.Lifecycle = LifecycleScheduled
			tc.config.Deployment = DeploymentExternal
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			err = listener.lifecycleManager.Start(context.Background(), listener)

			assert.NoError(t, err, "Expected no error on scheduled lifecycle manager start")

			assert.Equal(t, tc.endStatus, listener.Config.Status,
				"expected status to be updated to %s, got %s", tc.endStatus, listener.Config.Status)
		})
	}
}

func TestManagedLifecycle_Stop(t *testing.T) {
	tests := []struct {
		name      string
		config    ListenerConfig
		endStatus string
		stopError error
	}{
		{
			name:      "successful-stop",
			config:    makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			endStatus: StatusReady,
			stopError: nil,
		},
		{
			name:      "stop-error",
			config:    makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			endStatus: StatusStopping,
			stopError: errors.New("listener failed to stop"),
		},
		{
			name:      "nil-listener",
			config:    makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			endStatus: StatusRunning,
			stopError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.config.Lifecycle = LifecycleManaged
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			if tc.name == "nil-listener" {
				listener.listener = nil
				assert.Panics(t, func() { listener.Stop(context.Background()) }, "the function did not panic on a nil listener")
				return
			}

			mockListenerImplementation := &MockListenerImplementation{}
			mockListenerImplementation.On("Stop", mock.Anything).Return(tc.stopError).Once()
			listener.listener = mockListenerImplementation

			err = listener.lifecycleManager.Stop(context.Background(), listener)
			if tc.stopError != nil {
				assert.Error(t, err, "Expected an error on lifecycle manager stop")
			} else {
				assert.NoError(t, err, "Expected no error on lifecycle manager stop")
			}

			mockListenerImplementation.AssertExpectations(t)

			assert.Equal(t, tc.endStatus, listener.Config.Status,
				"expected status to be updated to %s, got %s", tc.endStatus, listener.Config.Status)
		})
	}
}

func TestScheduledLifecycle_Stop(t *testing.T) {
	tests := []struct {
		name      string
		config    ListenerConfig
		endStatus string
	}{
		{
			name:      "successful-scheduled-stop",
			config:    makeValidListenerConfig("scheduled-listener-id", "http", StatusRunning),
			endStatus: StatusStopping,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create scheduled listener
			tc.config.Deployment = DeploymentExternal
			tc.config.Lifecycle = LifecycleScheduled
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			err = listener.lifecycleManager.Stop(context.Background(), listener)

			assert.NoError(t, err, "Expected no error on scheduled lifecycle manager stop")

			assert.Equal(t, tc.endStatus, listener.Config.Status,
				"expected status to be updated to %s, got %s", tc.endStatus, listener.Config.Status)
		})
	}
}

func TestManagedLifecycle_Terminate(t *testing.T) {
	tests := []struct {
		name           string
		config         ListenerConfig
		endStatus      string
		terminateError error
	}{
		{
			name:           "successful-terminate",
			config:         makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			endStatus:      StatusTerminating,
			terminateError: nil,
		},
		{
			name:           "terminate-error",
			config:         makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			endStatus:      StatusTerminating,
			terminateError: errors.New("listener failed to terminate"),
		},
		{
			name:           "nil-listener",
			config:         makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			endStatus:      StatusTerminating,
			terminateError: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.config.Lifecycle = LifecycleManaged
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			if tc.name == "nil-listener" {
				listener.listener = nil
				assert.Panics(t, func() { listener.Stop(context.Background()) }, "the function did not panic on a nil listener")
				return
			}

			mockListenerImplementation := &MockListenerImplementation{}
			mockListenerImplementation.On("Terminate", mock.Anything).Return(tc.terminateError).Once()
			listener.listener = mockListenerImplementation

			err = listener.lifecycleManager.Terminate(context.Background(), listener)
			if tc.terminateError != nil {
				assert.Error(t, err, "Expected an error on lifecycle manager terminate")
			} else {
				assert.NoError(t, err, "Expected no error on lifecycle manager terminate")
			}

			mockListenerImplementation.AssertExpectations(t)

			assert.Equal(t, tc.endStatus, listener.Config.Status,
				"expected status to be updated to %s, got %s", tc.endStatus, listener.Config.Status)
		})
	}
}

func TestScheduledLifecycle_Terminate(t *testing.T) {
	tests := []struct {
		name      string
		config    ListenerConfig
		endStatus string
	}{
		{
			name:      "successful-scheduled-terminate",
			config:    makeValidListenerConfig("scheduled-listener-id", "http", StatusRunning),
			endStatus: StatusTerminating,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create scheduled listener
			tc.config.Deployment = DeploymentExternal
			tc.config.Lifecycle = LifecycleScheduled
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			err = listener.lifecycleManager.Terminate(context.Background(), listener)

			assert.NoError(t, err, "Expected no error on scheduled lifecycle manager terminate")

			assert.Equal(t, tc.endStatus, listener.Config.Status,
				"expected status to be updated to %s, got %s", tc.endStatus, listener.Config.Status)
		})
	}
}
