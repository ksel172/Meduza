package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/conf"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

// CertificateHandler handles certificate-related requests
type CertificateHandler struct {
	certDAL dal.CertificateDAL
}

// NewCertificateHandler creates a new certificate handler
func NewCertificateHandler(certDAL dal.CertificateDAL) *CertificateHandler {
	// Create certificate upload directory if it doesn't exist
	uploadPath := conf.GetCertUploadPath()
	if err := os.MkdirAll(uploadPath, 0755); err != nil {
		logger.Error("Failed to create certificate upload directory", err)
	}

	return &CertificateHandler{
		certDAL: certDAL,
	}
}

// UploadCertificate handles the certificate upload
func (h *CertificateHandler) UploadCertificate(ctx *gin.Context) {
	certType := ctx.Param(models.ParamCertificateType)

	// Validate certificate type
	if certType != "cert" && certType != "key" {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid certificate type", "Must be 'cert' or 'key'")
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "No file provided", err.Error())
		return
	}

	// Validate file extension based on certificate type
	var allowedExts string
	if certType == "cert" {
		allowedExts = ".crt,.pem,.cer"
	} else {
		allowedExts = ".key,.pem"
	}

	ext := filepath.Ext(file.Filename)
	if !strings.Contains(allowedExts, ext) {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid file extension",
			fmt.Sprintf("Allowed extensions for %s are: %s", certType, allowedExts))
		return
	}

	// Create a unique filename to prevent overwriting
	uploadPath := conf.GetCertUploadPath()
	filename := fmt.Sprintf("%s-%s", certType, filepath.Base(file.Filename))
	
	// Validate filename to ensure it does not contain any path separators or parent directory references
	if strings.Contains(filename, "/") || strings.Contains(filename, "\\") || strings.Contains(filename, "..") {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid file name", "File name contains invalid characters")
		return
	}
	
	filePath := filepath.Join(uploadPath, filename)
	
	// Ensure the resolved path is within the safe directory
	absPath, err := filepath.Abs(filePath)
	if err != nil || !strings.HasPrefix(absPath, uploadPath) {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid file path", "File path is outside the allowed directory")
		return
	}
	
	// Save the file
	if err := ctx.SaveUploadedFile(file, absPath); err != nil {
		logger.Error("Failed to save certificate file", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to save certificate file", err.Error())
		return
	}

	// Save certificate info to database
	if err := h.certDAL.SaveCertificate(ctx, certType, filePath, file.Filename); err != nil {
		logger.Error("Failed to store certificate information", err)
		// Clean up the file if we failed to save to database
		os.Remove(filePath)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to store certificate information", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Certificate uploaded successfully", gin.H{
		"type":     certType,
		"filename": file.Filename,
		"path":     filePath,
	})
}

// GetCertificates retrieves all certificates
func (h *CertificateHandler) GetCertificates(ctx *gin.Context) {
	certs, err := h.certDAL.GetAllCertificates(ctx)
	if err != nil {
		logger.Error("Failed to retrieve certificates", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to retrieve certificates", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Certificates retrieved successfully", certs)
}

// DeleteCertificate deletes a certificate
func (h *CertificateHandler) DeleteCertificate(ctx *gin.Context) {
	certID := ctx.Param(models.ParamCertificateID)

	err := h.certDAL.DeleteCertificate(ctx, certID)
	if err != nil {
		logger.Error("Failed to delete certificate", err)
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete certificate", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Certificate deleted successfully", nil)
}
