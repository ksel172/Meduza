package handler_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAgentDAL struct {
	mock.Mock
}

func (m *MockAgentDAL) GetAgent(agentID string) (models.Agent, error) {
	args := m.Called(agentID)
	return args.Get(0).(models.Agent), args.Error(1)
}

func (m *MockAgentDAL) UpdateAgent(ctx context.Context, agent models.Agent) error {
	args := m.Called(ctx, agent)
	return args.Error(0)
}

func (m *MockAgentDAL) DeleteAgent(ctx context.Context, agentID string) error {
	args := m.Called(ctx, agentID)
	return args.Error(0)
}

func (m *MockAgentDAL) CreateAgentTask(ctx context.Context, task models.AgentTask) error {
	args := m.Called(ctx, task)
	return args.Error(0)
}

func (m *MockAgentDAL) GetAgentTasks(ctx context.Context, agentID string) ([]models.AgentTask, error) {
	args := m.Called(ctx, agentID)
	return args.Get(0).([]models.AgentTask), args.Error(1)
}

func (m *MockAgentDAL) DeleteAgentTask(ctx context.Context, agentID string, taskID string) error {
	args := m.Called(ctx, agentID, taskID)
	return args.Error(0)
}

func (m *MockAgentDAL) DeleteAgentTasks(ctx context.Context, agentID string) error {
	args := m.Called(ctx, agentID)
	return args.Error(0)
}

func TestGetAgent(t *testing.T) {
	mockDAL := new(MockAgentDAL)
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
				ID:     "test-agent-id",
				Name:   "test-agent",
				Status: "active",
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
			c.Params = gin.Params{{Key: "agent_id", Value: tt.agentID}}

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

func TestCreateAgentTask(t *testing.T) {
	mockDAL := new(MockAgentDAL)
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
				mockDAL.On("CreateAgentTask", mock.Anything, mock.AnythingOfType("models.AgentTask")).Return(tt.mockError).Once()
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
	mockDAL := new(MockAgentDAL)
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
				mockDAL.On("DeleteAgentTasks", mock.Anything, tt.agentID).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: "agent_id", Value: tt.agentID}}

			handler.DeleteAgentTasks(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockDAL.AssertExpectations(t)
		})
	}
}
