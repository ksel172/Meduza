package models

import (
	"time"

	"github.com/google/uuid"
)

// C2Request represents any request sent to a C2 server by an Agent.
// Valid agent statuses: ["uninitialized", "inactive", "active"]
type C2Request struct {
	Reason      string `json:"reason"`
	AgentID     string `json:"agent_id"`
	ConfigID    string `json:"config_id"`
	AgentStatus string `json:"agent_status"`
	Message     string `json:"message"`
	// Hmac        string `json:"hmac"`
}

type C2RequestReason string

const (
	TaskReason     C2RequestReason = "task"
	ResponseReason C2RequestReason = "response"
	RegisterReason C2RequestReason = "register"
	UnknownReason  C2RequestReason = "unknown"
)

// IsValid checks if the reason is one of the defined constants
func (r C2RequestReason) IsValid() bool {
	switch r {
	case TaskReason, ResponseReason, RegisterReason:
		return true
	default:
		return false
	}
}

// String returns the string representation of the reason
func (r C2RequestReason) String() string {
	return string(r)
}

// Initialize a new C2Request with status uninitialized, for use when creating a new agent
func NewC2Request() C2Request {
	return C2Request{AgentStatus: "uninitialized"} // Default to uninitialized, if not provided
}

// Validates if the C2Request contains valid data
func (r C2Request) Valid() bool {

	// Validate AgentID is uuid
	if _, err := uuid.Parse(r.AgentID); err != nil {
		return false
	}

	// Validate other fields are one of the valid values
	return (r.AgentStatus == "uninitialized" || r.AgentStatus == "inactive" || r.AgentStatus == "active")
}

// Converts a C2Request into a new Agent for registration
func (r C2Request) IntoNewAgent() Agent {
	return Agent{
		AgentID:       r.AgentID, // uuid generated at agent computer, sent with initial checkin request
		ConfigID:      r.ConfigID,
		Status:        r.AgentStatus,
		FirstCallback: time.Now(),
		ModifiedAt:    time.Now(),
	}
}
