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
