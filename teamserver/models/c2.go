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
