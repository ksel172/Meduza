package models

import (
	"time"
)

const (
	// URL parameter constants
	ParamAgentID string = "agent_id"
	ParamTaskID  string = "task_id"
)

// Contains all information required for controlling an agent.
type Agent struct {
	ID            string      `json:"id"`
	PayloadID     string      `json:"payload_id"`          // FK to payloads.ID
	ConfigID      string      `json:"config_id,omitempty"` // FK to agent_config.ID
	Name          string      `json:"name"`
	Note          string      `json:"note"`
	Status        AgentStatus `json:"status"`
	FirstCallback time.Time   `json:"first_callback"`
	LastCallback  time.Time   `json:"last_callback"`
	ModifiedAt    time.Time   `json:"modified_at"`

	AgentInfo `json:"agent_info"` // embeds AgentInfo, in db, table agents_info
}

// AgentInfo contains information about the agent device
// Inherits ID from Agent
type AgentInfo struct {
	AgentID    string `json:"agent_id"`
	HostName   string `json:"hostname"`
	IPAddress  string `json:"ip_address"`
	Username   string `json:"username"`
	SystemInfo string `json:"system_info"`
	OSInfo     string `json:"os_info"`
}

type AgentInfoRequest struct {
	HostName   string `json:"hostname"`
	IPAddress  string `json:"ip_address"`
	Username   string `json:"username"`
	SystemInfo string `json:"system_info"`
	OSInfo     string `json:"os_info"`
}

// AgentConfig controls how the agent operates
type AgentConfig struct {
	ID                string    `json:"id"`
	Sleep             uint      `json:"sleep"`
	Jitter            uint      `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart uint8     `json:"working_hours_start"`
	WorkingHoursEnd   uint8     `json:"working_hours_end"`
}

// AgentTask represents the information of a task sent to an Agent
type AgentTask struct {
	ID         string          `json:"task_id"`
	AgentID    string          `json:"agent_id"` // FK to agents.ID
	Type       AgentTaskType   `json:"type"`
	Status     AgentTaskStatus `json:"status"`
	Module     string          `json:"module"`
	Command    AgentCommand    `json:"command"`
	CreatedAt  time.Time       `json:"created_at"`
	StartedAt  time.Time       `json:"started_at"`
	FinishedAt time.Time       `json:"finished_at"`
}

// AgentCommand represents the information of a command sent to an Agent
type AgentCommand struct {
	Name       string    `json:"name"`
	Started    time.Time `json:"started"`
	Completed  time.Time `json:"completed"`
	Parameters []string  `json:"parameters"`
	Output     string    `json:"output"`
}

type AgentTaskType uint8

const (
	TaskLoadAssembly AgentTaskType = iota
	TaskUnloadAssembly
	TaskAgentCommand
	TaskShellCommand
	TaskModuleCommand
	TaskHelpCommand
	TaskSetDelay
	TaskSetJitter
	TaskGetTasks
	TaskKillTasks
	TaskExit
	TaskUnknown
)

type AgentTaskStatus uint8

const (
	TaskStatusUninitialized AgentTaskStatus = iota
	TaskStatusQueued
	TaskStatusSent
	TaskStatusRunning
	TaskStatusComplete
	TaskStatusFailed
	TaskStatusAborted
)

type AgentStatus uint8

const (
	AgentUninitialized AgentStatus = iota
	AgentStage0
	AgentStage1
	AgentStage2
	AgentActive
	AgentLost
	AgentExited
	AgentDisconnected
	AgentHidden
)
