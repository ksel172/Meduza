package http_listener

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	gin.SetMode("test")

	// Setup helper functions
	// setupBlockedPort := func(t *testing.T) (net.Listener, int) {
	// 	// Find an available port and block it
	// 	l, err := net.Listen("tcp", "127.0.0.1:0")
	// 	require.NoError(t, err)
	// 	_, portStr, _ := net.SplitHostPort(l.Addr().String())
	// 	port := 0
	// 	_, err = fmt.Sscanf(portStr, "%d", &port)
	// 	require.NoError(t, err)
	// 	return l, port
	// }

	tests := []struct {
		name           string
		setupFunc      func(t *testing.T) *HTTPListener
		contextTimeout time.Duration
		contextCancel  bool
		expectRunning  bool
		expectError    bool
		expectedError  error
	}{
		{
			name: "successful_start",
			setupFunc: func(t *testing.T) *HTTPListener {
				return buildDefaultTestHTTPListener(t)
			},
			contextTimeout: 5 * time.Second,
			contextCancel:  false,
			expectRunning:  true,
			expectError:    false,
		},
		{
			name: "context_canceled",
			setupFunc: func(t *testing.T) *HTTPListener {
				return buildDefaultTestHTTPListener(t)
			},
			contextTimeout: 5 * time.Second,
			contextCancel:  true,
			expectRunning:  false,
			expectError:    true,
			expectedError:  context.Canceled,
		},
		{
			name: "timeout_waiting_for_start",
			setupFunc: func(t *testing.T) *HTTPListener {
				listener := buildDefaultTestHTTPListener(t)

				// Replace the Start method to simulate a slow startup
				listener.startServerGoroutine = func(errChan chan<- error) {
					go func() {
						time.Sleep(10 * time.Second)
					}()
				}

				return listener
			},
			contextTimeout: 1 * time.Second,
			expectRunning:  false,
			expectError:    true,
			expectedError:  context.DeadlineExceeded,
		},
		// { // Requires port manager implementation for port distribution to listeners
		// 	name: "port_already_in_use",
		// 	setupFunc: func(t *testing.T) *HTTPListener {
		// 		blockedListener, port := setupBlockedPort(t)
		// 		t.Cleanup(func() { blockedListener.Close() })

		// 		mockCheckinController := &checkin.CheckInController{}
		// 		listener, err := NewHTTPListener(
		// 			HTTPListenerConfig{
		// 				Host:         "127.0.0.1",
		// 				Port:         port,
		// 				ReadTimeout:  30,
		// 				WriteTimeout: 30,
		// 			},
		// 			mockCheckinController,
		// 		)
		// 		require.NoError(t, err)
		// 		return listener
		// 	},
		// 	contextTimeout: 5 * time.Second,
		// 	expectRunning:  false,
		// 	expectError:    true,
		// 	expectedError:  nil, //TODO: add error kinds
		// },
	}

	// Execute test grid
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup listener based on test case
			listener := tt.setupFunc(t)

			// Prepare listener clean up
			isListenerCleanedUp := false
			defer func() {
				if !isListenerCleanedUp {
					listener.Stop(context.Background())
				}
			}()

			// Create context with timeout
			ctx, cancel := context.WithTimeout(context.Background(), tt.contextTimeout)
			defer cancel()

			// Apply context cancellation if specified
			if tt.contextCancel {
				cancel()
			}

			// Test the Start method
			err := listener.Start(ctx)

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected %s, got: %s", tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}

			// Verify running state
			assert.Equal(t, tt.expectRunning, listener.isRunning)

			// If expected to be running, verify it's actually listening
			if tt.expectRunning && !tt.contextCancel {
				// Determine the address to check
				addr := fmt.Sprintf("%s:%d", listener.Config.Host, listener.Config.Port)
				connected := isServerListening(addr, 2*time.Second)
				assert.True(t, connected, "Server should be listening on %s", addr)
			}

			// Cleanup listener on every loop
			isListenerCleanedUp = true // for not trying to clean up again on defer
			listener.Stop(context.Background())
		})
	}
}

func TestStop(t *testing.T) {
	gin.SetMode("test")

	// Setup helper functions
	// setupBlockedPort := func(t *testing.T) (net.Listener, int) {
	// 	// Find an available port and block it
	// 	l, err := net.Listen("tcp", "127.0.0.1:0")
	// 	require.NoError(t, err)
	// 	_, portStr, _ := net.SplitHostPort(l.Addr().String())
	// 	port := 0
	// 	_, err = fmt.Sscanf(portStr, "%d", &port)
	// 	require.NoError(t, err)
	// 	return l, port
	// }

	tests := []struct {
		name          string
		setupFunc     func(t *testing.T) *HTTPListener
		expectRunning bool
		expectError   bool
		expectedError error
	}{
		{
			name: "successful_stop",
			setupFunc: func(t *testing.T) *HTTPListener {
				return buildStartDefaultTestHTTPListener(t)
			},
			expectRunning: false,
			expectError:   false,
		},
	}

	// Execute test grid
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup listener based on test case
			listener := tt.setupFunc(t)

			err := listener.Stop(context.Background())

			// Verify results
			if tt.expectError {
				assert.Error(t, err)
				if !errors.Is(err, tt.expectedError) {
					t.Errorf("expected %s, got: %s", tt.expectedError, err)
				}
			} else {
				assert.NoError(t, err)
			}

			// t.Logf("%+v", listener)

			// Verify running state
			assert.Equal(t, tt.expectRunning, listener.isRunning)
		})
	}
}
