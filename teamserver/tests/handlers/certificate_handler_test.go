package handler_tests

import (
	"bytes"
	"errors"
	"mime/multipart"
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

func TestUploadCertificate(t *testing.T) {
	mockCertDAL := &mocks.MockCertificateDAL{}
	handler := handlers.NewCertificateHandler(mockCertDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		certType       string
		fileName       string
		fileContent    []byte
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful upload cert",
			certType:       "cert",
			fileName:       "test-cert.crt",
			fileContent:    []byte("test-cert-content"),
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "successful upload key",
			certType:       "key",
			fileName:       "test-key.key",
			fileContent:    []byte("test-key-content"),
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid cert type",
			certType:       "invalid",
			fileName:       "test-cert.crt",
			fileContent:    []byte("test-cert-content"),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid file extension",
			certType:       "cert",
			fileName:       "test-cert.txt",
			fileContent:    []byte("test-cert-content"),
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "dal error",
			certType:       "cert",
			fileName:       "test-cert.crt",
			fileContent:    []byte("test-cert-content"),
			mockError:      errors.New("failed dal op"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.expectedStatus == http.StatusOK {
				mockCertDAL.On("SaveCertificate", mock.Anything, tt.certType, mock.Anything, tt.fileName).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)
			part, _ := writer.CreateFormFile("file", tt.fileName)
			part.Write(tt.fileContent)
			writer.Close()

			c.Request = httptest.NewRequest(http.MethodPost, "/certificates/"+tt.certType, body)
			c.Request.Header.Set("Content-Type", writer.FormDataContentType())

			handler.UploadCertificate(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCertDAL.AssertExpectations(t)
		})
	}
}

func TestGetCertificates(t *testing.T) {
	mockCertDAL := &mocks.MockCertificateDAL{}
	handler := handlers.NewCertificateHandler(mockCertDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockCerts      []models.Certificate
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful get certificates",
			mockCerts:      []models.Certificate{{ID: "test-cert-id", Type: "cert", Path: "test-path", Filename: "test-cert.crt"}},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "dal error",
			mockCerts:      nil,
			mockError:      errors.New("failed dal op"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCertDAL.On("GetAllCertificates", mock.Anything).Return(tt.mockCerts, tt.mockError).Once()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/certificates", nil)

			handler.GetCertificates(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCertDAL.AssertExpectations(t)
		})
	}
}

func TestDeleteCertificate(t *testing.T) {
	mockCertDAL := &mocks.MockCertificateDAL{}
	handler := handlers.NewCertificateHandler(mockCertDAL)
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		certID         string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful delete certificate",
			certID:         "test-cert-id",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "dal error",
			certID:         "test-cert-id",
			mockError:      errors.New("failed dal op"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCertDAL.On("DeleteCertificate", mock.Anything, tt.certID).Return(tt.mockError).Once()

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Params = gin.Params{{Key: models.ParamCertificateID, Value: tt.certID}}

			handler.DeleteCertificate(c)

			assert.Equal(t, tt.expectedStatus, w.Code)
			mockCertDAL.AssertExpectations(t)
		})
	}
}
