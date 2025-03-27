package http_listener

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestHTTPListenerEndToEnd(t *testing.T) {
	// Use a random port to avoid conflicts
	testPort := 8765
	testHost := "localhost"

	// Create a new HTTP listener
	t.Log("Creating HTTP listener")
	listener := NewHTTPListener(testHost, testPort, false, "", "")
	if listener == nil {
		t.Fatal("Failed to create HTTP listener")
	}

	// Create contexts for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start the listener
	t.Log("Starting HTTP listener")
	err := listener.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start HTTP listener: %v", err)
	}

	// Give the server a moment to fully start
	time.Sleep(100 * time.Millisecond)

	// Test connectivity by making a request to the health endpoint
	t.Log("Testing HTTP listener connectivity")
	url := fmt.Sprintf("http://%s:%d/health", testHost, testPort)

	// Create a client with timeout
	client := &http.Client{Timeout: 5 * time.Second}

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to connect to HTTP listener: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	t.Logf("Successfully connected to HTTP listener at %s", url)

	// Stop the listener
	t.Log("Stopping HTTP listener")
	err = listener.Stop(ctx)
	if err != nil {
		t.Fatalf("Failed to stop HTTP listener: %v", err)
	}

	// Verify listener is stopped by trying to connect again
	t.Log("Verifying listener is stopped")
	_, err = client.Get(url)
	if err == nil {
		t.Fatal("HTTP listener is still accepting connections after stopping")
	}
	t.Log("Confirmed listener is stopped")

	// Test restarting the listener
	t.Log("Restarting HTTP listener")
	err = listener.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to restart HTTP listener: %v", err)
	}

	// Give the server a moment to fully start
	time.Sleep(100 * time.Millisecond)

	// Test connectivity again
	t.Log("Testing connectivity after restart")
	resp, err = client.Get(url)
	if err != nil {
		t.Fatalf("Failed to connect to restarted HTTP listener: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}
	t.Log("Successfully connected to restarted listener")

	// Terminate the listener (force close)
	t.Log("Terminating HTTP listener")
	err = listener.Terminate(ctx)
	if err != nil {
		t.Fatalf("Failed to terminate HTTP listener: %v", err)
	}

	// Verify listener is terminated by trying to connect again
	t.Log("Verifying listener is terminated")
	_, err = client.Get(url)
	if err == nil {
		t.Fatal("HTTP listener is still accepting connections after termination")
	}
	t.Log("Confirmed listener is terminated")
}

func TestHTTPListenerWithTLS(t *testing.T) {
	// Skip this test if no certificates are available
	// You could generate test certificates for this test or skip if not available
	t.Skip("Skipping TLS test - provide test certificates to enable")

	// Similar to the above test but with TLS enabled
	// This would require test certificates
}

func TestHTTPListenerUpdateConfig(t *testing.T) {
	// Test the configuration update functionality
	testPort := 8766
	testHost := "localhost"

	listener := NewHTTPListener(testHost, testPort, false, "", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Start the listener
	err := listener.Start(ctx)
	if err != nil {
		t.Fatalf("Failed to start HTTP listener: %v", err)
	}

	// Update configuration
	listener.Port = 8767 // Change the port

	// Apply the configuration update
	err = listener.UpdateConfig(ctx)
	if err != nil {
		t.Fatalf("Failed to update configuration: %v", err)
	}

	// Verify the new configuration works
	url := fmt.Sprintf("http://%s:%d/health", testHost, 8767)
	client := &http.Client{Timeout: 5 * time.Second}

	// Give the server a moment to fully restart with new config
	time.Sleep(100 * time.Millisecond)

	resp, err := client.Get(url)
	if err != nil {
		t.Fatalf("Failed to connect to reconfigured HTTP listener: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	}

	// Clean up
	listener.Stop(ctx)
}
