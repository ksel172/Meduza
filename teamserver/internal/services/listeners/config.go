package services

import (
	"errors"

	_ "github.com/go-playground/validator/v10"
)

// type ListenerConfig struct {

// 	// metadata about the listener
// 	ID   string `json:"id"`
// 	Type string `json:"type"` // http, tcp, smb, custom, etc

// 	Name        string `json:"name"`
// 	Description string `json:"description"`
// 	Status      string `json:"status"` // running, stopped etc

// 	// external deployment specific specs
// 	Heartbeat int `json:"heartbeat"`

// 	// fixed specs, cannot be modified after set unless listener is restarted
// 	Lifecycle  string `json:"lifecycle" validate:"oneof:scheduled managed"` // scheduled || managed
// 	Deployment string `json:"deployment" validate:"oneof:local external"`

// 	// listener specs shared by both deployment types
// 	Host string `json:"host"`
// 	Port int    `json:"port"`
// }

type HttpListenerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (h *HttpListenerConfig) Validate() error {
	var errs []error

	if h.Host == "" {
		errs = append(errs, errors.New("host cannot be empty"))
	}

	if h.Port < 1 || h.Port > 65535 {
		errs = append(errs, errors.New("port must be between 1 and 65535"))
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (c *Listener) validate() error {
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
