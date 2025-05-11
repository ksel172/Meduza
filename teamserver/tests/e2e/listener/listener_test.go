package listener_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	listener_service "github.com/ksel172/Meduza/teamserver/internal/services/listener"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/internal/storage/repos"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/stretchr/testify/require"
)

type Container struct {
	ListenerController *handlers.ListenerController
	ListenerService    *listener_service.ListenerService
}

func setup(t *testing.T) (Container, error) {
	// Setup db
	schema := conf.GetMeduzaDbSchema()
	db, err := repos.Setup()
	if err != nil {
		t.Fatalf("failed to setup database")
	}

	// Create required data layers
	listenerDAL := dal.NewListenerDAL(db, schema)

	// Create service and controller
	listenerService := listener_service.NewListenerService(listenerDAL)
	listenerController := handlers.NewListenersHandler(listenerService, listenerDAL)

	return Container{
		ListenerController: listenerController,
		ListenerService:    listenerService,
	}, nil
}

func createListener(t *testing.T, container Container, listenerModel models.Listener) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Marshal the listener model into the request body
	body, err := json.Marshal(listenerModel)
	if err != nil {
		t.Fatalf("failed to marshal listener model")
	}

	// Create the request
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	// Create the listener by sending a request to the controller
	container.ListenerController.CreateListener(c)

	// Define response wrapper
	var response struct {
		Status  int               `json:"status"`
		Message string            `json:"message"`
		Data    []models.Listener `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	fmt.Printf("%+v", response)

	// Ensure it was created
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 CREATED response")
}

func getListener(t *testing.T, container Container) models.Listener {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Make request
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	container.ListenerController.GetAllListeners(c)

	require.Equal(t, http.StatusOK, w.Code, "expected 200 OK from GetAllListeners")

	// Define response wrapper
	var response struct {
		Status  int               `json:"status"`
		Message string            `json:"message"`
		Data    []models.Listener `json:"data"`
	}

	// Parse response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "failed to parse response body")
	require.Len(t, response.Data, 1, "expected exactly one listener")

	return response.Data[0]
}

// e2e test for the listener service
func TestHTTPListenerEndToEnd(t *testing.T) {

	container, err := setup(t)
	if err != nil {
		t.Fatalf("failed to prepare test dependencies container")
	}

	// seed the database with the listener
	createListener(t, container, models.Listener{
		Kind:        models.HTTPListenerKind,
		Name:        "test-listener",
		Description: "listener for end to end testings",
	})
	t.Log("inserted listened into the database succesfully")

	createdListener := getListener(t, container)
	t.Logf("retrieved listener from database: %+v", createdListener)

	// Create contexts for operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Start the listener that was created
	if err := container.ListenerService.StartListener(ctx, createdListener.ID); err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}

	// // Start the listener
	// t.Log("Starting HTTP listener")
	// err = listener.Start(ctx)
	// if err != nil {
	// 	t.Fatalf("Failed to start HTTP listener: %v", err)
	// }

	// // Give the server a moment to fully start
	// time.Sleep(100 * time.Millisecond)

	// // Test connectivity by pinging the
	// t.Log("Testing HTTP listener connectivity")
	// url := fmt.Sprintf("http://%s:%d/", testHost, testPort)

	// // Create a client with timeout
	// client := &http.Client{Timeout: 5 * time.Second}

	// resp, err := client.Get(url)
	// if err != nil {
	// 	t.Fatalf("Failed to connect to HTTP listener: %v", err)
	// }
	// defer resp.Body.Close()

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
