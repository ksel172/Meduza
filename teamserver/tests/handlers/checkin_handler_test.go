package handler_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAgentAuthController(t *testing.T) {
	mockPayloadDal := new(mocks.MockPayloadDAL)

	// checkInController := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	agentAuthController := services.NewAgentAuthController(mockPayloadDal)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		c2Request      models.C2Request
		expectedStatus int
	}{
		{
			name:           "agent auth: success",
			c2Request:      models.C2Request{ConfigID: "test-config-id"},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "agent auth: fail get payload token",
			c2Request:      models.NewC2Request(),
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "agent auth: stored token does not match provided token",
			c2Request:      models.NewC2Request(),
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			switch tt.name {
			case "agent auth: success":
				mockPayloadDal.On("GetPayloadToken", mock.AnythingOfType("string")).Return("testPayloadToken", nil).Once()
			case "agent auth: fail get payload token":
				mockPayloadDal.On("GetPayloadToken", mock.AnythingOfType("string")).Return("", errors.New("failed")).Once()
			case "agent auth: stored token does not match provided token":
				mockPayloadDal.On("GetPayloadToken", mock.AnythingOfType("string")).Return("failTestPayloadToken", nil).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.c2Request)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Request.Header.Set("Authorization", "testPayloadToken")

			// If Authorization header is provided, then request will be handled at encryption key request
			// Otherwise, authentication will be skipped and agent/server communication will already be encrypted
			agentAuthController.AuthenticateAgent(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockPayloadDal.AssertExpectations(t)
		})
	}
}

// From here below, all tests will assume agent authentication handler has already been called and the values have been set in the context
// This is validated by test TestAgentAuthController

func TestAgentEncryptionKeyRequest(t *testing.T) {
	mockCheckInDal := new(mocks.MockCheckInDal)
	mockAgentDAL := new(mocks.MockAgentDAL)

	// checkInController := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	checkInController := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		c2Request      models.C2Request
		expectedStatus int
	}{
		{
			name:           "encryption key request: success",
			c2Request:      models.C2Request{ConfigID: "test-config-id"},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.c2Request)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

			// Simulate the agentAuth handler, it sets 2 values in the request as below
			c.Set(services.AuthToken, "testPayloadToken")
			c.Set("c2request", tt.c2Request)

			checkInController.Checkin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestAgentRegisterRequest(t *testing.T) {
	mockAgentDAL := new(mocks.MockAgentDAL)
	mockCheckInDal := new(mocks.MockCheckInDal)
	controller := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	gin.SetMode(gin.TestMode)

	c2request := models.C2Request{
		Reason:  models.Register,
		AgentID: "test-agent-id",
		Message: `{"agent_id": "test-agent-id"}`,
	}

	tests := []struct {
		name           string
		c2Request      models.C2Request
		expectedStatus int
	}{
		{
			name:           "register agent: success",
			c2Request:      c2request,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "register agent: agent already exists",
			c2Request:      c2request,
			expectedStatus: http.StatusConflict,
		},
		{
			name:           "register agent: create agent error",
			c2Request:      c2request,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:           "register agent: create agent info error",
			c2Request:      c2request,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "register agent: success":
				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, errors.New("agent does not exist")).Once()
				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(nil).Once()
				mockAgentDAL.On("CreateAgentInfo", mock.AnythingOfType("models.AgentInfo")).Return(nil).Once()
			case "register agent: agent already exists":
				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, nil).Once()
			case "register agent: create agent error":
				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, errors.New("agent does not exist")).Once()
				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(errors.New("failed to create agent")).Once()
			case "register agent: create agent info error":
				mockAgentDAL.On("GetAgent", tt.c2Request.AgentID).Return(models.Agent{}, errors.New("agent does not exist")).Once()
				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(nil).Once()
				mockAgentDAL.On("CreateAgentInfo", mock.AnythingOfType("models.AgentInfo")).Return(errors.New("failed to create agent info")).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.c2Request)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

			// Simulate the agentAuth handler, it sets the c2request in the request as below
			c.Set("c2request", tt.c2Request)

			controller.Checkin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInDal.AssertExpectations(t)
			mockAgentDAL.AssertExpectations(t)
		})
	}
}

func TestAgentTasksRequest(t *testing.T) {
	mockAgentDAL := new(mocks.MockAgentDAL)
	mockCheckInDal := new(mocks.MockCheckInDal)
	handler := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	gin.SetMode(gin.TestMode)

	c2request := models.C2Request{
		Reason:  models.Task,
		AgentID: "test-agent-id",
		Message: `{"agent_id": "test-agent-id"}`,
	}

	tests := []struct {
		name           string
		c2request      models.C2Request
		expectedStatus int
	}{
		{
			name:           "agent task: success",
			c2request:      c2request,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "agent task: get agent tasks error",
			c2request:      c2request,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "agent task: callback update error",
			c2request:      c2request,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "agent task: success":
				mockAgentDAL.On("GetAgentTasks", tt.c2request.AgentID).Return(make([]models.AgentTask, 0), nil).Once()
				mockAgentDAL.On("UpdateAgentLastCallback", tt.c2request.AgentID, mock.AnythingOfType("string")).Return(nil).Once()
			case "agent task: get agent tasks error":
				mockAgentDAL.On("GetAgentTasks", tt.c2request.AgentID).Return(make([]models.AgentTask, 0), errors.New("failed to get agent tasks")).Once()
			case "agent task: callback update error":
				mockAgentDAL.On("GetAgentTasks", tt.c2request.AgentID).Return(make([]models.AgentTask, 0), nil).Once()
				mockAgentDAL.On("UpdateAgentLastCallback", tt.c2request.AgentID, mock.AnythingOfType("string")).Return(errors.New("failed to update agent last callback")).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.c2request)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Set("c2request", tt.c2request)

			handler.Checkin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInDal.AssertExpectations(t)
			mockAgentDAL.AssertExpectations(t)
		})
	}
}

func TestAgentResponseRequest(t *testing.T) {
	mockAgentDAL := new(mocks.MockAgentDAL)
	mockCheckInDal := new(mocks.MockCheckInDal)
	handler := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	gin.SetMode(gin.TestMode)

	c2request := models.C2Request{
		Reason:  models.Response,
		AgentID: "test-agent-id",
		Message: `{"agent_id": "test-agent-id"}`,
	}

	tests := []struct {
		name           string
		c2request      models.C2Request
		expectedStatus int
	}{
		{
			name:           "agent response: success",
			c2request:      c2request,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "agent response: update agent task error",
			c2request:      c2request,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.name {
			case "agent response: success":
				mockAgentDAL.On("UpdateAgentTask", mock.AnythingOfType("models.AgentTask")).Return(nil).Once()
			case "agent response: update agent task error":
				mockAgentDAL.On("UpdateAgentTask", mock.AnythingOfType("models.AgentTask")).Return(errors.New("failed to update agent task")).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.c2request)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Set("c2request", tt.c2request)

			handler.Checkin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInDal.AssertExpectations(t)
			mockAgentDAL.AssertExpectations(t)
		})
	}
}
