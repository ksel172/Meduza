package listener_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
	PayloadController  *handlers.PayloadController
	AgentController    *handlers.AgentController
	ListenerService    *listener_service.ListenerService
}

func setup(t *testing.T, createLocalListenerRequest models.CreateLocalListenerRequest) (Container, models.Listener, string, error) {
	// Connect to database, prepare dependencies container
	t.Log("creating dependencies container")
	container := createContainer(t)

	// Seed db - 1. Listener
	switch createLocalListenerRequest.Kind {
	case models.HTTPListenerKind:
		createHTTPListener(t, container, createLocalListenerRequest)
	default:
		t.Fatalf("listener kind not yet implemented")
	}
	t.Log("created listener in the database succesfully")

	// Retrive the listener that was created
	createdListener := getListener(t, container, createLocalListenerRequest.Name)
	t.Logf("retrieved listener from database: %+v", createdListener)

	// Seed db - 2. Payload
	createPayload(t, container, models.PayloadRequest{
		PayloadName: "test-payload",
		ListenerID:  createdListener.ID,
	})
	t.Log("craeted payload in the database")

	createdPayload := getPayload(t, container)
	t.Logf("retrieved payload from database: %+v", createdPayload)

	authToken := getPayloadToken(t, container, createdPayload.ConfigID)

	return container, createdListener, authToken, nil
}

func createContainer(t *testing.T) Container {
	// Setup db
	schema := conf.GetMeduzaDbSchema()
	db, err := repos.Setup()
	if err != nil {
		t.Fatalf("failed to setup database: %v", err)
	}

	// Create required data layers
	listenerDAL := dal.NewListenerDAL(db, schema)
	agentDAL := dal.NewAgentDAL(db, schema)
	payloadDAL := dal.NewPayloadDAL(db, schema)
	moduleDAL := dal.NewModuleDAL(db, schema)

	// Create services
	listenerService := listener_service.NewListenerService(listenerDAL)

	// Create controllers
	listenerController := handlers.NewListenersHandler(listenerService, listenerDAL)
	payloadController := handlers.NewPayloadController(agentDAL, listenerDAL, payloadDAL)
	agentController := handlers.NewAgentController(agentDAL, moduleDAL)

	return Container{
		ListenerController: listenerController,
		PayloadController:  payloadController,
		AgentController:    agentController,
		ListenerService:    listenerService,
	}
}

// Targets C2 server client API and creates a listener in the database
// bypasses launching the server and directly hits the endpoint, as the server is NOT launched in these e2e tests
func createHTTPListener(t *testing.T, container Container, createLocalListenerRequest models.CreateLocalListenerRequest) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Marshal the listener model into the request body
	body, err := json.Marshal(createLocalListenerRequest)
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
	t.Logf("CreateListener response: %+v", response)

	// Ensure it was created
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 CREATED response")
}

// Retrieves a listener from the database by its name
func getListener(t *testing.T, container Container, listenerName string) models.Listener {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Make request
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{{Key: models.ParamListenerName, Value: listenerName}}
	container.ListenerController.GetListenerByName(c)

	require.Equal(t, http.StatusOK, w.Code, "expected 200 OK from GetAllListeners")

	// Define response wrapper
	var response struct {
		Status  int             `json:"status"`
		Message string          `json:"message"`
		Data    models.Listener `json:"data"`
	}

	// Parse response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "failed to parse response body")

	return response.Data
}

func createPayload(t *testing.T, container Container, payloadRequest models.PayloadRequest) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Marshal the listener model into the request body
	body, err := json.Marshal(payloadRequest)
	if err != nil {
		t.Fatalf("failed to marshal payload request")
	}

	// Create the request
	c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	// Create the listener by sending a request to the controller
	container.PayloadController.CreatePayload(c)

	// Define response wrapper
	var response struct {
		Status  int               `json:"status"`
		Message string            `json:"message"`
		Data    []models.Listener `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)
	t.Logf("CreatePayload response: %+v", response)

	// Ensure it was created
	require.Equal(t, http.StatusCreated, w.Code, "expected 201 CREATED response")
}

// Return the first payload for now
func getPayload(t *testing.T, container Container) models.PayloadConfig {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Make request
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	container.PayloadController.GetAllPayloads(c)

	require.Equal(t, http.StatusOK, w.Code, "expected 200 OK from GetAllListeners")

	// Define response wrapper
	var response struct {
		Status  int                    `json:"status"`
		Message string                 `json:"message"`
		Data    []models.PayloadConfig `json:"data"`
	}

	// Parse response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "failed to parse response body")

	return response.Data[0]
}

func getPayloadToken(t *testing.T, container Container, payloadID string) string {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	// Make request
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Params = gin.Params{{Key: models.ParamPayloadID, Value: payloadID}}
	container.PayloadController.GetToken(c)

	require.Equal(t, http.StatusOK, w.Code, "expected 200 OK from GetAllListeners")

	// Define response wrapper
	var response struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
		Data    string `json:"data"`
	}

	// Parse response
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err, "failed to parse response body")

	return response.Data
}
