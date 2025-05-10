package http_listener

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	controller_mocks "github.com/ksel172/Meduza/teamserver/internal/mocks/controller"
	"github.com/ksel172/Meduza/teamserver/internal/services/listener/checkin"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/assert"
)

func TestAgentRegisterRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testAuthToken := "test-auth-token"
	sessionToken, aesKey, err := authenticateAgent(testAuthToken)
	if err != nil {
		t.Fatalf("failed agent authentication: %v", err)
	}
	encodedSessionToken := base64.StdEncoding.EncodeToString([]byte(sessionToken))

	// For this test, we only simulate valid c2requests.
	// In the checkin module, we will validate parsing of multiple kinds of c2requests
	c2request := models.C2Request{
		Reason:  models.Register,
		AgentID: "test-agent-id",
	}

	tests := []struct {
		name           string
		setupFunc      func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController)
		c2request      models.C2Request
		sessionToken   string
		expectedStatus int
	}{
		{
			name: "register: success",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleRegisterRequest", c2request).Return(nil).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusCreated,
		},
		{
			name: "register: invalid data response",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleRegisterRequest", c2request).Return(checkin.ErrInvalidData).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "register: agent already exists conflict",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleRegisterRequest", c2request).Return(checkin.ErrConflict).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "register: invalid data response",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleRegisterRequest", c2request).Return(checkin.ErrInternalServer).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "register: checkin server error",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleRegisterRequest", c2request).Return(checkin.ErrInternalServer).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "register: invalid session token",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   "invalid-session-token",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener, mockCheckinController := tt.setupFunc(tt.c2request)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, err := encryptAgentRequest(tt.c2request, aesKey)
			if err != nil {
				t.Fatalf("failed to encrypt c2 request")
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Request.Header.Add("Session-Token", tt.sessionToken)

			listener.HandleCheckIn(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckinController.AssertExpectations(t)
		})
	}
}

func TestAgentTaskRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testAuthToken := "test-auth-token"
	sessionToken, aesKey, err := authenticateAgent(testAuthToken)
	if err != nil {
		t.Fatalf("failed agent authentication: %v", err)
	}
	encodedSessionToken := base64.StdEncoding.EncodeToString([]byte(sessionToken))

	// For this test, we only simulate valid c2requests.
	// In the checkin module, we will validate parsing of multiple kinds of c2requests
	c2request := models.C2Request{
		Reason:  models.Task,
		AgentID: "test-agent-id",
	}

	tests := []struct {
		name           string
		setupFunc      func(models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController)
		c2request      models.C2Request
		sessionToken   string
		expectedStatus int
	}{
		{
			name: "task: success",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleTaskRequest", c2request, sessionToken).Return([]byte{}, nil).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusOK,
		},
		{
			name: "task: database error",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleTaskRequest", c2request, sessionToken).Return([]byte{}, checkin.ErrDatabase).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusNotFound,
		},
		{
			name: "task: agent already exists conflict",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleTaskRequest", c2request, sessionToken).Return([]byte{}, checkin.ErrInternalServer).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "task: invalid data response",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleTaskRequest", c2request, sessionToken).Return([]byte{}, checkin.ErrUnauthorized).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener, mockCheckInController := tt.setupFunc(tt.c2request)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, err := encryptAgentRequest(tt.c2request, aesKey)
			if err != nil {
				t.Fatalf("failed to encrypt c2 request")
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Request.Header.Add("Session-Token", tt.sessionToken)

			listener.HandleCheckIn(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInController.AssertExpectations(t)
		})
	}
}

func TestAgentResponseRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	testAuthToken := "test-auth-token"
	sessionToken, aesKey, err := authenticateAgent(testAuthToken)
	if err != nil {
		t.Fatalf("failed agent authentication: %v", err)
	}
	encodedSessionToken := base64.StdEncoding.EncodeToString([]byte(sessionToken))

	// For this test, we only simulate valid c2requests.
	// In the checkin module, we will validate parsing of multiple kinds of c2requests
	c2request := models.C2Request{
		Reason:  models.Response,
		AgentID: "test-agent-id",
	}

	tests := []struct {
		name           string
		setupFunc      func(models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController)
		c2request      models.C2Request
		sessionToken   string
		expectedStatus int
	}{
		{
			name: "response: success",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleResponseRequest", c2request).Return(nil).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusOK,
		},
		{
			name: "response: database error",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleResponseRequest", c2request).Return(checkin.ErrInvalidData).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "response: agent already exists conflict",
			setupFunc: func(c2request models.C2Request) (*HTTPListener, *controller_mocks.MockCheckInController) {
				mockCheckInController := &controller_mocks.MockCheckInController{}
				listener := buildDefaultTestHTTPListenerWithMock(t, mockCheckInController)
				mockCheckInController.On("HandleResponseRequest", c2request).Return(checkin.ErrInternalServer).Once()
				return listener, mockCheckInController
			},
			c2request:      c2request,
			sessionToken:   encodedSessionToken,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			listener, mockCheckInController := tt.setupFunc(tt.c2request)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body, err := encryptAgentRequest(tt.c2request, aesKey)
			if err != nil {
				t.Fatalf("failed to encrypt c2 request")
			}

			c.Request = httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
			c.Request.Header.Add("Session-Token", tt.sessionToken)

			listener.HandleCheckIn(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCheckInController.AssertExpectations(t)
		})
	}
}
