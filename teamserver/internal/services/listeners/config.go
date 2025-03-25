package services

import (
	"errors"

	_ "github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ListenerConfig struct {

	// metadata about the listener
	ID   uuid.UUID `json:"id"`
	Type string    `json:"type"` // http, tcp, smb, custom, etc

	// listener specs shared by both deployment types
	Host   string `json:"host"`
	Port   int    `json:"port"`
	Status string `json:"status"` // running, stopped etc

	// external deployment specific specs
	Heartbeat int `json:"heartbeat"`

	// fixed specs, cannot be modified after set unless listener is restarted
	Lifecycle  string `json:"lifecycle" validate:"oneof:scheduled managed"` // scheduled || managed
	Deployment string `json:"deployment" validate:"oneof:local external"`
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
