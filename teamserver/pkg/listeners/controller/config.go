package controller

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"

	_ "github.com/go-playground/validator/v10"
)

type ControllerConfig struct {
	ID         string           `json:"id"`          // Controller ID
	ServerURL  string           `json:"server_url"`  // C2 server URL
	APIKey     string           `json:"api_key"`     // API Key to register controller with the c2 server
	ListenAddr string           `json:"listen_addr"` // Endpoint controller's REST API is listening at
	Heartbeat  int              `json:"heartbeat"`   // Controller heartbeat frequency
	Listeners  []ListenerConfig `json:"listeners"`   // Controller managed listeners when it was shutdown
}

type ListenerConfig struct {

	// metadata about the listener
	ID   string
	Kind string // http, tcp, smb, custom, etc

	// listener specs shared by both deployment types
	Host   string
	Port   int
	Status string // running, stopped etc

	// external deployment specific specs
	Heartbeat int

	// fixed specs, cannot be modified after set unless listener is restarted
	Lifecycle  string `json:"lifecycle" validate:"oneof:scheduled managed"` // scheduled || managed
	Deployment string `json:"deployment" validate:"oneof:local external"`
}

func NewControllerConfig(configBase64 string) (ControllerConfig, error) {
	// Read json base64 encoded argument
	configBytes, err := base64.StdEncoding.DecodeString(configBase64)
	if err != nil {
		return ControllerConfig{}, fmt.Errorf("failed to decode ControllerConfig: %v", err)
	}
	var config ControllerConfig
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return ControllerConfig{}, fmt.Errorf("failed to marshal into ControllerConfig: %v", err)
	}

	// validate
	if err := config.validate(); err != nil {
		return ControllerConfig{}, fmt.Errorf("config validation error: %v", err)
	}

	return config, nil
}

func (c *ControllerConfig) validate() error {
	// Accumulate errors
	var errs []error

	// Validate required config fields
	if c.ServerURL == "" || c.ID == "" {
		errs = append(errs, fmt.Errorf("invalid config: missing required fields"))
	}

	if c.Heartbeat < 30 {
		errs = append(errs, errors.New("heartbeat cannot be less than 30"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c *ListenerConfig) validate() error {
	var errs []error

	// Validate lifecycle and deployment combinations
	if c.Deployment == DeploymentLocal && c.Lifecycle == LifecycleScheduled {
		errs = append(errs, errors.New("local deployments cannot be scheduled"))
	}

	if c.Heartbeat < 30 {
		errs = append(errs, errors.New("heartbeat cannot be less than 30"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
