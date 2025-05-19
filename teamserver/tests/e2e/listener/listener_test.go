package listener_test

import (
	"context"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/require"
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

	testAgent, err := newTestHTTPAgent(listener.Host, listener.Port, authToken)
	if err != nil {
		t.Fatalf("failed to create agent: %v", err)
	}
	t.Logf("created test agent")

	// Start the listener that was created
	t.Log("Starting listener...")
	if err := container.ListenerService.StartListener(context.Background(), listener.ID); err != nil {
		t.Fatalf("failed to start listener: %v", err)
	}
	t.Logf("Listener started succesfully: %s:%d...", listener.Host, listener.Port)

	// Agent sends authentication request to listener
	t.Run("agent authenticate", func(t *testing.T) {
		err := testAgent.Authenticate(t, listener.ID)
		require.NoError(t, err)
	})
	t.Run("agent register", func(t *testing.T) {
		testAgent.Register(t)
		require.NoError(t, err)
	})

	// Retrieve the registered agent and store it (the testAgent knows its ID because its returned in the register, however, that's all it knows)
	agent := getAgent(t, container, testAgent.ID)
	testAgent.Agent = agent

	// Seed db with AgentTask
	createAgentTask(t, container, testAgent.ID, models.AgentTaskRequest{
		Type:   models.TaskShellCommand,
		Status: models.TaskStatusQueued,
		Command: models.AgentCommand{
			Name:       "test-command-name",
			Parameters: []string{"test-parameter-one", "test-parameter-two"},
		},
	})
	t.Log("Created AgentTask in database")

	// Test tasks & response endpoint
	// Unfinished implementations at the moment
	t.Run("agent tasks", func(t *testing.T) {
		err := testAgent.GetTasks(t)
		require.NoError(t, err)
	})
	// t.Run("agent response: ", testAgent.SendResponse(t))

	// Stop the listener
	container.ListenerService.StopListener(context.Background(), listener.ID)
	t.Run("stopped listener: authenticate attempt", func(t *testing.T) {
		err := testAgent.Authenticate(t, listener.ID)
		require.Error(t, err)
	})

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
