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
	AgentID       string      `json:"agent_id" binding:"required,uuid"`
	Name          string      `json:"name" binding:"omitempty,max=100"`
	Note          string      `json:"note"`
	Status        AgentStatus `json:"status"`
	ConfigID      string      `json:"config_id,omitempty"`
	FirstCallback time.Time   `json:"first_callback"`
	LastCallback  time.Time   `json:"last_callback"`
	ModifiedAt    time.Time   `json:"modified_at"`
}

// AgentInfo contains information about the agent computer
type AgentInfo struct {
	AgentID    string `json:"agent_id" binding:"required,uuid"`
	HostName   string `json:"host_name"`
	IPAddress  string `json:"ip_address"`
	Username   string `json:"username"`
	SystemInfo string `json:"system_info"`
	OSInfo     string `json:"os_info"`
}

// AgentConfig controls how the agent operates
type AgentConfig struct {
	ConfigID          string    `json:"config_id"`
	ListenerID        string    `json:"listener_id"`
	Arch              string    `json:"architecture"`
	Sleep             uint      `json:"sleep"`
	Jitter            uint      `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart uint8     `json:"working_hours_start"`
	WorkingHoursEnd   uint8     `json:"working_hours_end"`
}

// AgentTask represents the information of a task sent to an Agent
type AgentTask struct {
	AgentID  string          `json:"agent_id"`
	TaskID   string          `json:"task_id"`
	Type     AgentTaskType   `json:"type"`
	Status   AgentTaskStatus `json:"status"`
	Module   string          `json:"module"`
	Command  AgentCommand    `json:"command"`
	Created  time.Time       `json:"created"`
	Started  time.Time       `json:"started"`
	Finished time.Time       `json:"finished"`
}

// AgentCommand represents the information of a command sent to an Agent
type AgentCommand struct {
	AgentID    string    `json:"agent_id"`
	Name       string    `json:"name"`
	Started    time.Time `json:"started"`
	Completed  time.Time `json:"completed"`
	Parameters []string  `json:"parameters"`
	Output     string    `json:"output"`
}

type AgentTaskType uint8

const (
	LoadAssembly AgentTaskType = iota
	UnloadAssembly
	AgentCommandType
	ShellCommand
	ModuleCommand
	HelpCommand
	SetDelay
	SetJitter
	GetTasks
	KillTasks
	Exit
	Unknown
)

type AgentTaskStatus uint8

const (
	TaskUninitialized AgentTaskStatus = iota
	TaskQueued
	TaskSent
	TaskRunning
	TaskComplete
	TaskFailed
	TaskAborted
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
