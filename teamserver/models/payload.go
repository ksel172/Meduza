package models

import (
	"time"

	"encoding/json"
	"errors"
	"fmt"
)

const (
	ParamPayloadID    string = "id"
	ParamManifestID   string = "manifest_id"
	ParamPayloadJobID string = "payload_job_id"
	ParamPayloadToken string = "token"
)

type Payload struct {
	ID         string `json:"id"`
	ListenerID string `json:"listener_id"` // FK to listeners.ID
	ConfigID   string `json:"config_id"`   // FK to agent_config.ID
	ManifestID string `json:"manifest_id"` // FK to payload_manifest.ID
	Name       string `json:"name"`
	Arch       string `json:"architecture"`

	// Fields used by agents for checkin
	PublicKey  []byte `json:"-"`
	PrivateKey []byte `json:"-"`
	Token      string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}

type PayloadManifestV1 struct {
	ID              string    `json:"id"`
	ManifestVersion float32   `json:"manifest_version"`
	Body            string    `json:"body"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

type PayloadManifestBodyV1 struct {
	PayloadName        string `json:"name"`
	PayloadVersion     string `json:"version"`
	PayloadAuthor      string `json:"author"`
	PayloadDescription string `json:"description"`

	PayloadBuildConfig PayloadBuildConfig     `json:"payload_build_config"` // PayloadBuildConfig represents basic params for building a payload
	PayloadParameters  PayloadBuildParameters `json:"parameters"`           // PayloadBuildParameters represents the parameters for the payload

	SourcePath string `json:"source_path"` // SourcePath is the path to the source code for the payload
}

type PayloadBuildConfig struct {
	DockerImage    string   `json:"docker_image"`
	BuildArgs      string   `json:"build_args"`
	OutputPath     string   `json:"output_path"`
	OutputFile     string   `json:"output_file"`
	SupportedArchs []string `json:"supported_arch"`
}

type PayloadBuildParameters struct {
	// Necessary parameters for building a payload
	ListenerID string `json:"listener_id"`

	// Custom parameters for building a payload
	CustomParameters []BuildParameter `json:"custom_parameters"`
}

// TODO: Need to add validation for the parameters (especially for the data type)
type BuildParameter struct {
	DataType      string `json:"data_type"`
	ParameterName string `json:"parameter_name"`
	Description   string `json:"description"`
	DefaultValue  string `json:"default_value"`
	Required      bool   `json:"required"`
}

type PayloadJob struct {
	ID           string                 `json:"id"`
	PayloadID    string                 `json:"payload_id"`
	Status       string                 `json:"status"`
	Architecture string                 `json:"architecture"`
	Parameters   map[string]interface{} `json:"parameters"`
	StartTime    time.Time              `json:"start_time"`
	EndTime      time.Time              `json:"end_time,omitempty"`
	OutputPath   string                 `json:"output_path,omitempty"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	BuildLog     string                 `json:"build_log,omitempty"`
}

func (m *PayloadManifestV1) Validate() error {
	var errs []error

	// Check base required fields
	if m.ManifestVersion == 0 {
		errs = append(errs, errors.New("manifest_version cannot be empty"))
	}

	if m.Body == "" {
		errs = append(errs, errors.New("manifest body cannot be empty"))
	}

	// Validate body according to manifest version
	if m.ManifestVersion <= 1.0 {
		var body PayloadManifestBodyV1
		if err := json.Unmarshal([]byte(m.Body), &body); err != nil {
			errs = append(errs, fmt.Errorf("failed to unmarshal manifest body: %w", err))
		}

		// Validate V1 body fields
		if body.PayloadName == "" {
			errs = append(errs, errors.New("payload name cannot be empty"))
		}

		if body.PayloadVersion == "" {
			errs = append(errs, errors.New("payload version cannot be empty"))
		}

		if body.PayloadAuthor == "" {
			errs = append(errs, errors.New("payload author cannot be empty"))
		}

		if body.PayloadDescription == "" {
			errs = append(errs, errors.New("payload description cannot be empty"))
		}

		if body.SourcePath == "" {
			errs = append(errs, errors.New("source path cannot be empty"))
		}

		if len(body.PayloadBuildConfig.SupportedArchs) == 0 {
			errs = append(errs, errors.New("at least one supported architecture must be specified"))
		}

		if err := body.PayloadParameters.ValidateCustomParameters(); err != nil {
			errs = append(errs, fmt.Errorf("invalid parameters: %w", err))
		}
	} else {
		errs = append(errs, fmt.Errorf("unsupported manifest version: %f", m.ManifestVersion))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}

func ValidateManifestJSON(data []byte) error {
	var manifest PayloadManifestV1
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("invalid manifest JSON: %w", err)
	}

	return manifest.Validate()
}

// BuildParameterTypes returns allowed data types for build parameters
func BuildParameterTypes() []string {
	return []string{"string", "integer", "boolean", "float"}
}

// ValidateBuildParameter validates a single build parameter
func (p *BuildParameter) ValidateBuildParameter() error {
	// Check required fields
	if p.ParameterName == "" {
		return fmt.Errorf("parameter missing required field: parameter_name")
	}

	// Validate parameter data type
	validTypes := BuildParameterTypes()
	isValidType := false
	for _, t := range validTypes {
		if p.DataType == t {
			isValidType = true
			break
		}
	}

	if !isValidType {
		return fmt.Errorf("invalid parameter data type for %s: %s", p.ParameterName, p.DataType)
	}

	// If parameter is required, it should have a default value
	if p.Required && p.DefaultValue == "" {
		return fmt.Errorf("required parameter %s should have a default value", p.ParameterName)
	}

	return nil
}

// ValidateCustomParameters validates all custom parameters
func (p *PayloadBuildParameters) ValidateCustomParameters() error {
	paramNames := make(map[string]bool)

	for i, param := range p.CustomParameters {
		if err := param.ValidateBuildParameter(); err != nil {
			return fmt.Errorf("invalid parameter at index %d: %w", i, err)
		}

		// Check for duplicate parameter names
		if paramNames[param.ParameterName] {
			return fmt.Errorf("duplicate parameter name: %s", param.ParameterName)
		}
		paramNames[param.ParameterName] = true
	}

	return nil
}

// FromJSON creates a PayloadManifestV1 from JSON data
func PayloadManifestFromJSON(data []byte) (*PayloadManifestV1, error) {
	var manifest PayloadManifestV1
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	if err := manifest.Validate(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// ToJSON converts PayloadManifestV1 to JSON
func (m *PayloadManifestV1) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
