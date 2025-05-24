package http_listener

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"testing"
	"time"

	controller_mocks "github.com/ksel172/Meduza/teamserver/internal/mocks/controller"
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
func createAgentSession(sessionToken string) ([]byte, error) {
	_, agentPublicKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate agent keys: %v", err)
	}
	serverPrivKey, _, err := utils.GenerateECDHKeyPair()
	if err != nil {
		return nil, fmt.Errorf("failed to generate agent keys: %v", err)
	}
	aesKey, err := utils.DeriveECDHSharedSecret(serverPrivKey, agentPublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to derive shared key: %v", err)
	}
	storage.KeyRegistry.WriteKey(sessionToken, aesKey)

	return aesKey, nil
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
	return encryptedc2request, nil
}
