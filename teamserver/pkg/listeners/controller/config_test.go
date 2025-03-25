package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListenerConfig_Validate(t *testing.T) {
	tests := []struct {
		name   string
		config ListenerConfig
	}{
		{
			name: "local-scheduled",
			config: ListenerConfig{
				Lifecycle:  LifecycleScheduled,
				Deployment: DeploymentLocal,
			},
		},
		{
			name: "heartbeat-lower-than-minimum",
			config: ListenerConfig{
				Lifecycle:  LifecycleScheduled,
				Deployment: DeploymentExternal,
				Heartbeat:  29,
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.config.validate()
			assert.Error(t, err, "expected config validation error")
		})
	}
}
