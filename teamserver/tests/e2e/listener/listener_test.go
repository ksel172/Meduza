package listener_test

import (
	"context"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
)

func TestListenerService(t *testing.T) {
	gin.SetMode("test")

	container, listener, authToken, err := setup(t, models.CreateLocalListenerRequest{
		Kind:        models.HTTPListenerKind,
		Name:        "test-listener",
		Description: "listener for end to end testings",
	})
	if err != nil {
		t.Fatalf("failed to prepare test dependencies container")
	}

	// Create contexts for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start the listener that was created
	t.Log("Starting listener...")
	if err := container.ListenerService.StartListener(ctx, listener.ID); err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	// Create test agent to send requests to listener
	testAgent := newTestHTTPAgent(listener.Host, listener.Port, authToken)

	// Agent sends authentication request to listener
	testAgent.Authenticate(t, listener.ID)

	// Test connectivity by pinging the listener
	// t.Log("Testing HTTP listener connectivity")
	// response, err := testClient.Ping()
	// if err != nil {
	// 	t.Fatalf("Failed to connect to HTTP listener: %v", err)
	// }
	// t.Logf("Ping response: %s", response)

	// if resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	// }
	// t.Logf("Successfully connected to HTTP listener at %s", url)

	// // Stop the listener
	// t.Log("Stopping HTTP listener")
	// err = listener.Stop(ctx)
	// if err != nil {
	// 	t.Fatalf("Failed to stop HTTP listener: %v", err)
	// }

	// // Verify listener is stopped by trying to connect again
	// t.Log("Verifying listener is stopped")
	// _, err = client.Get(url)
	// if err == nil {
	// 	t.Fatal("HTTP listener is still accepting connections after stopping")
	// }
	// t.Log("Confirmed listener is stopped")

	// // Test restarting the listener
	// t.Log("Restarting HTTP listener")
	// err = listener.Start(ctx)
	// if err != nil {
	// 	t.Fatalf("Failed to restart HTTP listener: %v", err)
	// }

	// // Give the server a moment to fully start
	// time.Sleep(100 * time.Millisecond)

	// // Test connectivity again
	// t.Log("Testing connectivity after restart")
	// resp, err = client.Get(url)
	// if err != nil {
	// 	t.Fatalf("Failed to connect to restarted HTTP listener: %v", err)
	// }
	// defer resp.Body.Close()

	// if resp.StatusCode != http.StatusOK {
	// 	t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
	// }
	// t.Log("Successfully connected to restarted listener")

	// // Terminate the listener (force close)
	// t.Log("Terminating HTTP listener")
	// err = listener.Terminate(ctx)
	// if err != nil {
	// 	t.Fatalf("Failed to terminate HTTP listener: %v", err)
	// }

	// // Verify listener is terminated by trying to connect again
	// t.Log("Verifying listener is terminated")
	// _, err = client.Get(url)
	// if err == nil {
	// 	t.Fatal("HTTP listener is still accepting connections after termination")
	// }
	// t.Log("Confirmed listener is terminated")
}

func TestHTTPListenerWithTLS(t *testing.T) {
	t.Skip("Skipping TLS test - provide test certificates to enable")
}

// func TestHTTPListenerUpdateConfig(t *testing.T) {
// 	// Test the configuration update functionality
// 	testPort := 8766
// 	testHost := "localhost"

// 	listener := NewHTTPListener(testHost, testPort, false, "", "")
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()

// 	// Start the listener
// 	err := listener.Start(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to start HTTP listener: %v", err)
// 	}

// 	// Update configuration
// 	listener.Port = 8767 // Change the port

// 	// Apply the configuration update
// 	err = listener.UpdateConfig(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to update configuration: %v", err)
// 	}

// 	// Verify the new configuration works
// 	url := fmt.Sprintf("http://%s:%d/health", testHost, 8767)
// 	client := &http.Client{Timeout: 5 * time.Second}

// 	// Give the server a moment to fully restart with new config
// 	time.Sleep(100 * time.Millisecond)

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to reconfigured HTTP listener: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
// 	}

// 	// Clean up
// 	listener.Stop(ctx)
// }
