package http_listener

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	controller_mocks "github.com/ksel172/Meduza/teamserver/internal/mocks/controller"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

// Helper function to build default test listener
func buildDefaultTestHTTPListener(t *testing.T) *HTTPListener {
	mockCheckInController := &controller_mocks.MockCheckInController{}

	listener, err := NewHTTPListener(
		HTTPListenerConfig{
			Host:         "127.0.0.1",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		mockCheckInController,
	)
	if err != nil {
		t.Fatalf("failed to build http listener: %v", err)
	}

	return listener
}

// Helper function to build default test listener
func buildDefaultTestHTTPListenerWithMock(t *testing.T, mockCheckInController *controller_mocks.MockCheckInController) *HTTPListener {

	listener, err := NewHTTPListener(
		HTTPListenerConfig{
			Host:         "127.0.0.1",
			Port:         8080,
			ReadTimeout:  30,
			WriteTimeout: 30,
		},
		mockCheckInController,
	)
	if err != nil {
		t.Fatalf("failed to build http listener: %v", err)
	}

	return listener
}

func buildStartDefaultTestHTTPListener(t *testing.T) *HTTPListener {
	listener := buildDefaultTestHTTPListener(t)
	if err := listener.Start(context.Background()); err != nil {
		t.Fatalf("failed to start http listener: %v", err)
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

// Function to simulate agent authentication and retrieve sessiontoken + AES key for message encryption
func authenticateAgent(authToken string) (string, []byte, error) {
	controller := checkin.CheckInController{}

	// Write into the registry the server keys for that authToken
	serverPrivKey, serverPubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate server keys: %v", err)
	}
	storage.AsymmetricKeyRegistry.WriteKey(authToken, storage.KeyPair{
		PublicKey:  serverPubKey,
		PrivateKey: serverPrivKey,
	})

	// Generate agent keys
	agentPrivKey, agentPubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		return "", nil, fmt.Errorf("failed to generate agent keys: %v", err)
	}

	// Send authenticate request to check in controller
	response, err := controller.Authenticate(string(agentPubKey), authToken)
	if err != nil {
		return "", nil, errors.New("failed to authenticate")
	}

	// Use the server public key to derive the shared key
	sharedKey, err := utils.DeriveECDHSharedSecret(agentPrivKey, response.PublicKey)
	if err != nil {
		return "", nil, fmt.Errorf("failed to derive shared key")
	}

	// SessionToken is sent in base64 in the requests anyway
	return response.SessionToken, sharedKey, nil
}

func encryptAgentRequest(c2request models.C2Request, key []byte) ([]byte, error) {
	c2requestBytes, err := json.Marshal(c2request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal c2request")
	}
	encryptedc2request, err := utils.AesEncrypt(key, c2requestBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt message")
	}
	return []byte(base64.StdEncoding.EncodeToString(encryptedc2request)), nil
}
