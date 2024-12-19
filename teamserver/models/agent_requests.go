package models

import (
	"time"

	"github.com/google/uuid"
)

// Contains only the fields that can be updated for any given agent
type UpdateAgentRequest struct {
	Name   string                   `json:"name"`
	Note   string                   `json:"note"`
	Status string                   `json:"status"`
	Config UpdateAgentConfigRequest `json:"config"`
}

// Contains only the fields that can be updated for any given agent configuration
type UpdateAgentConfigRequest struct {
	CallbackURLs    []string          `json:"callback_urls"`
	RotationType    string            `json:"rotation_type"`
	RotationRetries int               `json:"rotation_retries"`
	Sleep           time.Duration     `json:"sleep"`
	Jitter          int               `json:"jitter"`
	StartDate       time.Time         `json:"start_date"`
	KillDate        time.Time         `json:"kill_date"`
	WorkingHours    [2]int            `json:"working_hours"`
	Headers         map[string]string `json:"headers"`
}

// Conversion from UpdateAgentRequest to Agent
func (uar UpdateAgentRequest) IntoAgent(agent Agent) Agent {

	// Update main fields
	agent.Name = uar.Name
	agent.Note = uar.Note
	agent.Status = uar.Status

	// Update config fields
	agent.Config = uar.Config.IntoAgentConfig(agent.Config)

	// Set last modified time
	agent.ModifiedAt = time.Now()

	return agent
}

// Conversion from UpdateAgentConfigRequest to AgentConfig
func (uacr UpdateAgentConfigRequest) IntoAgentConfig(agentConfig AgentConfig) AgentConfig {

	agentConfig.CallbackURLs = uacr.CallbackURLs
	agentConfig.RotationType = uacr.RotationType
	agentConfig.RotationRetries = uacr.RotationRetries
	agentConfig.Sleep = uacr.Sleep
	agentConfig.Jitter = uacr.Jitter
	agentConfig.StartDate = uacr.StartDate
	agentConfig.KillDate = uacr.KillDate
	agentConfig.WorkingHours = uacr.WorkingHours
	agentConfig.Headers = uacr.Headers

	return agentConfig
}

// AgentTask request
type AgentTaskRequest struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Module  string `json:"module"`
	Command string `json:"command"`
}

// Initializes an AgentTask
func NewAgentTaskRequest() AgentTaskRequest {
	return AgentTaskRequest{}
}

// Returns an AgentTask model from an AgentTaskRequest
func (agr AgentTaskRequest) IntoAgentTask() AgentTask {
	return AgentTask{
		ID:      uuid.New().String(),
		Type:    agr.Type,
		Status:  agr.Status,
		Module:  agr.Module,
		Command: agr.Command,
		Created: time.Now(),
	}
}
