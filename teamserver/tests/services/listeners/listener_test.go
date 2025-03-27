package service_tests

// import (
// 	"context"
// 	"errors"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

// func TestListener_Start(t *testing.T) {
// 	tests := []struct {
// 		name       string
// 		config     ListenerConfig
// 		startError error
// 	}{
// 		{
// 			name:       "successful-start",
// 			config:     makeValidListenerConfig("test-listener-id", "http", StatusReady),
// 			startError: nil,
// 		},
// 		{
// 			name:       "listener-not-ready",
// 			config:     makeValidListenerConfig("pending-listener-id", "http", StatusPending),
// 			startError: errors.New("listener not ready to start"),
// 		},
// 		{
// 			name:       "start-error",
// 			config:     makeValidListenerConfig("test-listener-id", "http", StatusReady),
// 			startError: errors.New("lifecycle manager failed to start listener"),
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			listener, err := NewListenerFromConfig(tc.config)
// 			if !assert.NoError(t, err, "error creating listener from config") {
// 				assert.FailNow(t, "error creating listener from config")
// 			}

// 			mockLifecycleManager := &MockLifecycleManager{}
// 			if tc.name != "listener-not-ready" {
// 				mockLifecycleManager.On("Start", mock.Anything, listener).Return(tc.startError).Once()
// 			}
// 			listener.lifecycleManager = mockLifecycleManager

// 			err = listener.Start(context.Background())
// 			if tc.startError != nil {
// 				assert.Error(t, err, "Expected an error on lifecycle manager start")
// 			} else {
// 				assert.NoError(t, err, "Expected no error on lifecycle manager start")
// 			}

// 			mockLifecycleManager.AssertExpectations(t)
// 		})
// 	}
// }

// func TestListener_Stop(t *testing.T) {
// 	tests := []struct {
// 		name      string
// 		config    ListenerConfig
// 		stopError error
// 	}{
// 		{
// 			name:      "successful-stop",
// 			config:    makeValidListenerConfig("test-listener-id", "http", StatusRunning),
// 			stopError: nil,
// 		},
// 		{
// 			name:      "listener-not-running",
// 			config:    makeValidListenerConfig("pending-listener-id", "http", StatusPending),
// 			stopError: errors.New("listener not running, cannot be stopped"),
// 		},
// 		{
// 			name:      "stop-error",
// 			config:    makeValidListenerConfig("test-listener-id", "http", StatusRunning),
// 			stopError: errors.New("lifecycle manager failed to stop listener"),
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			listener, err := NewListenerFromConfig(tc.config)
// 			if !assert.NoError(t, err, "error creating listener from config") {
// 				assert.FailNow(t, "error creating listener from config")
// 			}

// 			mockLifecycleManager := &MockLifecycleManager{}
// 			if tc.name != "listener-not-running" {
// 				mockLifecycleManager.On("Stop", mock.Anything, listener).Return(tc.stopError).Once()
// 			}
// 			listener.lifecycleManager = mockLifecycleManager

// 			err = listener.Stop(context.Background())
// 			if tc.stopError != nil {
// 				assert.Error(t, err, "Expected an error on lifecycle manager stop")
// 			} else {
// 				assert.NoError(t, err, "Expected no error on lifecycle manager stop")
// 			}

// 			mockLifecycleManager.AssertExpectations(t)
// 		})
// 	}
// }

// func TestListener_Terminate(t *testing.T) {
// 	tests := []struct {
// 		name           string
// 		config         ListenerConfig
// 		terminateError error
// 	}{
// 		{
// 			name:           "successful-terminate",
// 			config:         makeValidListenerConfig("test-listener-id", "http", StatusReady),
// 			terminateError: nil,
// 		},
// 		{
// 			name:           "terminate-error",
// 			config:         makeValidListenerConfig("test-listener-id", "http", StatusReady),
// 			terminateError: errors.New("lifecycle manager failed to terminate listener"),
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			listener, err := NewListenerFromConfig(tc.config)
// 			if !assert.NoError(t, err, "error creating listener from config") {
// 				assert.FailNow(t, "error creating listener from config")
// 			}

// 			mockLifecycleManager := &MockLifecycleManager{}
// 			mockLifecycleManager.On("Terminate", mock.Anything, listener).Return(tc.terminateError).Once()
// 			listener.lifecycleManager = mockLifecycleManager

// 			err = listener.Terminate(context.Background())
// 			if tc.terminateError != nil {
// 				assert.Error(t, err, "Expected an error on lifecycle manager terminate")
// 			} else {
// 				assert.NoError(t, err, "Expected no error on lifecycle manager terminate")
// 			}

// 			mockLifecycleManager.AssertExpectations(t)
// 		})
// 	}
// }
