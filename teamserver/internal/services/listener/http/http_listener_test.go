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

	// Test cases grid
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
				listener.startServerGoroutine = func(errChan chan<- error, readyChan chan<- struct{}) {
					go func() {
						// Sleep longer than the timeout
						time.Sleep(10 * time.Second)

						// Simulate no error, if context cancellation is not working
						readyChan <- struct{}{}
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
					cleanupListener(t, listener)
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
					t.Errorf("expected %s, got: %s", err, tt.expectedError)
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

			// Cleanup listener on every
			isListenerCleanedUp = true // for not trying to clean up again on defer
			cleanupListener(t, listener)
		})
	}
}

// func TestHTTPListenerEndToEnd(t *testing.T) {
// 	// Use a random port to avoid conflicts
// 	testPort := 8765
// 	testHost := "localhost"

// 	// Create a new HTTP listener
// 	t.Log("Creating HTTP listener")
// 	listener := NewHTTPListener(testHost, testPort, false, "", "")
// 	if listener == nil {
// 		t.Fatal("Failed to create HTTP listener")
// 	}

// 	// Create contexts for operations
// 	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
// 	defer cancel()

// 	// Start the listener
// 	t.Log("Starting HTTP listener")
// 	err := listener.Start(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to start HTTP listener: %v", err)
// 	}

// 	// Give the server a moment to fully start
// 	time.Sleep(100 * time.Millisecond)

// 	// Test connectivity by making a request to the health endpoint
// 	t.Log("Testing HTTP listener connectivity")
// 	url := fmt.Sprintf("http://%s:%d/health", testHost, testPort)

// 	// Create a client with timeout
// 	client := &http.Client{Timeout: 5 * time.Second}

// 	resp, err := client.Get(url)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to HTTP listener: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
// 	}
// 	t.Logf("Successfully connected to HTTP listener at %s", url)

// 	// Stop the listener
// 	t.Log("Stopping HTTP listener")
// 	err = listener.Stop(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to stop HTTP listener: %v", err)
// 	}

// 	// Verify listener is stopped by trying to connect again
// 	t.Log("Verifying listener is stopped")
// 	_, err = client.Get(url)
// 	if err == nil {
// 		t.Fatal("HTTP listener is still accepting connections after stopping")
// 	}
// 	t.Log("Confirmed listener is stopped")

// 	// Test restarting the listener
// 	t.Log("Restarting HTTP listener")
// 	err = listener.Start(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to restart HTTP listener: %v", err)
// 	}

// 	// Give the server a moment to fully start
// 	time.Sleep(100 * time.Millisecond)

// 	// Test connectivity again
// 	t.Log("Testing connectivity after restart")
// 	resp, err = client.Get(url)
// 	if err != nil {
// 		t.Fatalf("Failed to connect to restarted HTTP listener: %v", err)
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		t.Fatalf("Expected status code 200, got %d", resp.StatusCode)
// 	}
// 	t.Log("Successfully connected to restarted listener")

// 	// Terminate the listener (force close)
// 	t.Log("Terminating HTTP listener")
// 	err = listener.Terminate(ctx)
// 	if err != nil {
// 		t.Fatalf("Failed to terminate HTTP listener: %v", err)
// 	}

// 	// Verify listener is terminated by trying to connect again
// 	t.Log("Verifying listener is terminated")
// 	_, err = client.Get(url)
// 	if err == nil {
// 		t.Fatal("HTTP listener is still accepting connections after termination")
// 	}
// 	t.Log("Confirmed listener is terminated")
// }

// func TestHTTPListenerWithTLS(t *testing.T) {
// 	// Skip this test if no certificates are available
// 	// You could generate test certificates for this test or skip if not available
// 	t.Skip("Skipping TLS test - provide test certificates to enable")

// 	// Similar to the above test but with TLS enabled
// 	// This would require test certificates
// }

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

// Function to simulate agent authentication and retrieve sessiontoken + AES key for message encryption
// func authenticateAgent() (string, []byte, error) {
// 	mockAgentDAL := new(mocks.MockAgentDAL)
// 	mockCheckInDal := new(mocks.MockCheckInDal)
// 	mockPayloaDAL := new(mocks.MockPayloadDAL)
// 	controller := NewCheckInController(mockAgentDAL)

// 	// Generates agent and server keys
// 	agentPrivKey, agentPubKey, err := utils.GenerateECDHKeyPair()
// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to generate agent keys: %v", err)
// 	}
// 	serverPrivKey, serverPubKey, err := utils.GenerateECDHKeyPair()
// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to generate server keys: %v", err)
// 	}

