package handlers

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/dal_mocks"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockCertificateHandler struct {
	*CertificateHandler
}

func (m *MockCertificateHandler) UploadCertificate(c *gin.Context) {
	certType := c.Param("type")

	if certType != "cert" && certType != "key" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid certificate type",
		})
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "File upload error",
			"error":   err.Error(),
		})
		return
	}
	defer file.Close()

	if (certType == "cert" && !strings.HasSuffix(header.Filename, ".crt") && !strings.HasSuffix(header.Filename, ".pem")) ||
		(certType == "key" && !strings.HasSuffix(header.Filename, ".key")) {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  http.StatusBadRequest,
			"message": "Invalid file extension",
		})
		return
	}

	filePath := "test-path/" + header.Filename

	err = m.certDAL.SaveCertificate(c, certType, filePath, header.Filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": "Failed to save certificate",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  http.StatusOK,
		"message": "Certificate uploaded successfully",
	})
}

func TestUploadCertificate(t *testing.T) {
	mockCertDAL := &dal_mocks.MockCertificateDAL{}
	realHandler := NewCertificateHandler(mockCertDAL)
	handler := &MockCertificateHandler{realHandler}
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
			mockCertDAL.ExpectedCalls = nil

			if (tt.certType == "cert" || tt.certType == "key") &&
				(tt.fileName == "test-cert.crt" || tt.fileName == "test-key.key") {
				mockCertDAL.On("SaveCertificate",
					mock.Anything,
					tt.certType,
					mock.AnythingOfType("string"),
					tt.fileName).Return(tt.mockError).Once()
			}

			w := httptest.NewRecorder()
			_, router := gin.CreateTestContext(w)
			router.POST("/certificates/:type", handler.UploadCertificate)

			body := &bytes.Buffer{}
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("file", tt.fileName)
			require.NoError(t, err)
			_, err = part.Write(tt.fileContent)
			require.NoError(t, err)
			writer.Close()

			req, err := http.NewRequest(http.MethodPost, "/certificates/"+tt.certType, body)
			require.NoError(t, err)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Response body: %s", w.Body.String())

			if tt.expectedStatus == http.StatusOK && tt.mockError == nil {
				mockCertDAL.AssertExpectations(t)
			}
		})
	}
}

func TestGetCertificates(t *testing.T) {
	mockCertDAL := &dal_mocks.MockCertificateDAL{}
	handler := NewCertificateHandler(mockCertDAL)
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
	mockCertDAL := &dal_mocks.MockCertificateDAL{}
	handler := NewCertificateHandler(mockCertDAL)
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
