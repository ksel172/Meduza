package handlers

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
	"github.com/ksel172/Meduza/teamserver/utils"
)

const (
	logLevel         = "INFO"
	logDetailPayload = "PayloadController"
)

type PayloadController struct {
	payloadDAL dal.IPayloadDAL
	agentDAL   dal.IAgentDAL
	buildDir   string
}

func NewPayloadController(payloadDAL dal.IPayloadDAL, agentDAL dal.IAgentDAL) *PayloadController {
	buildDir := "./teamserver/build"
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		if err := os.MkdirAll(buildDir, 0755); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to create build directory: %v", err))
		}
	}

	return &PayloadController{
		payloadDAL: payloadDAL,
		buildDir:   buildDir,
		agentDAL:   agentDAL,
	}
}

/* Payload Manifest endpoints */

// UploadPayloadManifest handles uploading a new payload zip file
func (pc *PayloadController) UploadPayloadManifest(ctx *gin.Context) {
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

	buildDir := pc.buildDir
	if err := os.MkdirAll(buildDir, 0755); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create build directory", err.Error())
		return
	}

	// Strip the .zip extension from the filename for the directory name
	baseFilename := strings.TrimSuffix(filename, filepath.Ext(filename))
	extractDir := fmt.Sprintf("%s/payload-%s", buildDir, baseFilename)

	// Check if this payload name already exists to allow overwriting
	if _, err := os.Stat(extractDir); err == nil {
		// Directory exists, remove it to allow overwriting
		if err := os.RemoveAll(extractDir); err != nil {
			models.ResponseError(ctx, http.StatusInternalServerError,
				"Failed to remove existing payload directory", err.Error())
			return
		}
		logger.Info(logLevel, logDetailPayload,
			fmt.Sprintf("Removed existing payload directory: %s", extractDir))
	}

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

	manifest, err := parseManifestFile(manifestFile, extractDir)
	if err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Failed to parse payload manifest file", err.Error())
		return
	}

	// Save the manifest to database
	if err := pc.payloadDAL.CreatePayloadManifest(ctx.Request.Context(), manifest); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to save payload manifest", err.Error())
		return
	}

	var manifestData map[string]interface{}
	if err := json.Unmarshal([]byte(manifest.Body), &manifestData); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to unmarshal manifest body: %v", err))
		manifestData = map[string]interface{}{}
	} else {
		if sourcePath, ok := manifestData["source_path"].(string); ok {
			// Replace any remaining .zip in the path with empty string
			manifestData["source_path"] = strings.ReplaceAll(sourcePath, ".zip", "")
		}
	}

	// Send successful response
	models.ResponseSuccess(ctx, http.StatusCreated, "Payload folder uploaded and extracted successfully", manifestData)
}

// GetPayloadFormat returns the format of a payload manifest, which includes parameter definitions
func (pc *PayloadController) GetPayloadManifest(ctx *gin.Context) {
	manifestID := ctx.Param(models.ParamManifestID)
	if manifestID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing manifest ID", "Manifest ID is required")
		return
	}

	// Get the manifest from the database
	manifest, err := pc.payloadDAL.GetPayloadManifest(ctx.Request.Context(), manifestID)
	if err != nil {
		models.ResponseError(ctx, http.StatusNotFound, "Failed to get payload manifest", err.Error())
		return
	}

	// Parse the body for a more structured response
	var manifestBody map[string]interface{}
	if err := json.Unmarshal([]byte(manifest.Body), &manifestBody); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to unmarshal manifest body: %v", err))
		// Return the raw manifest if unmarshal fails
		models.ResponseSuccess(ctx, http.StatusOK, "Payload manifest retrieved successfully", manifest)
		return
	}

	// Create a combined response with both metadata and body content
	response := map[string]interface{}{
		"id":               manifest.ID,
		"manifest_version": manifest.ManifestVersion,
		"created_at":       manifest.CreatedAt,
		"updated_at":       manifest.UpdatedAt,
		"content":          manifestBody,
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payload manifest retrieved successfully", response)
}