// 	// Create request body
// 	c2request := models.C2Request{Message: base64.StdEncoding.EncodeToString(agentPubKey)}
// 	bodyRawBytes, _ := json.Marshal(c2request)
// 	encodedBodyString := base64.StdEncoding.EncodeToString(bodyRawBytes)
// 	body := []byte(encodedBodyString)

// 	// Agent prepares its authentication request
// 	w := httptest.NewRecorder()
// 	c, _ := gin.CreateTestContext(w)
// 	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
// 	c.Request.Header.Add("Auth-Token", "testAuthToken")

// 	// Ensure auth-token is accepted using the payloadDAL mock, which will also return the server keys
// 	mockPayloaDAL.On("GetKeys", "test-auth-token").Return(serverPrivKey, serverPubKey, nil).Once()

// 	// Submit request to server
// 	controller.Checkin(c)

// 	if w.Code != http.StatusAccepted {
// 		return "", nil, fmt.Errorf("agent auth failed")
// 	}

// 	// Receive server response
// 	serverResponseBase64 := struct {
// 		PublicKey    string `json:"public_key"`
// 		SessionToken string `json:"session_token"`
// 	}{}
// 	if err := json.Unmarshal(w.Body.Bytes(), &serverResponseBase64); err != nil {
// 		return "", nil, fmt.Errorf("invalid response body")
// 	}
// 	serverPublicKey, _ := base64.StdEncoding.DecodeString(serverResponseBase64.PublicKey)

// 	// Use the server public key to derive the shared key
// 	sharedKey, err := utils.DeriveECDHSharedSecret(agentPrivKey, serverPublicKey)
// 	if err != nil {
// 		return "", nil, fmt.Errorf("failed to derive shared key")
// 	}

// 	// SessionToken is sent in base64 in the requests anyway
// 	return serverResponseBase64.SessionToken, sharedKey, nil
// }

// func encryptAgentRequest(c2request models.C2Request, key []byte) ([]byte, error) {
// 	c2requestBytes, err := json.Marshal(c2request)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to marshal c2request")
// 	}
// 	encryptedc2request, err := utils.AesEncrypt(key, c2requestBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to encrypt message")
// 	}
// 	return []byte(base64.StdEncoding.EncodeToString(encryptedc2request)), nil
// }

// func TestAgentRegisterRequest(t *testing.T) {
// 	mockAgentDAL := new(mocks.MockAgentDAL)
// 	mockCheckInDal := new(mocks.MockCheckInDal)
// 	mockPayloaDAL := new(mocks.MockPayloadDAL)
// 	controller := services.NewCheckInController(mockCheckInDal, mockAgentDAL, mockPayloaDAL)
// 	gin.SetMode(gin.TestMode)

// 	encodedSessionToken, aesKey, err := authenticateAgent()
// 	if err != nil {
// 		t.Fatalf("failed agent auth: %v", err)
// 	}

// 	c2request := models.C2Request{
// 		Reason:  models.Register,
// 		AgentID: "test-agent-id",
// 		Message: `{"agent_id": "test-agent-id"}`,
// 	}

