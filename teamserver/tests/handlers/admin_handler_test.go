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
	"github.com/joho/godotenv"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateAdmin(t *testing.T) {
	mockAdminDAL := &mocks.MockAdminDal{}
	handler := handlers.NewAdminController(mockAdminDAL)
	gin.SetMode(gin.TestMode)

	godotenv.Load("../../.env")
	adminSecret := conf.GetMeduzaAdminSecret()
	adminReq := models.ResAdmin{
		Adminname:    "testAdmin",
		PasswordHash: "test-admin-password",
	}

	tests := []struct {
		name           string
		adminRequest   models.ResAdmin
		mockError      error
		expectedStatus int
	}{
		{
			name:           "succesful create admin",
			adminRequest:   adminReq,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "dal error",
			adminRequest:   adminReq,
			mockError:      errors.New("failed dal op"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockAdminDAL.On("CreateDefaultAdmins", mock.AnythingOfType("*models.ResAdmin")).Return(tt.mockError).Once()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.adminRequest)
			c.Request = httptest.NewRequest(http.MethodGet, "/", bytes.NewReader(body))
			c.Request.Header = http.Header{
				"Authorization": []string{fmt.Sprintf("Bearer %s", adminSecret)},
			}
			handler.CreateAdmin(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockAdminDAL.AssertExpectations(t)
		})
	}
}
