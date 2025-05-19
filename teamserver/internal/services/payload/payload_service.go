package payload

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/internal/storage/dal"
	"github.com/ksel172/Meduza/teamserver/models"
	"github.com/ksel172/Meduza/teamserver/pkg/logger"
)

const (
	logLevel         = "INFO"
	logDetailPayload = "PayloadBuildService"
)

// BuildStatus represents the status of a build job
type BuildStatus string

const (
	BuildStatusQueued     BuildStatus = "queued"
	BuildStatusInProgress BuildStatus = "in_progress"
	BuildStatusCompleted  BuildStatus = "completed"
	BuildStatusFailed     BuildStatus = "failed"
)

// BuildJob represents a payload build job
type BuildJob struct {
	ID           string            `json:"id"`
	PayloadID    string            `json:"payload_id"`
	Status       BuildStatus       `json:"status"`
	Architecture string            `json:"architecture"`
	Parameters   map[string]string `json:"parameters"`
	StartTime    time.Time         `json:"start_time"`
	EndTime      time.Time         `json:"end_time,omitempty"`
	OutputPath   string            `json:"output_path,omitempty"`
	ErrorMessage string            `json:"error_message,omitempty"`
	BuildLog     string            `json:"build_log,omitempty"`
}

// ManifestData holds the parsed JSON from manifest body
type ManifestData struct {
	Name               string `json:"name"`
	Version            string `json:"version"`
	Author             string `json:"author"`
	Description        string `json:"description"`
	SourcePath         string `json:"source_path"`
	PayloadBuildConfig struct {
		DockerImage    string   `json:"docker_image"`
		BuildArgs      string   `json:"build_args"`
		OutputPath     string   `json:"output_path"`
		OutputFile     string   `json:"output_file"`
		SupportedArchs []string `json:"supported_arch"`
	} `json:"payload_build_config"`
	Parameters struct {
		ListenerID       string `json:"listener_id"`
		CustomParameters []struct {
			DataType      string `json:"data_type"`
			ParameterName string `json:"parameter_name"`
			Description   string `json:"description"`
			DefaultValue  string `json:"default_value"`
			Required      bool   `json:"required"`
		} `json:"custom_parameters"`
	} `json:"parameters"`
}

// PayloadBuildService handles payload building operations
type PayloadBuildService struct {
	payloadDAL dal.IPayloadDAL
	buildJobs  map[string]*BuildJob
	buildDir   string
}

// NewPayloadBuildService creates a new payload build service
func NewPayloadBuildService(payloadDAL dal.IPayloadDAL, buildDir string) *PayloadBuildService {
	// Create build directory if it doesn't exist
	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		err := os.MkdirAll(buildDir, 0755)
		if err != nil {
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Failed to create build directory: %v", err))
		}
	}

	return &PayloadBuildService{
		payloadDAL: payloadDAL,
		buildJobs:  make(map[string]*BuildJob),
		buildDir:   buildDir,
	}
}

