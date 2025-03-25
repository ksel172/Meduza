package controller

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func makeValidListenerConfig(id, kind, status string) ListenerConfig {
	// It doesn't matter if listeners are managed, scheduled
	// local or external for tests at this level
	return ListenerConfig{
		ID:         id,
		Kind:       kind,
		Host:       "localhost",
		Port:       8000,
		Status:     status,
		Heartbeat:  30,
		Lifecycle:  LifecycleManaged, // always managed for local listeners
		Deployment: DeploymentLocal,
	}
}

func TestManager_Start(t *testing.T) {

	tests := []struct {
		name        string
		config      ListenerConfig
		listenerID  string // ID provided by the user to look for
		expectError bool   // manager.Start error()
		startError  error  // listener.Start() error in goroutine
	}{
		{
			name:        "successful-start",
			config:      makeValidListenerConfig("test-listener-id", "http", StatusReady),
			listenerID:  "test-listener-id",
			expectError: false,
			startError:  nil,
		},
		{
			name:        "listener-not-found",
			config:      makeValidListenerConfig("existing-id", "http", StatusReady),
			listenerID:  "not-found-listener-id",
			expectError: true,
			startError:  nil,
		},
		{ // listener.Start returns an error, the error is written to errChan
			name:        "listener-start-error",
			config:      makeValidListenerConfig("error-listener-id", "http", StatusReady),
			listenerID:  "error-listener-id",
			expectError: false,
			startError:  errors.New("listener start error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create listener with mock lifecycle
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			// Define mock behavior for specific test, need to pass listener as argument on calls
			// Using mock.Anything for context, applies to all cases except listener-not-found
			mockLifecycleManager := &MockLifecycleManager{}
			if tc.name != "listener-not-found" {
				mockLifecycleManager.On("Start", mock.Anything, listener).Return(tc.startError).Once()
			}
			listener.lifecycleManager = mockLifecycleManager

			// Create manager
			m, err := NewListenerManager(map[string]*Listener{listener.ID: listener})
			if !assert.NoError(t, err, "error creating manager") {
				assert.FailNow(t, "error creating manager")
			}

			// Goroutine control
			errChan := make(chan error)
			ctx := context.Background()

			// Test manager
			err = m.startListener(ctx, tc.listenerID, errChan)
			if tc.expectError {
				assert.Error(t, err, "Expected an error on manager start")
				return // Won't check goroutine when error is returned early, it won't have spawned
			} else {
				assert.NoError(t, err, "Expected no error on manager start")
			}

			// Test goroutine execution
			select {
			case _, ok := <-errChan:
				// If channel was closed, it means goroutine finished execution with no errors
				if !ok {
					assert.Nil(t, tc.startError, "Expected no error from listener.Start method")
				} else {
					assert.NotNil(t, tc.startError, "Unexpected error from listener.Start method")
				}
			case <-ctx.Done():
				t.Fatal("Timeout waiting for goroutine to complete")
			}

			mockLifecycleManager.AssertExpectations(t)
		})
	}
}

func TestManager_Add(t *testing.T) {
	tests := []struct {
		name        string
		config      ListenerConfig
		listeners   map[string]*Listener
		expectedLen int
		expectError bool
	}{
		{
			name:        "successful-add",
			config:      makeValidListenerConfig("test-listener-id", "http", StatusReady),
			listeners:   make(map[string]*Listener),
			expectedLen: 1,
			expectError: false,
		},
		{
			name:        "listener-exists",
			config:      makeValidListenerConfig("existing-listener-id", "http", StatusReady),
			listeners:   map[string]*Listener{"existing-listener-id": {ID: "existing-listener-id"}},
			expectedLen: 1,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			// Create manager
			m, err := NewListenerManager(tc.listeners)
			if !assert.NoError(t, err, "error creating manager") {
				assert.FailNow(t, "error creating manager")
			}

			err = m.addListener(tc.config)
			if tc.expectError {
				assert.Error(t, err, "expected error adding listener")
			} else {
				assert.NoError(t, err, "expected no error adding listener")
			}

			assert.Equal(t, len(m.listeners), tc.expectedLen)
		})
	}
}