// DeletePayload deletes a payload manifest and all associated files
func (pc *PayloadController) DeletePayloadManifest(ctx *gin.Context) {
	manifestID := ctx.Param("id")
	if manifestID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing manifest ID", "Manifest ID is required")
		return
	}

	// Get the manifest to find the source path before deleting
	manifest, err := pc.payloadDAL.GetPayloadManifest(ctx.Request.Context(), manifestID)
	if err != nil {
		models.ResponseError(ctx, http.StatusNotFound, "Failed to get payload manifest", err.Error())
		return
	}

	// Extract source path from body
	var manifestData struct {
		SourcePath string `json:"source_path"`
	}
	if err := json.Unmarshal([]byte(manifest.Body), &manifestData); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to unmarshal manifest body: %v", err))
		// Continue with deletion even if we can't extract the source path
	}

	// Delete the manifest from the database
	if err := pc.payloadDAL.DeletePayloadManifest(ctx.Request.Context(), manifestID); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to delete payload manifest", err.Error())
		return
	}

	// Delete the source directory if it exists and is not empty
	sourcePath := manifestData.SourcePath
	if sourcePath != "" && sourcePath != "/" {
		if err := os.RemoveAll(sourcePath); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to delete payload source directory: %v", err))
			// Continue even if the directory deletion fails
		}
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Payload manifest deleted successfully", nil)
}

// GetAvailablePayloads returns a list of all available payload manifests
func (pc *PayloadController) GetAvailablePayloadManifests(ctx *gin.Context) {
	manifests, err := pc.payloadDAL.GetAllPayloadManifests(ctx.Request.Context())
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get payload manifests", err.Error())
		return
	}

	// Prepare a simplified response with only the necessary information
	var response []map[string]interface{}
	for _, manifest := range manifests {
		var bodyData map[string]interface{}
		if err := json.Unmarshal([]byte(manifest.Body), &bodyData); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to unmarshal manifest body: %v", err))
			// Skip this manifest if we can't parse its body
			continue
		}

		// Extract key fields
		name := getStringValue(bodyData, "name")
		version := getStringValue(bodyData, "version")
		author := getStringValue(bodyData, "author")
		description := getStringValue(bodyData, "description")

		// Extract supported architectures
		var supportedArch []string
		if buildConfig, ok := bodyData["payload_build_config"].(map[string]interface{}); ok {
			if archs, ok := buildConfig["supported_arch"].([]interface{}); ok {
				for _, arch := range archs {
					if archStr, ok := arch.(string); ok {
						supportedArch = append(supportedArch, archStr)
					}
				}
			}
		}

		response = append(response, map[string]interface{}{
			"id":               manifest.ID,
			"name":             name,
			"version":          version,
			"author":           author,
			"description":      description,
			"manifest_version": manifest.ManifestVersion,
			"supported_arch":   supportedArch,
		})
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Available payloads retrieved successfully", response)
}

