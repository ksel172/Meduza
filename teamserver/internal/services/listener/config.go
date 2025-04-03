package listener

import (
	"errors"

	http_listener "github.com/ksel172/Meduza/teamserver/internal/services/listener/http"
	"github.com/ksel172/Meduza/teamserver/utils"
)

func (l *Listener) ValidateConfig() error {
	var errs []error

	// Validate generic listener config
	if err := l.validateSetDefaults(); err != nil {
		errs = append(errs, err)
	}

	// Validate specific listener kind configs
	switch l.Kind {
	case "http":
		var config http_listener.HTTPListenerConfig
		if err := utils.MapToStruct(l.Config, &config); err != nil {
			errs = append(errs, errors.New("invalid HTTP listener config"))
		} else if err := config.ValidateConfig(); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (l *Listener) validateSetDefaults() error {
	var errs []error

	// Validate lifecycle and deployment combinations
	if l.Deployment == DeploymentLocal && l.Lifecycle == LifecycleScheduled {
		errs = append(errs, errors.New("local deployments cannot be scheduled"))
	}

	if l.Heartbeat == 0 {
		l.Heartbeat = 30
	} else if l.Heartbeat < 30 {
		errs = append(errs, errors.New("heartbeat cannot be less than 30"))
	}

	if l.Name == "" {
		l.Name = "listener-" + utils.RandomString(12)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}
