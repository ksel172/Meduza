package handler_tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/handlers"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/tests/mocks"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	mockUserDAL := &mocks.MockUserDAL{}
	mockAuthProvider := &mocks.MockJWTService{}
	handler := handlers.NewAuthController(mockUserDAL, mockAuthProvider)
	gin.SetMode(gin.TestMode)

	loginR := models.AuthRequest{
		Username: "test-username",
		Password: "test-password",
	}
	userFound := models.ResUser{
		ID:           "test-user-id",
		PasswordHash: "$2a$12$xrf0BNvHN7VUHGBKMZ.VHu9WKtdNgBmPlq3xfRQeBcUzjVFBo.QMq",
	}
	incorrectPasswordUser := models.ResUser{
		ID:           "test-user-id",
		PasswordHash: "test-incorrect-password",
	}

	tests := []struct {
		name             string
		loginRequest     models.AuthRequest
		userFound        models.ResUser
		mockUserDALError error
		mockAuthError    error
		reachUserDAL     bool
		reachJWTService  bool
		expectedStatus   int
	}{
		{
			name:             "successful logout",
			loginRequest:     loginR,
			userFound:        userFound,
			mockUserDALError: nil,
			mockAuthError:    nil,
			reachUserDAL:     true,
			reachJWTService:  true,
			expectedStatus:   http.StatusOK,
		},
		{
			name:             "dal error",
			loginRequest:     loginR,
			userFound:        userFound,
			mockUserDALError: errors.New("failed dal op"),
			mockAuthError:    nil,
			reachUserDAL:     true,
			reachJWTService:  false,
			expectedStatus:   http.StatusBadRequest,
		},
		{
			name:             "incorrect password",
			loginRequest:     loginR,
			userFound:        incorrectPasswordUser,
			mockUserDALError: nil,
			mockAuthError:    nil,
			reachUserDAL:     true,
			reachJWTService:  false,
			expectedStatus:   http.StatusUnauthorized,
		},
		{
			name:             "auth error",
			loginRequest:     loginR,
			userFound:        userFound,
			mockUserDALError: nil,
			mockAuthError:    errors.New("failed auth op"),
			reachUserDAL:     true,
			reachJWTService:  true,
			expectedStatus:   http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.reachUserDAL {
				mockUserDAL.On("GetUserByUsername", tt.loginRequest.Username).Return(&tt.userFound, tt.mockUserDALError).Once()
			}
			if tt.reachJWTService {
				mockAuthProvider.On("GenerateTokens", tt.userFound.ID, tt.userFound.Role).Return(&models.AuthResponse{}, tt.mockAuthError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, _ := json.Marshal(tt.loginRequest)
			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

			handler.LoginController(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUserDAL.AssertExpectations(t)
		})
	}
}

// TODO: Add cookie validation post request
func TestLogout(t *testing.T) {
	mockUserDAL := &mocks.MockUserDAL{}
	mockAuthProvider := &mocks.MockJWTService{}
	handler := handlers.NewAuthController(mockUserDAL, mockAuthProvider)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name                  string
		accessToken           string
		refreshToken          string
		mockValidateAuthError error
		mockRevokeAuthError   error
		expectedAccessToken   string
		expectedRefreshToken  string
		expectedStatus        int
	}{
		{
			name:                  "successful logout",
			accessToken:           "test-access-token",
			refreshToken:          "test-refresh-token",
			mockValidateAuthError: nil,
			mockRevokeAuthError:   nil,
			expectedAccessToken:   "",
			expectedRefreshToken:  "",
			expectedStatus:        http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if tt.accessToken != "" {
				mockAuthProvider.On("ValidateToken", tt.accessToken).Return(&models.UserClaim{}, tt.mockValidateAuthError).Once()
				if tt.mockValidateAuthError == nil {
					mockAuthProvider.On("RevokeToken", tt.accessToken, time.Now()).Once()
				}

			}

			if tt.refreshToken != "" {
				mockAuthProvider.On("ValidateToken", tt.refreshToken).Return(&models.UserClaim{}, tt.mockRevokeAuthError).Once()
				if tt.mockRevokeAuthError == nil {
					mockAuthProvider.On("RevokeToken", tt.refreshToken, time.Now()).Once()
				}

			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodPost, "/", nil)
			c.Header("Authorization", fmt.Sprintf("Bearer %s", tt.accessToken))

			handler.LogoutController(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockUserDAL.AssertExpectations(t)
		})
	}
}
