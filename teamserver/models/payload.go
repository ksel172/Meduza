package models

import (
	"encoding/json"
	"fmt"
	"time"
)

const (
	// URL parameter constants
	ParamPayloadID    string = "id"
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
	ID              string `json:"id"`
	ManifestVersion string `json:"manifest_version"`
	Name            string `json:"name"`
	Version         string `json:"version"`
	Author          string `json:"author"`
	Description     string `json:"description"`

	PayloadBuildConfig PayloadBuildConfig     `json:"payload_build_config"` // PayloadBuildConfig represents basic params for building a payload
	Parameters         PayloadBuildParameters `json:"parameters"`           // PayloadBuildParameters represents the parameters for the payload

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

func (m *PayloadManifestV1) ValidateManifest() error {
	// Check required string fields
	if m.Name == "" {
		return fmt.Errorf("manifest missing required field: name")
	}
	if m.Version == "" {
		return fmt.Errorf("manifest missing required field: version")
	}
	if m.Author == "" {
		return fmt.Errorf("manifest missing required field: author")
	}
	if m.Description == "" {
		return fmt.Errorf("manifest missing required field: description")
	}
	if m.ManifestVersion == "" {
		return fmt.Errorf("manifest missing required field: manifest_version")
	}

	// Validate build config
	if m.PayloadBuildConfig.OutputFile == "" {
		return fmt.Errorf("build config missing required field: output_file")
	}
	if len(m.PayloadBuildConfig.SupportedArchs) == 0 {
		return fmt.Errorf("build config missing required field: supported_arch")
	}

	return nil
}

func ValidateManifestJSON(data []byte) error {
	var manifest PayloadManifestV1
	if err := json.Unmarshal(data, &manifest); err != nil {
		return fmt.Errorf("invalid manifest JSON: %w", err)
	}

	return manifest.ValidateManifest()
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

	if err := manifest.ValidateManifest(); err != nil {
		return nil, err
	}

	if err := manifest.Parameters.ValidateCustomParameters(); err != nil {
		return nil, err
	}

	return &manifest, nil
}

// ToJSON converts PayloadManifestV1 to JSON
func (m *PayloadManifestV1) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

// type PayloadRequest struct {
// 	PayloadName       string    `json:"payload_name" validate:"required"`
// 	ListenerID        string    `json:"listener_id" validate:"required"`
// 	Arch              string    `json:"architecture" validate:"required,oneof=win-x64 win-x86 linux-x64 linux-x86"`
// 	SelfContained     bool      `json:"self_contained" validate:"required,oneof=true false"`
// 	Sleep             uint      `json:"sleep" validate:"required"`
// 	Jitter            uint      `json:"jitter" validate:"required"`
// 	StartDate         time.Time `json:"start_date" validate:"required"`
// 	KillDate          time.Time `json:"kill_date" validate:"required"`
// 	WorkingHoursStart uint8     `json:"working_hours_start" validate:"required"`
// 	WorkingHoursEnd   uint8     `json:"working_hours_end" validate:"required"`
// }

// type PayloadConfig struct {
// 	PayloadID         string    `json:"payload_id"`
// 	PayloadName       string    `json:"payload_name"`
// 	ConfigID          string    `json:"config_id"`
// 	ListenerID        string    `json:"listener_id"`
// 	PublicKey         []byte    `json:"-"`
// 	PrivateKey        []byte    `json:"-"`
// 	Token             string    `json:"token"`
// 	Arch              string    `json:"architecture"`
// 	ListenerConfig    any       `json:"config"`
// 	Sleep             uint      `json:"sleep"`
// 	Jitter            uint      `json:"jitter"` // Jitter as a percentage
// 	StartDate         time.Time `json:"start_date"`
// 	KillDate          time.Time `json:"kill_date"`
// 	WorkingHoursStart uint8     `json:"working_hours_start"`
// 	WorkingHoursEnd   uint8     `json:"working_hours_end"`
// 	CreatedAt         time.Time `json:"created_at"`
// 	//ListenerType   string        `json:"listenerType"`
// }

// func IntoPayloadConfig(payloadRequest PayloadRequest) PayloadConfig {
// 	return PayloadConfig{
// 		PayloadName:       payloadRequest.PayloadName,
// 		ConfigID:          "",
// 		ListenerID:        payloadRequest.ListenerID,
// 		Arch:              payloadRequest.Arch,
// 		ListenerConfig:    nil,
// 		Sleep:             payloadRequest.Sleep,
// 		Jitter:            payloadRequest.Jitter,
// 		StartDate:         payloadRequest.StartDate,
// 		KillDate:          payloadRequest.KillDate,
// 		WorkingHoursStart: payloadRequest.WorkingHoursStart,
// 		WorkingHoursEnd:   payloadRequest.WorkingHoursEnd,
// 		CreatedAt:         time.Now(),
// 	}
// }

// func IntoAgentConfig(payloadConfig PayloadConfig) AgentConfig {
// 	return AgentConfig{
// 		ConfigID:          payloadConfig.ConfigID,
// 		ListenerID:        payloadConfig.ListenerID,
// 		Arch:              payloadConfig.Arch,
// 		Sleep:             payloadConfig.Sleep,
// 		Jitter:            payloadConfig.Jitter,
// 		StartDate:         payloadConfig.StartDate,
// 		KillDate:          payloadConfig.KillDate,
// 		WorkingHoursStart: payloadConfig.WorkingHoursStart,
// 		WorkingHoursEnd:   payloadConfig.WorkingHoursEnd,
// 	}
// }
