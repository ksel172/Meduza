package models

import (
	"time"

	"github.com/google/uuid"
)

// C2Request represents any request sent to a C2 server by an Agent.
type C2Request struct {
	Reason      RequestReason `json:"reason"`
	AgentID     string        `json:"agent_id"`
	ConfigID    string        `json:"config_id"`
	AgentStatus AgentStatus   `json:"agent_status"`
	Message     string        `json:"message"`
}

// Initialize a new C2Request with status uninitialized, for use when creating a new agent
func NewC2Request() C2Request {
	return C2Request{AgentStatus: AgentUninitialized} // Default to uninitialized, if not provided
}

// Validates if the C2Request contains valid data
func (r C2Request) Valid() bool {
	return (r.AgentStatus == AgentUninitialized || r.AgentStatus == AgentActive || r.AgentStatus == AgentExited)
}

// Converts a C2Request into a new Agent for registration
func (r C2Request) IntoNewAgent() Agent {
	return Agent{
		AgentID:       uuid.NewString(),
		ConfigID:      r.ConfigID,
		Status:        r.AgentStatus,
		FirstCallback: time.Now(),
		ModifiedAt:    time.Now(),
	}
}

type RequestReason uint8

const (
	Register RequestReason = iota
	Task
	Response
)
