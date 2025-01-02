package handler_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUsers(t *testing.T) {
	mockUserDAL := &mocks.MockUserDAL{}
	handler := handlers.NewUserController(mockUserDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful get users",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "dal error",
			mockError:      errors.New("failed dal op"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserDAL.On("GetUsers").Return(make([]models.User, 0), tt.mockError).Once()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

			handler.GetUsers(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUserDAL.AssertExpectations(t)
		})
	}
}

func TestAddUsers(t *testing.T) {
	mockUserDAL := &mocks.MockUserDAL{}
	handler := handlers.NewUserController(mockUserDAL)
	gin.SetMode(gin.TestMode)

	requestUser := models.ResUser{
		ID:           "test-user-id",
		Username:     "testusername",
		PasswordHash: "test-password",
		Role:         "visitor",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	tests := []struct {
		name           string
		requestUser    models.ResUser
		mockError      error
		reachDAL       bool
		expectedStatus int
	}{
		{
			name:           "successful get users",
			requestUser:    requestUser,
			mockError:      nil,
			reachDAL:       true,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "dal error",
			requestUser:    requestUser,
			mockError:      errors.New("failed dal op"),
			reachDAL:       true,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reachDAL {
				mockUserDAL.On("AddUsers", mock.AnythingOfType("*models.ResUser")).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.requestUser)
			c.Request = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))

			handler.AddUsers(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUserDAL.AssertExpectations(t)
		})
	}
}