// SubmitBuildJob creates a new build job for a payload
func (s *PayloadBuildService) SubmitBuildJob(ctx context.Context, payloadID string, arch string, params map[string]string) (string, error) {
	// Fetch payload manifest
	manifest, err := s.payloadDAL.GetPayloadManifest(ctx, payloadID)
	if err != nil {
		return "", fmt.Errorf("failed to get payload manifest: %w", err)
	}

	// Parse manifest body
	var manifestData ManifestData
	if err := json.Unmarshal([]byte(manifest.Body), &manifestData); err != nil {
		return "", fmt.Errorf("failed to parse manifest body: %w", err)
	}

	// Validate architecture
	archSupported := false
	for _, supportedArch := range manifestData.PayloadBuildConfig.SupportedArchs {
		if supportedArch == arch {
			archSupported = true
			break
		}
	}

	if !archSupported {
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	// Create build job
	jobID := uuid.New().String()
	job := &BuildJob{
		ID:           jobID,
		PayloadID:    payloadID,
		Status:       BuildStatusQueued,
		Architecture: arch,
		Parameters:   params,
		StartTime:    time.Now(),
	}

	// Store job
	s.buildJobs[jobID] = job

	// Start build in goroutine
	go s.executeBuild(context.Background(), job, manifest)

	return jobID, nil
}

// GetBuildJob retrieves the status of a build job
func (s *PayloadBuildService) GetBuildJob(jobID string) (*BuildJob, error) {
	job, exists := s.buildJobs[jobID]
	if !exists {
		return nil, fmt.Errorf("build job not found: %s", jobID)
	}
	return job, nil
}

// GetPayloadBuildJobs retrieves all build jobs for a payload
func (s *PayloadBuildService) GetPayloadBuildJobs(payloadID string) ([]*BuildJob, error) {
	var jobs []*BuildJob
	for _, job := range s.buildJobs {
		if job.PayloadID == payloadID {
			jobs = append(jobs, job)
		}
	}
	return jobs, nil
}

// executeBuild performs the actual build process
func (s *PayloadBuildService) executeBuild(ctx context.Context, job *BuildJob, manifest *models.PayloadManifestV1) {
	logger.Info(logLevel, logDetailPayload, fmt.Sprintf("Starting build job %s for payload %s", job.ID, job.PayloadID))

	// Update job status
	job.Status = BuildStatusInProgress

	var logBuffer bytes.Buffer
	defer func() {
		job.BuildLog = logBuffer.String()

		// On any panic, mark the job as failed
		if r := recover(); r != nil {
			job.Status = BuildStatusFailed
			job.ErrorMessage = fmt.Sprintf("Build panic: %v", r)
			job.EndTime = time.Now()
			logger.Error(logLevel, logDetailPayload, fmt.Sprintf("Build panic: %v", r))
		}
	}()

	// Parse manifest body
	var manifestData ManifestData
	if err := json.Unmarshal([]byte(manifest.Body), &manifestData); err != nil {
		job.Status = BuildStatusFailed
		job.ErrorMessage = fmt.Sprintf("Failed to parse manifest body: %v", err)
		job.EndTime = time.Now()
		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	// Create job directory
	jobDir := filepath.Join(s.buildDir, job.ID)
	if err := os.MkdirAll(jobDir, 0755); err != nil {
		job.Status = BuildStatusFailed
		job.ErrorMessage = fmt.Sprintf("Failed to create job directory: %v", err)
		job.EndTime = time.Now()
		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	// Process parameters and create config file
	paramFile, err := s.processParameters(job, &manifestData, jobDir)
	if err != nil {
		job.Status = BuildStatusFailed
		job.ErrorMessage = fmt.Sprintf("Failed to process parameters: %v", err)
		job.EndTime = time.Now()
		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	fmt.Fprintf(&logBuffer, "Generated parameter file: %s\n", paramFile)

	// Launch Docker container to build the payload
	if err := s.runDockerBuild(ctx, job, &manifestData, jobDir, &logBuffer); err != nil {
		job.Status = BuildStatusFailed
		job.ErrorMessage = fmt.Sprintf("Build failed: %v", err)
		job.EndTime = time.Now()
		logger.Error(logLevel, logDetailPayload, job.ErrorMessage)
		return
	}

	// Set the output path
	outputFile := manifestData.PayloadBuildConfig.OutputFile
	if outputFile == "" {
		outputFile = fmt.Sprintf("payload-%s-%s", job.PayloadID, job.Architecture)
	}

	job.OutputPath = filepath.Join(jobDir, "output", outputFile)
	job.Status = BuildStatusCompleted
	job.EndTime = time.Now()

	logger.Info(logLevel, logDetailPayload, fmt.Sprintf("Build job %s completed successfully", job.ID))
}

// processParameters validates and processes build parameters
func (s *PayloadBuildService) processParameters(job *BuildJob, manifestData *ManifestData, jobDir string) (string, error) {
	// Validate required parameters
	for _, param := range manifestData.Parameters.CustomParameters {
		if param.Required {
			value, exists := job.Parameters[param.ParameterName]
			if !exists || value == "" {
				// Use default value if provided
				if param.DefaultValue != "" {
					job.Parameters[param.ParameterName] = param.DefaultValue
				} else {
					return "", fmt.Errorf("missing required parameter: %s", param.ParameterName)
				}
			}
		}
	}

	// Add listener ID if specified
	if manifestData.Parameters.ListenerID != "" {
		job.Parameters["listener_id"] = manifestData.Parameters.ListenerID
	}

	// Create parameters JSON file
	paramsData, err := json.MarshalIndent(job.Parameters, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal parameters: %w", err)
	}

	// Write parameters to file
	paramsFile := filepath.Join(jobDir, "parameters.json")
	if err := os.WriteFile(paramsFile, paramsData, 0644); err != nil {
		return "", fmt.Errorf("failed to write parameters file: %w", err)
	}

	return paramsFile, nil
}

// runDockerBuild executes the Docker build process
func (s *PayloadBuildService) runDockerBuild(ctx context.Context, job *BuildJob, manifestData *ManifestData, jobDir string, logBuffer io.Writer) error {
	// Create output directory
	outputDir := filepath.Join(jobDir, "output")
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Prepare Docker command
	dockerImage := manifestData.PayloadBuildConfig.DockerImage
	if dockerImage == "" {
		dockerImage = "golang:latest" // Default to golang if not specified
	}

	// Prepare build command
	buildArgs := []string{
		"run", "--rm",
		"-v", fmt.Sprintf("%s:/src", manifestData.SourcePath),
		"-v", fmt.Sprintf("%s:/parameters.json", filepath.Join(jobDir, "parameters.json")),
		"-v", fmt.Sprintf("%s:/output", outputDir),
		"-e", fmt.Sprintf("GOOS=%s", getGOOS(job.Architecture)),
		"-e", fmt.Sprintf("GOARCH=%s", getGOARCH(job.Architecture)),
		"-e", "CGO_ENABLED=0",
		"-w", "/src",
		dockerImage,
	}

	// Add custom build command
	buildCmd := fmt.Sprintf("go build %s -o /output/%s .",
		manifestData.PayloadBuildConfig.BuildArgs,
		manifestData.PayloadBuildConfig.OutputFile)

	buildArgs = append(buildArgs, "sh", "-c", buildCmd)

	// Log the command
	fmt.Fprintf(logBuffer, "Running Docker build: docker %s\n", strings.Join(buildArgs, " "))

	// Execute docker command
	cmd := exec.CommandContext(ctx, "docker", buildArgs...)
	cmd.Stdout = logBuffer
	cmd.Stderr = logBuffer

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker build failed: %w", err)
	}

	// Verify the output file exists
	outputFile := filepath.Join(outputDir, manifestData.PayloadBuildConfig.OutputFile)
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		return fmt.Errorf("build completed but output file not found: %s", outputFile)
	}

	return nil
}

// Template processing for source code
func (s *PayloadBuildService) processTemplates(sourceDir string, parameters map[string]string) error {
	// Walk through all files in the source directory
	return filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Only process template files (e.g., *.go.tmpl, *.c.tmpl)
		if !strings.HasSuffix(path, ".tmpl") {
			return nil
		}

		// Read the template file
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read template file %s: %w", path, err)
		}

		// Parse and execute the template
		tmpl, err := template.New(filepath.Base(path)).Parse(string(content))
		if err != nil {
			return fmt.Errorf("failed to parse template file %s: %w", path, err)
		}

		// Create the output file (remove .tmpl extension)
		outputPath := strings.TrimSuffix(path, ".tmpl")
		outputFile, err := os.Create(outputPath)
		if err != nil {
			return fmt.Errorf("failed to create output file %s: %w", outputPath, err)
		}
		defer outputFile.Close()

		// Execute the template with the parameters
		if err := tmpl.Execute(outputFile, parameters); err != nil {
			return fmt.Errorf("failed to execute template %s: %w", path, err)
		}

		return nil
	})
}

// Helper functions to map architecture to GOOS/GOARCH
func getGOOS(arch string) string {
	switch {
	case strings.HasPrefix(arch, "win-"):
		return "windows"
	case strings.HasPrefix(arch, "linux-"):
		return "linux"
	case strings.HasPrefix(arch, "darwin-"):
		return "darwin"
	default:
		return "windows" // Default to Windows
	}
}

func getGOARCH(arch string) string {
	switch {
	case strings.HasSuffix(arch, "-x64"):
		return "amd64"
	case strings.HasSuffix(arch, "-x86"):
		return "386"
	case strings.HasSuffix(arch, "-arm64"):
		return "arm64"
	default:
		return "amd64" // Default to amd64
	}
}
