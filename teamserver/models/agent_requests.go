package models

import (
	"time"

	"github.com/google/uuid"
)

// Contains only the fields that can be updated for any given agent
type UpdateAgentRequest struct {
	AgentID    string                   `json:"agent_id"`
	Name       string                   `json:"name"`
	Note       string                   `json:"note"`
	Status     string                   `json:"status"`
	Config     UpdateAgentConfigRequest `json:"config"`
	ModifiedAt time.Time                `json:"modified_at"`
}

// Contains only the fields that can be updated for any given agent configuration
type UpdateAgentConfigRequest struct {
	Sleep             int       `json:"sleep"`
	Jitter            int       `json:"jitter"`
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart int       `json:"working_hours_start"`
	WorkingHoursEnd   int       `json:"working_hours_end"`
}

// Conversion from UpdateAgentConfigRequest to AgentConfig
func (uacr UpdateAgentConfigRequest) IntoAgentConfig(agentConfig AgentConfig) AgentConfig {

	agentConfig.Sleep = uacr.Sleep
	agentConfig.Jitter = uacr.Jitter
	agentConfig.StartDate = uacr.StartDate
	agentConfig.KillDate = uacr.KillDate
	agentConfig.WorkingHoursStart = uacr.WorkingHoursStart
	agentConfig.WorkingHoursEnd = uacr.WorkingHoursEnd

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
		AgentID: uuid.New().String(),
		Type:    agr.Type,
		Status:  agr.Status,
		Module:  agr.Module,
		Command: agr.Command,
		Created: time.Now(),
	}
}