// 	tests := []struct {
// 		name           string
// 		c2Request      models.C2Request
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "register agent: success",
// 			c2Request:      c2request,
// 			expectedStatus: http.StatusCreated,
// 		},
// 		{
// 			name:           "register agent: agent already exists",
// 			c2Request:      c2request,
// 			expectedStatus: http.StatusConflict,
// 		},
// 		{
// 			name:           "register agent: create agent error",
// 			c2Request:      c2request,
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 		{
// 			name:           "register agent: create agent info error",
// 			c2Request:      c2request,
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			switch tt.name {
// 			case "register agent: success":
// 				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, errors.New("agent does not exist")).Once()
// 				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(nil).Once()
// 				mockAgentDAL.On("CreateAgentInfo", mock.AnythingOfType("models.AgentInfo")).Return(nil).Once()
// 			case "register agent: agent already exists":
// 				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, nil).Once()
// 			case "register agent: create agent error":
// 				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, errors.New("agent does not exist")).Once()
// 				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(errors.New("failed to create agent")).Once()
// 			case "register agent: create agent info error":
// 				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, errors.New("agent does not exist")).Once()
// 				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(nil).Once()
// 				mockAgentDAL.On("CreateAgentInfo", mock.AnythingOfType("models.AgentInfo")).Return(errors.New("failed to create agent info")).Once()
// 			}

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			body, err := encryptAgentRequest(c2request, aesKey)
// 			if err != nil {
// 				t.Fatal(err.Error())
// 			}
// 			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
// 			c.Request.Header.Add("Session-Token", encodedSessionToken)

// 			controller.Checkin(c)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 			mockCheckInDal.AssertExpectations(t)
// 			mockAgentDAL.AssertExpectations(t)
// 		})
// 	}
// }

// func TestAgentTasksRequest(t *testing.T) {
// 	mockAgentDAL := new(mocks.MockAgentDAL)
// 	mockCheckInDal := new(mocks.MockCheckInDal)
// 	mockPayloaDAL := new(mocks.MockPayloadDAL)
// 	controller := services.NewCheckInController(mockCheckInDal, mockAgentDAL, mockPayloaDAL)
// 	gin.SetMode(gin.TestMode)

// 	encodedSessionToken, aesKey, err := authenticateAgent()
// 	if err != nil {
// 		t.Fatalf("failed agent auth: %v", err)
// 	}

// 	c2request := models.C2Request{
// 		Reason:  models.Task,
// 		AgentID: "test-agent-id",
// 		Message: `{"agent_id": "test-agent-id"}`,
// 	}

// 	tests := []struct {
// 		name           string
// 		c2request      models.C2Request
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "agent task: success",
// 			c2request:      c2request,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "agent task: get agent tasks error",
// 			c2request:      c2request,
// 			expectedStatus: http.StatusNotFound,
// 		},
// 		{
// 			name:           "agent task: callback update error",
// 			c2request:      c2request,
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			switch tt.name {
// 			case "agent task: success":
// 				mockAgentDAL.On("GetAgentTasks", tt.c2request.AgentID).Return(make([]models.AgentTask, 0), nil).Once()
// 				mockAgentDAL.On("UpdateAgentLastCallback", tt.c2request.AgentID, mock.AnythingOfType("string")).Return(nil).Once()
// 			case "agent task: get agent tasks error":
// 				mockAgentDAL.On("GetAgentTasks", tt.c2request.AgentID).Return(make([]models.AgentTask, 0), errors.New("failed to get agent tasks")).Once()
// 			case "agent task: callback update error":
// 				mockAgentDAL.On("GetAgentTasks", tt.c2request.AgentID).Return(make([]models.AgentTask, 0), nil).Once()
// 				mockAgentDAL.On("UpdateAgentLastCallback", tt.c2request.AgentID, mock.AnythingOfType("string")).Return(errors.New("failed to update agent last callback")).Once()
// 			}

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			body, err := encryptAgentRequest(c2request, aesKey)
// 			if err != nil {
// 				t.Fatal(err.Error())
// 			}
// 			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
// 			c.Request.Header.Add("Session-Token", encodedSessionToken)

// 			controller.Checkin(c)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 			mockCheckInDal.AssertExpectations(t)
// 			mockAgentDAL.AssertExpectations(t)
// 		})
// 	}
// }

// func TestAgentResponseRequest(t *testing.T) {
// 	mockAgentDAL := new(mocks.MockAgentDAL)
// 	mockCheckInDal := new(mocks.MockCheckInDal)
// 	mockPayloaDAL := new(mocks.MockPayloadDAL)
// 	controller := services.NewCheckInController(mockCheckInDal, mockAgentDAL, mockPayloaDAL)
// 	gin.SetMode(gin.TestMode)

// 	encodedSessionToken, aesKey, err := authenticateAgent()
// 	if err != nil {
// 		t.Fatalf("failed agent auth: %v", err)
// 	}

// 	c2request := models.C2Request{
// 		Reason:  models.Response,
// 		AgentID: "test-agent-id",
// 		Message: `{"agent_id": "test-agent-id"}`,
// 	}

// 	tests := []struct {
// 		name           string
// 		c2request      models.C2Request
// 		expectedStatus int
// 	}{
// 		{
// 			name:           "agent response: success",
// 			c2request:      c2request,
// 			expectedStatus: http.StatusOK,
// 		},
// 		{
// 			name:           "agent response: update agent task error",
// 			c2request:      c2request,
// 			expectedStatus: http.StatusInternalServerError,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			switch tt.name {
// 			case "agent response: success":
// 				mockAgentDAL.On("UpdateAgentTask", mock.AnythingOfType("models.AgentTask")).Return(nil).Once()
// 			case "agent response: update agent task error":
// 				mockAgentDAL.On("UpdateAgentTask", mock.AnythingOfType("models.AgentTask")).Return(errors.New("failed to update agent task")).Once()
// 			}

// 			w := httptest.NewRecorder()
// 			c, _ := gin.CreateTestContext(w)

// 			body, err := encryptAgentRequest(c2request, aesKey)
// 			if err != nil {
// 				t.Fatal(err.Error())
// 			}
// 			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
// 			c.Request.Header.Add("Session-Token", encodedSessionToken)

// 			controller.Checkin(c)

// 			assert.Equal(t, tt.expectedStatus, w.Code)
// 			mockCheckInDal.AssertExpectations(t)
// 			mockAgentDAL.AssertExpectations(t)
// 		})
// 	}
// }
