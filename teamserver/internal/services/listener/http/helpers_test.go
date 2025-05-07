package http_listener

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
)

// Helper function to build default test listener
func buildDefaultTestHTTPListener(t *testing.T) *HTTPListener {
	mockCheckinController := &checkin.CheckInController{}

	listener, err := NewHTTPListener(
		HTTPListenerConfig{
			Host:         "127.0.0.1",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		mockCheckinController,
	)
	if err != nil {
		t.Fatalf("failed to build http listener: %v", err)
	}

	return listener
}

// Helper function to check if a server is listening
func isServerListening(addr string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			conn.Close()
			return true
		}
		time.Sleep(100 * time.Millisecond)
	}
	return false
}

// Helper function to clean up a listener
func cleanupListener(t *testing.T, listener *HTTPListener) {
	if listener != nil && listener.isRunning {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := listener.Stop(ctx)
		if err != nil {
			t.Logf("Failed to shut down listener: %v", err)
		}
	}
}
