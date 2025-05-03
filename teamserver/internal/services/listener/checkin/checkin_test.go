package checkin

import (
	"testing"

	"github.com/gin-gonic/gin"
	// services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/internal/mocks"
	"github.com/ksel172/Meduza/teamserver/internal/storage"

	// "github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/ksel172/Meduza/teamserver/utils"
	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Requirements
	mockAgentDAL := new(mocks.MockAgentDAL)
	controller := NewCheckInController(mockAgentDAL)

	// Test agent auth token
	testAuthToken := "test-auth-token"

	// Create a mock public key for the agent and server
	_, agentPubKeyBytes, err := utils.GenerateECDHKeyPair()
	if err != nil {
		t.Fatal("failed to generate agent ecdh key pair")
	}
	testAgentPubKey := string(agentPubKeyBytes)

	serverPrivKey, serverPubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		t.Fatal("failed to generate server ecdh key pair")
	}

	// Prepare the key store for use on this test
	storage.AsymmetricKeyRegistry.WriteKey(testAuthToken, storage.KeyPair{
		PublicKey:  serverPubKey,
		PrivateKey: serverPrivKey,
	})

	tests := []struct {
		name        string
		agentPubKey string
		authToken   string
		expectError bool
	}{
		{
			name:        "agent authentication: success",
			agentPubKey: testAgentPubKey,
			authToken:   testAuthToken,
			expectError: false,
		},
		{
			name:        "agent authentication: invalid auth token",
			agentPubKey: testAgentPubKey,
			authToken:   "invalid-auth-token",
			expectError: true,
		},
		{ // This must be the last test in the grid
			name:        "agent authentication: key not in registry",
			agentPubKey: testAgentPubKey,
			authToken:   testAuthToken,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare mock DAL calls in order
			switch tt.name {
			case "agent authentication: key not in registry":
				storage.AsymmetricKeyRegistry.DeleteKey(testAuthToken)
			}

			// Submit request
			response, err := controller.Authenticate(tt.agentPubKey, tt.authToken)

			// Check error
			if !tt.expectError {
				assert.Nil(t, err)
			}

			// Verify returned values
			switch tt.name {
			case "agent authentication: success":
				assert.Equal(t, serverPubKey, response.PublicKey)
			}
		})
	}
}

/*
Below here could be reused for HTTP Listener tests
Looks more like that than checkInController functionality currently
*/
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
