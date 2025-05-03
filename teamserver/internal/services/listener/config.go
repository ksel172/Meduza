package listener

import (
	"errors"

	"github.com/ksel172/Meduza/teamserver/utils"
)

func (l *Listener) ValidateConfig() error {
	var errs []error

	// Validate generic listener config
	if err := l.validateSetDefaults(); err != nil {
		errs = append(errs, err)
	}

	// Validate listener implementation config for local deployments
	if err := l.listener.Validate(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

func (l *Listener) validateSetDefaults() error {
	var errs []error

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
