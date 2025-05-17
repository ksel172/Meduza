package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/utils"
)

type PayloadController struct {
	payloadDAL dal.IPayloadDAL
}

func NewPayloadController(payloadDAL dal.IPayloadDAL) *PayloadController {
	return &PayloadController{
		payloadDAL: payloadDAL,
	}
}

func (pc *PayloadController) UploadPayload(ctx *gin.Context) {
	sourceType := ctx.PostForm("sourceType")
	var data []byte
	var filename string
	var err error

	if sourceType == "file" {
		// Handle file upload
		file, fileHeader, err := ctx.Request.FormFile("payload")
		if err != nil {
			models.ResponseError(ctx, http.StatusBadRequest, "Failed to get file", err.Error())
			return
		}
		defer file.Close()

		filename = fileHeader.Filename

		// Verify file is a ZIP archive
		if !strings.HasSuffix(strings.ToLower(filename), ".zip") {
			models.ResponseError(ctx, http.StatusBadRequest, "File must be a ZIP archive for folder uploads", "Invalid file extension")
			return
		}

		// Read file bytes
		buffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(buffer, file); err != nil {
			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to read file", err.Error())
			return
		}
		data = buffer.Bytes()
	} else if sourceType == "url" {
		// Handle URL source
		url := ctx.PostForm("url")
		if url == "" {
			models.ResponseError(ctx, http.StatusBadRequest, "URL is required for URL source type", "Missing URL parameter")
			return
		}

		// Download file from URL
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			models.ResponseError(ctx, http.StatusBadRequest, "Failed to download from URL", err.Error())
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			models.ResponseError(ctx, http.StatusBadRequest, "Bad response from URL", resp.Status)
			return
		}

		// Read response body
		buffer := bytes.NewBuffer(nil)
		if _, err := io.Copy(buffer, resp.Body); err != nil {
			models.ResponseError(ctx, http.StatusInternalServerError, "Failed to read response body", err.Error())
			return
		}
		data = buffer.Bytes()

		// Extract filename from URL
		urlParts := strings.Split(url, "/")
		filename = urlParts[len(urlParts)-1]
		if filename == "" || !strings.HasSuffix(strings.ToLower(filename), ".zip") {
			filename = "payload-" + time.Now().Format("20060102-150405") + ".zip"
		}
	} else {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid source type", "Must be 'file' or 'url'")
		return
	}

	// Generate a unique ID for the payload
	payloadID := uuid.New().String()

	// Create directories
	// This will be mounted in the Docker volume
	buildDir := "./teamserver/build"
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create build directory", err.Error())
		return
	}

	extractDir := fmt.Sprintf("%s/payload-%s", buildDir, payloadID)
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create extraction directory", err.Error())
		return
	}

	// Create a temporary file to store the ZIP
	tempFile := filepath.Join(os.TempDir(), filename)
	if err := os.WriteFile(tempFile, data, 0644); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to write temporary file", err.Error())
		return
	}
	defer os.Remove(tempFile) // Clean up temp file when done

	// Extract the ZIP file
	if err := utils.Unzip(tempFile, extractDir); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to extract ZIP file", err.Error())
		return
	}

	// Look for manifest file in the root directory of the extracted payload
	manifestFile, err := findManifestFile(extractDir)
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to find payload manifest file", err.Error())
		return
	}

	// Parse the manifest file
	manifest, err := parseManifestFile(manifestFile)
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to parse payload manifest file", err.Error())
		return
	}

	// Set payload ID and create the payload record
	manifest.ID = payloadID
	manifest.SourcePath = extractDir

	// Save the manifest to database (this now stores the entire object as JSON)
	if err := pc.payloadDAL.CreatePayloadManifest(ctx.Request.Context(), manifest); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to save payload manifest", err.Error())
		return
	}

	// Send successful response
	models.ResponseSuccess(ctx, http.StatusCreated, "Payload folder uploaded and extracted successfully", map[string]interface{}{
		"payloadID": payloadID,
		"filename":  filename,
		"path":      extractDir,
		"manifest":  manifest,
	})
}

// findManifestFile searches for a manifest.json file in the root directory
func findManifestFile(rootDir string) (string, error) {
	// Common manifest file names to check
	manifestNames := []string{
		"manifest.json",
		"payload.json",
		"config.json",
	}

	// Read the directory contents
	files, err := os.ReadDir(rootDir)
	if err != nil {
		return "", fmt.Errorf("failed to read directory: %w", err)
	}

	// First, check for the exact manifest names
	for _, file := range files {
		if !file.IsDir() {
			fileName := strings.ToLower(file.Name())
			for _, manifestName := range manifestNames {
				if fileName == manifestName {
					return filepath.Join(rootDir, file.Name()), nil
				}
			}
		}
	}

	// If no exact matches, look for any JSON file in the root directory
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
			return filepath.Join(rootDir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no manifest file found in payload root directory")
}

// parseManifestFile reads and parses a manifest file into a RawPayload structure
func parseManifestFile(manifestPath string) (*models.PayloadManifestV1, error) {
	// Read the manifest file
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	// Parse the JSON into a RawPayload structure
	var rawPayload models.PayloadManifestV1
	if err := json.Unmarshal(data, &rawPayload); err != nil {
		return nil, fmt.Errorf("failed to parse manifest JSON: %w", err)
	}

	// Validate required fields
	if rawPayload.Name == "" {
		return nil, fmt.Errorf("manifest missing required field: name")
	}

	return &rawPayload, nil
}

func (pc *PayloadController) GetPayloadFormat(ctx *gin.Context) {

}

func (pc *PayloadController) DeletePayload(ctx *gin.Context) {

}

func (pc *PayloadController) GetAvailablePayloads(ctx *gin.Context) {

}