// SubmitBuildJob handles submitting a new build job for a payload
func (pc *PayloadController) SubmitBuildJob(ctx *gin.Context) {
	var request struct {
		ManifestID   string                 `json:"manifest_id" binding:"required"`
		ListenerID   string                 `json:"listener_id" binding:"required"`
		Architecture string                 `json:"architecture" binding:"required"`
		Name         string                 `json:"name" binding:"required"`
		Parameters   map[string]interface{} `json:"parameters"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid request", err.Error())
		return
	}

	// Get the manifest to validate architecture and parameters
	manifest, err := pc.payloadDAL.GetPayloadManifest(ctx.Request.Context(), request.ManifestID)
	if err != nil {
		models.ResponseError(ctx, http.StatusNotFound, "Failed to get payload manifest", err.Error())
		return
	}

	// Parse the manifest body to extract build configuration
	var manifestData struct {
		PayloadBuildConfig struct {
			SupportedArchs []string `json:"supported_arch"`
		} `json:"payload_build_config"`
	}

	if err := json.Unmarshal([]byte(manifest.Body), &manifestData); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to parse manifest configuration", err.Error())
		return
	}

	// Validate architecture
	archValid := false
	for _, arch := range manifestData.PayloadBuildConfig.SupportedArchs {
		if arch == request.Architecture {
			archValid = true
			break
		}
	}
	if !archValid {
		models.ResponseError(ctx, http.StatusBadRequest, "Invalid architecture",
			fmt.Sprintf("Architecture %s is not supported by this payload. Supported: %v",
				request.Architecture, manifestData.PayloadBuildConfig.SupportedArchs))
		return
	}

	// Create payload record
	payloadID := uuid.New().String()
	payload := &models.Payload{
		ID:         payloadID,
		ListenerID: request.ListenerID,
		ConfigID:   uuid.New().String(),
		ManifestID: request.ManifestID,
		Name:       request.Name,
		Arch:       request.Architecture,
		CreatedAt:  time.Now(),
	}

	agentConfig := payload.IntoAgentConfig()

	err = pc.agentDAL.CreateAgentConfig(ctx, agentConfig)
	if err != nil {
		models.ResponseError(
			ctx,
			http.StatusInternalServerError,
			"Failed to create an agent config from the payload config",
			err.Error())
		return
	}

	// Generate RSA key pair for payload
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to generate key pair", err.Error())
		return
	}

	// Convert private key to PEM format
	privateKeyPEM := &bytes.Buffer{}
	pem.Encode(privateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// Convert public key to DER format
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to marshal public key", err.Error())
		return
	}

	// Store keys in payload
	payload.PrivateKey = privateKeyPEM.Bytes()
	payload.PublicKey = publicKeyDER
	payload.Token = uuid.New().String() // Generate unique token for agent authentication

	// Save payload to database
	if err := pc.payloadDAL.CreatePayload(ctx.Request.Context(), payload); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create payload", err.Error())
		return
	}

	// Create build job
	job := &models.PayloadJob{
		ID:           uuid.New().String(),
		PayloadID:    payloadID,
		Status:       "queued",
		Architecture: request.Architecture,
		Parameters:   request.Parameters,
		StartTime:    time.Now(),
	}

	// Add listener ID to parameters if not already present
	if job.Parameters == nil {
		job.Parameters = make(map[string]interface{})
	}

	if _, ok := job.Parameters["listener_id"]; !ok {
		job.Parameters["listener_id"] = request.ListenerID
	}

	// Add token to parameters if needed by the build
	job.Parameters["token"] = payload.Token

	// Save job to database
	if err := pc.payloadDAL.CreateBuildJob(ctx.Request.Context(), job); err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to create build job", err.Error())
		return
	}

	// Launch the build process in a goroutine
	go pc.executeBuild(context.Background(), job, manifest)

	models.ResponseSuccess(ctx, http.StatusAccepted, "Build job submitted successfully", map[string]interface{}{
		"job_id":     job.ID,
		"payload_id": payloadID,
		"status":     job.Status,
	})
}

// GetBuildJob returns the status of a build job
func (pc *PayloadController) GetBuildJob(ctx *gin.Context) {
	jobID := ctx.Param("id")
	if jobID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing job ID", "Job ID is required")
		return
	}

	job, err := pc.payloadDAL.GetBuildJob(ctx.Request.Context(), jobID)
	if err != nil {
		models.ResponseError(ctx, http.StatusNotFound, "Failed to get build job", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Build job retrieved successfully", job)
}

// GetPayloadBuildJobs returns all build jobs for a payload
func (pc *PayloadController) GetPayloadBuildJobs(ctx *gin.Context) {
	payloadID := ctx.Param("id")
	if payloadID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing payload ID", "Payload ID is required")
		return
	}

	jobs, err := pc.payloadDAL.GetPayloadBuildJobs(ctx.Request.Context(), payloadID)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to get build jobs", err.Error())
		return
	}

	models.ResponseSuccess(ctx, http.StatusOK, "Build jobs retrieved successfully", jobs)
}

// DownloadBuildOutput handles downloading the output of a completed build job
func (pc *PayloadController) DownloadBuildOutput(ctx *gin.Context) {
	jobID := ctx.Param("id")
	if jobID == "" {
		models.ResponseError(ctx, http.StatusBadRequest, "Missing job ID", "Job ID is required")
		return
	}

	// Get the build job to check status and get output path
	job, err := pc.payloadDAL.GetBuildJob(ctx.Request.Context(), jobID)
	if err != nil {
		models.ResponseError(ctx, http.StatusNotFound, "Build job not found", err.Error())
		return
	}

	// Check if build is completed
	if job.Status != "completed" {
		models.ResponseError(ctx, http.StatusBadRequest,
			"Build not available for download",
			fmt.Sprintf("Current build status is: %s", job.Status))
		return
	}

	// Verify output path exists
	if job.OutputPath == "" {
		models.ResponseError(ctx, http.StatusInternalServerError,
			"Build output path not specified", "No output file available")
		return
	}

	// Check if file exists
	if _, err := os.Stat(job.OutputPath); os.IsNotExist(err) {
		models.ResponseError(ctx, http.StatusInternalServerError,
			"Build output file not found", fmt.Sprintf("File not found at path: %s", job.OutputPath))
		return
	}

	// Read the file
	data, err := os.ReadFile(job.OutputPath)
	if err != nil {
		models.ResponseError(ctx, http.StatusInternalServerError, "Failed to read output file", err.Error())
		return
	}

	// Generate a friendly filename
	filename := filepath.Base(job.OutputPath)

	// Get payload info to make a more descriptive filename
	payload, err := pc.payloadDAL.GetPayload(ctx.Request.Context(), job.PayloadID)
	if err == nil {
		// Access the payload struct properly
		payloadName := payload.Name
		if payloadName != "" {
			// Sanitize name for filename
			cleanName := strings.Map(func(r rune) rune {
				if r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z' || r >= '0' && r <= '9' || r == '-' || r == '_' {
					return r
				}
				return '_'
			}, payloadName)

			// Add appropriate extension based on architecture
			// Fixed: Using dummy empty map since we don't have access to the manifest here
			extension := getFileExtension(job.Architecture, nil)
			filename = fmt.Sprintf("%s-%s%s", cleanName, job.Architecture, extension)
		}
	}

	// Set headers for file download
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Content-Length", fmt.Sprintf("%d", len(data)))

	ctx.Data(http.StatusOK, "application/octet-stream", data)
}

// Helper function to execute a build job with direct parameter file creation in container
func (pc *PayloadController) executeBuild(ctx context.Context, job *models.PayloadJob, manifest *models.PayloadManifestV1) {
	buildStartTime := time.Now()
	logger.Info(logLevel, logDetailPayload, fmt.Sprintf("Starting build job %s for payload %s", job.ID, job.PayloadID))

	// Update job status to in progress
	job.Status = "in_progress"
	if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		return
	}

	// Create a buffer for in-memory log storage
	var logBuffer bytes.Buffer

	// Create job directory with unique path based on payload ID
	jobDir := filepath.Join(pc.buildDir, "payload-"+job.PayloadID, "jobs", job.ID)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Failed to create job directory: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	// Create log file
	logFilePath := filepath.Join(jobDir, "build.log")
	logFile, err := os.Create(logFilePath)
	if err != nil {
		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Failed to create log file: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}
	defer logFile.Close()

	// Create a multi-writer to write to both in-memory buffer and file
	multiWriter := io.MultiWriter(&logBuffer, logFile)

	// Helper to write log sections
	logSection := func(title string) {
		section := fmt.Sprintf("\n=== %s ===\n", strings.ToUpper(title))
		multiWriter.Write([]byte(section))
	}

	defer func() {
		// Update job with log and status
		job.BuildLog = logBuffer.String()

		// On any panic, mark the job as failed
		if r := recover(); r != nil {
			job.Status = "failed"
			job.ErrorMessage = fmt.Sprintf("Build panic: %v", r)
			job.EndTime = time.Now()

			if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
				logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job after panic: %v", err))
			}

			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Build panic: %v", r))
		}
	}()

	// Start log with build information
	logSection("Build Information")
	fmt.Fprintf(multiWriter, "Build Job ID: %s\n", job.ID)
	fmt.Fprintf(multiWriter, "Payload ID: %s\n", job.PayloadID)
	fmt.Fprintf(multiWriter, "Architecture: %s\n", job.Architecture)
	fmt.Fprintf(multiWriter, "Start Time: %s\n", buildStartTime.Format(time.RFC3339))
	fmt.Fprintf(multiWriter, "Log File: %s\n", logFilePath)

	// Parse the manifest body to extract build configuration
	logSection("Manifest Configuration")
	var manifestBody struct {
		Name               string `json:"name"`
		Version            string `json:"version"`
		Author             string `json:"author"`
		Description        string `json:"description"`
		SourcePath         string `json:"source_path"`
		PayloadBuildConfig struct {
			DockerImage      string            `json:"docker_image"`
			BuildCommand     string            `json:"build_command"`
			BuildArgs        string            `json:"build_args"`
			OutputPath       string            `json:"output_path"`
			OutputFile       string            `json:"output_file"`
			SupportedArchs   []string          `json:"supported_arch"`
			EnvVars          map[string]string `json:"env_vars"`
			ArchMappings     map[string]string `json:"arch_mappings"`
			OutputExtensions map[string]string `json:"output_extensions"`
		} `json:"payload_build_config"`
	}

	if err := json.Unmarshal([]byte(manifest.Body), &manifestBody); err != nil {
		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Failed to parse manifest body: %v", err)
		job.EndTime = time.Now()
		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}
		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	fmt.Fprintf(multiWriter, "Payload Name: %s\n", manifestBody.Name)
	fmt.Fprintf(multiWriter, "Payload Version: %s\n", manifestBody.Version)
	fmt.Fprintf(multiWriter, "Payload Author: %s\n", manifestBody.Author)
	fmt.Fprintf(multiWriter, "Payload Description: %s\n", manifestBody.Description)

	// Create output directory
	logSection("Directory Setup")
	fmt.Fprintf(multiWriter, "Job Directory: %s\n", jobDir)

	outputDir := filepath.Join(jobDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Failed to create output directory: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}
	fmt.Fprintf(multiWriter, "Output Directory: %s\n", outputDir)

	// Process parameters and create config file
	logSection("Build Parameters")
	fmt.Fprintf(multiWriter, "Parameters: %+v\n", job.Parameters)

	paramsJSON, err := json.MarshalIndent(job.Parameters, "", "  ")
	if err != nil {
		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Failed to marshal parameters: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	// Store parameters as a string to pass directly to the container
	paramsJSONStr := string(paramsJSON)
	fmt.Fprintf(multiWriter, "Parameter Contents:\n%s\n", paramsJSONStr)

	// Also save parameters to disk for reference
	paramFile := filepath.Join(jobDir, "parameters.json")
	if err := os.WriteFile(paramFile, paramsJSON, 0644); err != nil {
		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Failed to write parameters file: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	fmt.Fprintf(multiWriter, "Generated parameter file: %s\n", paramFile)

	// Prepare Docker command
	logSection("Build Environment")
	dockerImage := manifestBody.PayloadBuildConfig.DockerImage
	if dockerImage == "" {
		dockerImage = "ubuntu:latest" // Default to ubuntu base image if not specified
	}
	fmt.Fprintf(multiWriter, "Docker Image: %s\n", dockerImage)

	// Set up environment variables based on architecture
	arch := job.Architecture
	envVars := make(map[string]string)

	// Add environment variable for the selected architecture
	envVars["TARGET_ARCH"] = arch

	// Add any custom environment variables from the manifest
	for key, value := range manifestBody.PayloadBuildConfig.EnvVars {
		envVars[key] = value
	}

	// Apply architecture mappings from manifest if defined
	if archValue, ok := manifestBody.PayloadBuildConfig.ArchMappings[arch]; ok {
		envVars["TARGET_ARCH"] = archValue
	}

	// Log all environment variables
	fmt.Fprintf(multiWriter, "Build Environment Variables:\n")
	for key, value := range envVars {
		fmt.Fprintf(multiWriter, "  %s=%s\n", key, value)
	}

	// Define output filename with proper extension
	outputFile := manifestBody.PayloadBuildConfig.OutputFile
	if outputFile == "" {
		// Generate a default filename if none specified
		extension := getFileExtension(arch, manifestBody.PayloadBuildConfig.OutputExtensions)
		outputFile = fmt.Sprintf("payload-%s-%s%s", job.PayloadID[:8], arch, extension)
	}
	fmt.Fprintf(multiWriter, "Output File: %s\n", outputFile)

	// Get paths
	sourcePath, _ := filepath.Abs(manifestBody.SourcePath)
	outputDirPath, _ := filepath.Abs(outputDir)

	fmt.Fprintf(multiWriter, "Source Path: %s\n", sourcePath)
	fmt.Fprintf(multiWriter, "Source Path: %s%s\n", sourcePath, "/simple-agent-test")
	fmt.Fprintf(multiWriter, "Output Path: %s\n", outputDirPath)

	// Prepare base Docker command with volume mounts
	buildArgs := []string{
		"run", "--rm",
		"-v", fmt.Sprintf("%s:/src", sourcePath),
		"-v", fmt.Sprintf("%s:/output", outputDirPath),
		"-w", "/src",
	}

	// Add environment variables to Docker command
	for key, value := range envVars {
		buildArgs = append(buildArgs, "-e", fmt.Sprintf("%s=%s", key, value))
	}

	// Add Docker image
	buildArgs = append(buildArgs, dockerImage)

	// Escape special characters in JSON for shell
	escapedParams := strings.ReplaceAll(paramsJSONStr, `"`, `\"`)
	escapedParams = strings.ReplaceAll(escapedParams, "$", "\\$")
	escapedParams = strings.ReplaceAll(escapedParams, "`", "\\`")

	// Get build command from manifest or use a simple echo command if none provided
	buildCmd := manifestBody.PayloadBuildConfig.BuildCommand
	if buildCmd == "" {
		// Default to a simple shell script if no command is provided
		buildCmd = "echo 'No build command specified in manifest. Create a build script or specify build_command.' && exit 1"
	} else {
		// Replace template variables in build command
		buildCmd = strings.ReplaceAll(buildCmd, "{OUTPUT_FILE}", outputFile)
		buildCmd = strings.ReplaceAll(buildCmd, "{BUILD_ARGS}", manifestBody.PayloadBuildConfig.BuildArgs)
		buildCmd = strings.ReplaceAll(buildCmd, "{OUTPUT_DIR}", "/output")
		buildCmd = strings.ReplaceAll(buildCmd, "{ARCH}", arch)
	}

	// Create full shell command that first creates parameters file
	fullCmd := fmt.Sprintf("echo \"%s\" > ./params.json && "+
		"echo 'Using parameters:' && cat ./params.json && "+
		"echo 'Build environment:' && env | sort && "+
		"echo 'Source directory contents:' && ls -la && "+
		"echo 'Starting build with command: %s' && %s",
		escapedParams,
		buildCmd,
		buildCmd)

	buildArgs = append(buildArgs, "sh", "-c", fullCmd)

	// Log the command
	logSection("Build Command")
	fmt.Fprintf(multiWriter, "Running Docker build: docker %s\n", strings.Join(buildArgs, " "))

	// Time the actual build execution
	buildExecStart := time.Now()
	logSection("Build Process Output")

	// Execute docker command
	cmd := exec.CommandContext(ctx, "docker", buildArgs...)
	cmd.Stdout = multiWriter
	cmd.Stderr = multiWriter
	logger.Info(logLevel, logDetailPayload, fmt.Sprintf("Running Docker command: docker %s", strings.Join(buildArgs, " ")))

	if err := cmd.Run(); err != nil {
		buildDuration := time.Since(buildExecStart)
		logSection("Build Failure")
		fmt.Fprintf(multiWriter, "Build failed after %s: %v\n", buildDuration, err)

		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Docker build failed: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	buildDuration := time.Since(buildExecStart)

	// Verify the output file exists
	outputFilePath := filepath.Join(outputDir, outputFile)
	logSection("Build Verification")
	fmt.Fprintf(multiWriter, "Verifying output file: %s\n", outputFilePath)

	fileInfo, err := os.Stat(outputFilePath)
	if os.IsNotExist(err) {
		logSection("Verification Failed")
		fmt.Fprintf(multiWriter, "Build completed but output file not found: %s\n", outputFilePath)

		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Build completed but output file not found: %s", outputFilePath)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	} else if err != nil {
		logSection("Verification Failed")
		fmt.Fprintf(multiWriter, "Error checking output file: %v\n", err)

		job.Status = "failed"
		job.ErrorMessage = fmt.Sprintf("Error verifying output file: %v", err)
		job.EndTime = time.Now()

		if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		}

		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	// File exists and is accessible - get size and permissions
	fmt.Fprintf(multiWriter, "Output file created successfully\n")
	fmt.Fprintf(multiWriter, "File size: %d bytes\n", fileInfo.Size())
	fmt.Fprintf(multiWriter, "File permissions: %s\n", fileInfo.Mode().String())
	fmt.Fprintf(multiWriter, "File modified time: %s\n", fileInfo.ModTime().Format(time.RFC3339))

	// Mark job as completed
	totalDuration := time.Since(buildStartTime)
	logSection("Build Summary")
	fmt.Fprintf(multiWriter, "Build completed successfully\n")
	fmt.Fprintf(multiWriter, "Compilation time: %s\n", buildDuration)
	fmt.Fprintf(multiWriter, "Total build time: %s\n", totalDuration)
	fmt.Fprintf(multiWriter, "Output file: %s\n", outputFilePath)
	fmt.Fprintf(multiWriter, "Output size: %d bytes\n", fileInfo.Size())
	fmt.Fprintf(multiWriter, "Full build log saved to: %s\n", logFilePath)

	job.Status = "completed"
	job.OutputPath = outputFilePath
	job.EndTime = time.Now()

	if err := pc.payloadDAL.UpdateBuildJob(ctx, job); err != nil {
		logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to update job status: %v", err))
		return
	}

	logger.Info(logLevel, logDetailPayload, fmt.Sprintf("Build job %s completed successfully in %s", job.ID, totalDuration))
}

// Helper functions

// findManifestFile searches for a manifest.json file in the root directory and subdirectories
func findManifestFile(rootDir string) (string, error) {
	// Common manifest file names to check
	manifestNames := []string{
		"manifest.json",
		"payload.json",
		"config.json",
	}

	// Recursive function to find manifest file with depth limit
	var searchInDir func(dir string, depth int) (string, error)
	searchInDir = func(dir string, depth int) (string, error) {
		// Limit recursion depth to prevent excessive searching
		if depth > 3 {
			return "", nil
		}

		files, err := os.ReadDir(dir)
		if err != nil {
			return "", fmt.Errorf("failed to read directory %s: %w", dir, err)
		}

		// First check for manifest files in current directory
		for _, file := range files {
			if !file.IsDir() {
				fileName := strings.ToLower(file.Name())
				for _, manifestName := range manifestNames {
					if fileName == manifestName {
						return filepath.Join(dir, file.Name()), nil
					}
				}
			}
		}

		// Then check subdirectories
		for _, file := range files {
			if file.IsDir() {
				subdir := filepath.Join(dir, file.Name())
				if found, err := searchInDir(subdir, depth+1); err == nil && found != "" {
					return found, nil
				}
			}
		}

		// For depth 0 (root dir), also check for any JSON file as a fallback
		if depth == 0 {
			for _, file := range files {
				if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".json") {
					return filepath.Join(dir, file.Name()), nil
				}
			}
		}

		return "", nil
	}

	// Start search from root directory
	manifestPath, err := searchInDir(rootDir, 0)
	if err != nil {
		return "", err
	}

	if manifestPath == "" {
		return "", fmt.Errorf("no manifest file found in payload root directory")
	}

	return manifestPath, nil
}

// parseManifestFile reads and parses a manifest file into a PayloadManifestV1 structure
func parseManifestFile(manifestPath string, sourcePath string) (*models.PayloadManifestV1, error) {
	// Read the manifest file
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest file: %w", err)
	}

	// First make sure it's valid JSON
	var rawManifest map[string]interface{}
	if err := json.Unmarshal(data, &rawManifest); err != nil {
		return nil, fmt.Errorf("invalid JSON in manifest file: %w", err)
	}

	// Determine the actual source path
	manifestDir := filepath.Dir(manifestPath)
	actualSourcePath := sourcePath

	// If manifest is one level down, use that directory as the source path
	if filepath.Dir(manifestDir) == sourcePath {
		actualSourcePath = manifestDir
	}

	// Add source path if not present or override it with the correct one
	rawManifest["source_path"] = actualSourcePath

	// Re-marshal to get the complete JSON with source_path
	updatedData, err := json.Marshal(rawManifest)
	if err != nil {
		return nil, fmt.Errorf("failed to re-marshal manifest data: %w", err)
	}

	// Create and populate the manifest
	manifest := &models.PayloadManifestV1{
		ID:              uuid.New().String(),
		ManifestVersion: getVersionFloat(rawManifest["manifest_version"]),
		Body:            string(updatedData),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Add validation
	if manifest.ManifestVersion == 0 {
		return nil, fmt.Errorf("manifest version is missing or invalid")
	}

	// Validate required fields in the manifest body
	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return manifest, nil
}

// getVersionString converts a manifest_version value to a string representation
func getVersionFloat(version interface{}) float32 {
	switch v := version.(type) {
	case string:
		// Try to parse string as float
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f)
		}
		// Handle version strings like "v1.0"
		v = strings.TrimPrefix(v, "v")
		if f, err := strconv.ParseFloat(v, 32); err == nil {
			return float32(f)
		}
		return 1.0 // Default if parsing fails
	case float64:
		return float32(v)
	case float32:
		return v
	case int:
		return float32(v)
	default:
		return 1.0 // Default if not specified
	}
}

// getStringValue safely extracts a string from a map
func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		switch v := val.(type) {
		case string:
			return v
		case float64:
			return fmt.Sprintf("%.1f", v)
		case int:
			return fmt.Sprintf("%d", v)
		default:
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

// getFileExtension gets the appropriate file extension for the given architecture
// Consolidated function to handle file extensions with manifest configuration
func getFileExtension(arch string, extensionMap map[string]string) string {
	// First check if there's a specific mapping in the manifest
	if extensionMap != nil {
		if ext, ok := extensionMap[arch]; ok {
			return ext
		}
	}

	// Default extensions based on OS part of architecture
	if strings.HasPrefix(arch, "win-") {
		return ".exe"
	} else if strings.HasPrefix(arch, "darwin-") {
		return ".macho"
	} else if strings.HasPrefix(arch, "linux-") || strings.HasPrefix(arch, "freebsd-") {
		return ".bin"
	}

	return ""
}
