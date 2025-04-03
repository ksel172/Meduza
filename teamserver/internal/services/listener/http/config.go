package http_listener

import (
	"errors"
	"time"
)

// Concrete listener implementation configurations
type HTTPListenerConfig struct {
	Host         string        `json:"host"`
	Port         int           `json:"port"`
	EnableTLS    bool          `json:"enable_tls"`
	CertPath     string        `json:"cert_path"`
	KeyPath      string        `json:"key_path"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
}

func (h *HTTPListenerConfig) ValidateConfig() error {
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
