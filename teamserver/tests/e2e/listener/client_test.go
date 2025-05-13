package listener_test

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/require"
)

type TestAgent interface {
	Authenticate() error
}

type TestHTTPAgent struct {
	models.Agent
	AuthToken    string
	SessionToken string
	CallbackURL  string
	client       *http.Client
}

func newTestHTTPAgent(host string, port int, token string) TestHTTPAgent {
	return TestHTTPAgent{
		Agent: models.Agent{
			AgentID: uuid.New().String(),
		},
		AuthToken:   token,
		CallbackURL: fmt.Sprintf("http://%s:%d", host, port),
		client:      &http.Client{},
	}
}

// TODO: prepare key store with server public and private keys
// also retrieve the agent keys to send with c2request message
func (a *TestHTTPAgent) Authenticate(t *testing.T, listenerID string) {
	// Prepare request body
	c2request := models.C2Request{
		AgentID: a.AgentID,
	}
	c2requestBytes, err := json.Marshal(c2request)
	if err != nil {
		t.Fatalf("failed to marshal c2request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest(http.MethodPost, a.CallbackURL+"/", bytes.NewReader(c2requestBytes))
	if err != nil {
		t.Fatalf("failed to create HTTP request: %v", err)
	}

	req.Header.Add("Auth-Token", base64.StdEncoding.EncodeToString([]byte(a.AuthToken)))
	req.Header.Set("Content-Type", "application/json")

	// Make request
	resp, err := a.client.Do(req)
	if err != nil {
		t.Fatalf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Read body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	// Parse response
	var response struct {
		PublicKey    int    `json:"public_key"`
		SessionToken string `json:"session_token"`
	}
	err = json.Unmarshal(body, &response)
	require.NoError(t, err, "failed to parse response body")

	t.Logf("response: %+v", response)

	// Save session token
	a.SessionToken = response.SessionToken
}
