package listener_test

import (
	"context"
	"testing"
	"time"

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
	t.Run("stop listener", func(t *testing.T) {
		err := container.ListenerService.StopListener(context.Background(), listener.ID)
		require.NoErrorf(t, err, "failed to stop listener")

		// Listener stauts is updated async, must wait for a moment for updates to make it to database
		time.Sleep(250 * time.Millisecond)

		err = testAgent.Authenticate(t, listener.ID)
		require.Error(t, err)

		stoppedListener := getListener(t, container, listener.Name)
		if stoppedListener.Status != models.StatusReady {
			t.Errorf("stopped listener status is not back to ready")
		}
	})

	// Terminate the listener
	t.Run("terminate listener", func(t *testing.T) {
		err := container.ListenerService.TerminateListener(context.Background(), listener.ID)
		require.NoErrorf(t, err, "failed to terminate listener")

		// Listener stauts is updated async, must wait for a moment for updates to make it to database
		time.Sleep(250 * time.Millisecond)

		err = testAgent.Authenticate(t, listener.ID)
		require.Error(t, err)

		terminatedListener := getListener(t, container, listener.Name)
		if terminatedListener.Status != models.StatusPending {
			t.Errorf("terminated listener status is not back to ready")
		}
	})
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
