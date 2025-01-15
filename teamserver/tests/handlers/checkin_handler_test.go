package handler_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	services "github.com/ksel172/Meduza/teamserver/internal/services/listeners"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAgent(t *testing.T) {
	mockAgentDAL := new(mocks.MockAgentDAL)
	mockCheckInDal := new(mocks.MockCheckInDal)
	handler := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	gin.SetMode(gin.TestMode)

	// Create c2 request and give it an UUID
	c2Request := models.NewC2Request()
	c2Request.AgentID = uuid.New().String()

	tests := []struct {
		name           string
		c2Request      models.C2Request
		mockError      error
		reachDAL       bool
		expectedStatus int
	}{
		{
			name:           "successful create agent",
			c2Request:      c2Request,
			mockError:      nil,
			reachDAL:       true,
			expectedStatus: http.StatusCreated,
		},
		{ // Will get parsed from JSON but fail on validation check
			name:           "non-uuid c2request ID",
			c2Request:      models.NewC2Request(),
			mockError:      nil,
			reachDAL:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			c2Request:      c2Request,
			mockError:      errors.New("failed dal op"),
			reachDAL:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reachDAL {
				mockCheckInDal.On("CreateAgent", mock.AnythingOfType("models.Agent")).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.c2Request)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

			handler.Checkin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInDal.AssertExpectations(t)
		})
	}
}

func TestGetTasks(t *testing.T) {
	mockAgentDAL := new(mocks.MockAgentDAL)
	mockCheckInDal := new(mocks.MockCheckInDal)
	handler := services.NewCheckInController(mockCheckInDal, mockAgentDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		agentID        string
		mockError      error
		reachDAL       bool
		expectedStatus int
	}{
		{
			name:           "successful create agent",
			agentID:        "test-agent-id",
			mockError:      nil,
			reachDAL:       true,
			expectedStatus: http.StatusOK,
		},
		{ // Will get parsed from JSON but fail on validation check
			name:           "missing agent id",
			agentID:        "",
			mockError:      nil,
			reachDAL:       false,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			agentID:        "test-agent-id",
			mockError:      errors.New("failed dal op"),
			reachDAL:       true,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reachDAL {
				mockAgentDAL.On("GetAgentTasks", tt.agentID).Return(make([]models.AgentTask, 0), tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: models.ParamAgentID, Value: tt.agentID}}

			handler.Checkin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInDal.AssertExpectations(t)
		})
	}
}
