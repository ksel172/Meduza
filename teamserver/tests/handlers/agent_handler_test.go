package handler_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAgent(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		agentID        string
		mockAgent      models.Agent
		mockError      error
		expectedStatus int
	}{
		{
			name:    "successful get agent",
			agentID: "test-agent-id",
			mockAgent: models.Agent{
				AgentID: "test-agent-id",
				Name:    "test-agent",
				Status:  "active",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing agent id",
			agentID:        "",
			mockAgent:      models.Agent{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "agent not found",
			agentID:        "non-existent",
			mockAgent:      models.Agent{},
			mockError:      fmt.Errorf("agent not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agentID != "" {
				mockDAL.On("GetAgent", tt.agentID).Return(tt.mockAgent, tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: models.ParamAgentID, Value: tt.agentID}}

			handler.GetAgent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			if tt.expectedStatus == http.StatusOK {
				var response models.Agent
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockAgent, response)
			}
			mockDAL.AssertExpectations(t)
		})
	}
}

func TestUpdateAgent(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	// Below agent is sent as JSON to the handler
	agentUpdateRequest := models.UpdateAgentRequest{
		AgentID: "test-agent-id",
		Name:    "updated-agent-name",
	}

	// Handler returns the below agent from db
	updatedAgent := models.Agent{
		AgentID: "test-agent-id",
		Name:    "updated-agent-name",
	}

	tests := []struct {
		name               string
		agentUpdateRequest models.UpdateAgentRequest
		updatedAgent       models.Agent
		mockError          error
		expectedStatus     int
	}{
		{
			name:               "successful agent update",
			agentUpdateRequest: agentUpdateRequest,
			updatedAgent:       updatedAgent,
			mockError:          nil,
			expectedStatus:     http.StatusOK,
		},
		{
			name:               "agent update server error",
			agentUpdateRequest: agentUpdateRequest,
			updatedAgent:       updatedAgent,
			mockError:          errors.New("example failed dal op"),
			expectedStatus:     http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// If functionality of the update handler is extended in the future, need to make sure to
			// only setup the mockDAL on test cases that actually reach the dal
			mockDAL.On("UpdateAgent", mock.AnythingOfType("models.UpdateAgentRequest")).Return(tt.updatedAgent, tt.mockError).Once()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			body, _ := json.Marshal(tt.agentUpdateRequest)
			c.Request, _ = http.NewRequest(http.MethodPut, "/", bytes.NewReader(body))

			handler.UpdateAgent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}

func TestDeleteAgent(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		agentID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful delete agent",
			agentID:        "test-agent-id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing agent id",
			agentID:        "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			agentID:        "non-existent",
			mockError:      fmt.Errorf("failed dal op"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agentID != "" {
				mockDAL.On("DeleteAgent", tt.agentID).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: models.ParamAgentID, Value: tt.agentID}}

			handler.DeleteAgent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}

/* Agent tasks tests */

func TestGetAgentTasks(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	tasks := []models.AgentTask{
		{AgentID: "found-agent-task-id"},
	}

	tests := []struct {
		name           string
		agentID        string
		foundTasks     []models.AgentTask
		mockError      error
		reachDAL       bool
		expectedStatus int
	}{
		{
			name:           "successful get agent tasks",
			agentID:        "test-agent-id",
			foundTasks:     tasks,
			mockError:      nil,
			reachDAL:       true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing agent id",
			agentID:        "",
			mockError:      nil,
			reachDAL:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			agentID:        "test-agent-id",
			foundTasks:     tasks,
			mockError:      errors.New("failed dal op"),
			reachDAL:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reachDAL {
				mockDAL.On("GetAgentTasks", tt.agentID).Return(tt.foundTasks, tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: models.ParamAgentID, Value: tt.agentID}}

			handler.GetAgentTasks(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}

func TestCreateAgentTask(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		agentID        string
		taskRequest    models.AgentTaskRequest
		mockError      error
		expectedStatus int
	}{
		{
			name:    "successful task creation",
			agentID: "test-agent-id",
			taskRequest: models.AgentTaskRequest{
				Type:    "command",
				Command: "whoami",
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing agent id",
			agentID:        "",
			taskRequest:    models.AgentTaskRequest{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "dal error",
			agentID: "test-agent-id",
			taskRequest: models.AgentTaskRequest{
				Type:    "command",
				Command: "whoami",
			},
			mockError:      fmt.Errorf("dal error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agentID != "" {
				mockDAL.On("CreateAgentTask", mock.AnythingOfType("models.AgentTask")).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "agent_id", Value: tt.agentID}}

			body, _ := json.Marshal(tt.taskRequest)
			c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
			c.Request.Header.Set("Content-Type", "application/json")

			handler.CreateAgentTask(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}
func TestDeleteAgentTasks(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		agentID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful tasks deletion",
			agentID:        "test-agent-id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing agent id",
			agentID:        "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			agentID:        "test-agent-id",
			mockError:      fmt.Errorf("dal error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agentID != "" {
				mockDAL.On("DeleteAgentTasks", tt.agentID).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: models.ParamAgentID, Value: tt.agentID}}

			handler.DeleteAgentTasks(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}

func TestDeleteAgentTask(t *testing.T) {
	mockDAL := new(mocks.MockAgentDAL)
	handler := handlers.NewAgentController(mockDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		agentID        string
		taskID         string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful tasks deletion",
			agentID:        "test-agent-id",
			taskID:         "test-task-id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "missing agent id",
			agentID:        "",
			taskID:         "test-task-id",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "missing task id",
			agentID:        "test-agent-id",
			taskID:         "",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			agentID:        "test-agent-id",
			taskID:         "test-task-id",
			mockError:      fmt.Errorf("dal error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.agentID != "" && tt.taskID != "" {
				mockDAL.On("DeleteAgentTask", tt.agentID, tt.taskID).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{
				{Key: models.ParamAgentID, Value: tt.agentID},
				{Key: models.ParamTaskID, Value: tt.taskID},
			}

			handler.DeleteAgentTask(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}
