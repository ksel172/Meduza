package listener_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/ksel172/Meduza/teamserver/internal/storage"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
	"github.com/stretchr/testify/require"
)

type TestAgent interface {
	Authenticate(*testing.T, string) error
}

type TestHTTPAgent struct {
	models.Agent
	storage.KeyPair
	authToken    string
	sharedAesKey []byte
	sessionToken []byte
	callbackURL  string
	client       *http.Client
}

func newTestHTTPAgent(host string, port int, token string) (TestHTTPAgent, error) {
	// Create test agent to send requests to listener
	privKey, pubKey, err := utils.GenerateECDHKeyPair()
	if err != nil {
		return TestHTTPAgent{}, fmt.Errorf("failed to generate server keys: %v", err)
	}

	return TestHTTPAgent{
		KeyPair: storage.KeyPair{
			PublicKey:  pubKey,
			PrivateKey: privKey,
		},
		authToken:   token,
		callbackURL: fmt.Sprintf("http://%s:%d", host, port),
		client:      &http.Client{},
	}, nil
}

// TODO: prepare key store with server public and private keys
// also retrieve the agent keys to send with c2request message
func (a *TestHTTPAgent) Authenticate(t *testing.T, listenerID string) {
	// Prepare request body by encoding the c2request as base64 bytes
	c2request := models.C2Request{
		Message: base64.StdEncoding.EncodeToString(a.PublicKey),
	}
	c2requestBytes, err := json.Marshal(c2request)
	if err != nil {
		t.Fatalf("failed to marshal c2request: %v", err)
	}
	requestBody := make([]byte, base64.StdEncoding.EncodedLen(len(c2requestBytes)))
	base64.StdEncoding.Encode(requestBody, c2requestBytes)
	t.Logf("Prepared c2request: %+v", c2request)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, a.callbackURL+"/", bytes.NewReader(requestBody))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("Auth-Token", base64.StdEncoding.EncodeToString([]byte(a.authToken)))
	req.Header.Add("Content-Type", "application/json")

	// Make request
	// fmt.Printf("1. Test - BASE 64 Agent Public Key: %s\n", c2request.Message)
	resp, err := a.client.Do(req)
	if err != nil {
		t.Fatalf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("unexpected status code, expected: %d, got: %d", http.StatusAccepted, resp.StatusCode)
	}

	// Read body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Parse response
	var response struct {
		PublicKey    string `json:"public_key"`
		SessionToken string `json:"session_token"`
	}
	err = json.Unmarshal(respBody, &response)
	require.NoError(t, err, "failed to parse response body")

	t.Logf("server authentication response: %+v", response)

	// Save session token
	a.sessionToken, err = base64.StdEncoding.DecodeString(response.SessionToken)
	if err != nil {
		t.Fatalf("failed to decode session token: %v", err)
	}

	// Derive shared key for encrypting follow up requests
	decodedServerPublicKey, err := base64.StdEncoding.DecodeString(response.PublicKey)
	if err != nil {
		t.Fatalf("failed to decode server public key: %v", err)
	}
	a.sharedAesKey, err = utils.DeriveECDHSharedSecret(a.PrivateKey, decodedServerPublicKey)
	if err != nil {
		t.Fatalf("failed to derive shared key: %v", err)
	}
}

func (a *TestHTTPAgent) Register(t *testing.T) {
	c2request := models.C2Request{
		Reason: models.Register,
		Message: fmt.Sprintf(`{"hostname":"%s", "ip_address":"%s", "username":"%s", "system_info":"%s", "os_info":"%s"}`,
			"test-host", "192.168.0.1", "test-username", "test-system-info", "test-os-info"),
	}
	c2requestBytes, err := json.Marshal(c2request)
	if err != nil {
		t.Fatalf("failed to marshal c2request: %v", err)
	}

	// Encrypt c2 request
	c2requestEncrypted, err := utils.AesEncrypt(a.sharedAesKey, c2requestBytes)
	if err != nil {
		t.Fatalf("failed to encrypt c2request: %v", err)
	}
	t.Logf("Prepared encrypted c2request: %+v", c2request)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, a.callbackURL+"/", bytes.NewReader(c2requestEncrypted))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("Session-Token", base64.StdEncoding.EncodeToString(a.sessionToken))
	req.Header.Add("Auth-Token", base64.StdEncoding.EncodeToString([]byte(a.authToken)))
	req.Header.Add("Content-Type", "application/json")

	// Make request
	resp, err := a.client.Do(req)
	if err != nil {
		t.Fatalf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status code, expected: %d, got: %d", http.StatusCreated, resp.StatusCode)
	}

	// The server returns the agent ID
	var agentID string
	if err := json.NewDecoder(resp.Body).Decode(&agentID); err != nil {
		t.Fatalf("failed to decode agent ID from response: %v", err)
	}
	t.Logf("Received agent ID = %s", agentID)

	// Temporarily store the ID
	a.ID = agentID
}

func (a *TestHTTPAgent) GetTasks(t *testing.T) {
	c2request := models.C2Request{
		AgentID: a.ID,
		Reason:  models.Task,
	}
	c2requestBytes, err := json.Marshal(c2request)
	if err != nil {
		t.Fatalf("failed to marshal c2request: %v", err)
	}
	c2requestEncrypted, err := utils.AesEncrypt(a.sharedAesKey, c2requestBytes)
	if err != nil {
		t.Fatalf("failed to encrypt c2request: %v", err)
	}
	t.Logf("Prepared encrypted c2request: %+v", c2request)

	req, err := http.NewRequest(http.MethodPost, a.callbackURL+"/", bytes.NewReader(c2requestEncrypted))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}
	req.Header.Add("Session-Token", base64.StdEncoding.EncodeToString(a.sessionToken))
	req.Header.Add("Content-Type", "application/json")

	resp, err := a.client.Do(req)
	if err != nil {
		t.Fatalf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code, expected: %d, got: %d", http.StatusOK, resp.StatusCode)
	}
}