func TestManager_Terminate(t *testing.T) {
	tests := []struct {
		name           string
		config         ListenerConfig
		listenerID     string
		expectError    bool
		terminateError error
		shouldRemove   bool
	}{
		{
			name:           "successful-terminate",
			config:         makeValidListenerConfig("test-listener-id", "http", StatusReady),
			listenerID:     "test-listener-id",
			expectError:    false,
			terminateError: nil,
			shouldRemove:   true,
		},
		{
			name:           "listener-not-found",
			config:         makeValidListenerConfig("existing-id", "http", StatusReady),
			listenerID:     "not-found-listener-id",
			expectError:    true,
			terminateError: nil,
		},
		{
			name:           "listener-terminate-error",
			config:         makeValidListenerConfig("error-listener-id", "http", StatusReady),
			listenerID:     "error-listener-id",
			expectError:    true,
			terminateError: errors.New("listener start error"),
			shouldRemove:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			mockLifecycleManager := &MockLifecycleManager{}
			if tc.name != "listener-not-found" {
				mockLifecycleManager.On("Terminate", mock.Anything, listener).Return(tc.terminateError).Once()
			}
			listener.lifecycleManager = mockLifecycleManager

			listeners := map[string]*Listener{listener.ID: listener}
			m, err := NewListenerManager(listeners)
			if !assert.NoError(t, err, "error creating manager") {
				assert.FailNow(t, "error creating manager")
			}

			err = m.terminateListener(context.Background(), tc.listenerID)
			if tc.expectError {
				assert.Error(t, err, "Expected error on terminate")
			} else {
				assert.NoError(t, err, "Expected no error on terminate")
			}

			// Check if listener was removed or not from the manager listeners map
			if tc.name != "listener-not-found" {
				if _, ok := listeners[tc.listenerID]; !ok { // listener was removed
					if !tc.shouldRemove {
						t.Errorf("listener was removed from manager but shouldn't")
					}
				} else { // listener was not removed
					if tc.shouldRemove {
						t.Errorf("listener was not removed from manager")
					}
				}
			}

			mockLifecycleManager.AssertExpectations(t)

			// Test
		})
	}
}

func TestManager_Stop(t *testing.T) {
	tests := []struct {
		name        string
		config      ListenerConfig
		listenerID  string
		expectError bool
		stopError   error
	}{
		{
			name:        "successful-stop",
			config:      makeValidListenerConfig("test-listener-id", "http", StatusRunning),
			listenerID:  "test-listener-id",
			expectError: false,
			stopError:   nil,
		},
		{
			name:        "listener-not-found",
			config:      makeValidListenerConfig("existing-id", "http", StatusRunning),
			listenerID:  "not-found-listener-id",
			expectError: true,
			stopError:   nil,
		},
		{
			name:        "listener-stop-error",
			config:      makeValidListenerConfig("error-listener-id", "http", StatusRunning),
			listenerID:  "error-listener-id",
			expectError: false,
			stopError:   errors.New("listener stop error"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			listener, err := NewListenerFromConfig(tc.config)
			if !assert.NoError(t, err, "error creating listener from config") {
				assert.FailNow(t, "error creating listener from config")
			}

			mockLifecycleManager := &MockLifecycleManager{}
			if tc.name != "listener-not-found" {
				mockLifecycleManager.On("Stop", mock.Anything, listener).Return(tc.stopError).Once()
			}
			listener.lifecycleManager = mockLifecycleManager

			m, err := NewListenerManager(map[string]*Listener{listener.ID: listener})
			if !assert.NoError(t, err, "error creating manager") {
				assert.FailNow(t, "error creating manager")
			}

			errChan := make(chan error)
			ctx := context.Background()

			err = m.stopListener(ctx, tc.listenerID, errChan)
			if tc.expectError {
				assert.Error(t, err, "Expected an error on manager stop")
				return
			} else {
				assert.NoError(t, err, "Expected no error on manager stop")
			}

			select {
			case err, ok := <-errChan:
				if !ok {
					assert.Nil(t, tc.stopError, "Expected no error from listener.Stop method")
				} else {
					assert.NotNil(t, tc.stopError, "Expected error from listener.Stop method")
					assert.Equal(t, "failed to stop listener: listener stop error", err.Error())
				}
			case <-ctx.Done():
				t.Fatal("Timeout waiting for goroutine to complete")
			}

			mockLifecycleManager.AssertExpectations(t)
		})
	}
}
