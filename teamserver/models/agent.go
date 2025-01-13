package models

import (
	"time"
)

// AgentTaskId is
type AgentTaskId = string

const (
	ParamAgentID AgentTaskId = "agent_id" // Agents Id
	ParamTaskID  AgentTaskId = "task_id"  // Tasks Id of a agent.
)

// Contains all information required for controlling an agent.
type Agent struct {
	AgentID       string    `json:"agent_id"`
	Name          string    `json:"name"`
	Note          string    `json:"note"`
	Status        string    `json:"status"`
	ConfigID      string    `json:"config_id"`
	Info          AgentInfo `json:"agent_info"`
	FirstCallback time.Time `json:"first_callback"`
	LastCallback  time.Time `json:"last_callback"`
	ModifiedAt    time.Time `json:"modified_at"`
}

// AgentInfo contains information about the agent computer
type AgentInfo struct {
	AgentID    string `json:"agent_id"`
	HostName   string `json:"host_name"`
	IPAddress  string `json:"ip_address"`
	Username   string `json:"username"`
	SystemInfo string `json:"system_info"`
	OSInfo     string `json:"os_info"`
}

// AgentConfig controls how the agent operates
type AgentConfig struct {
	ConfigID          string    `json:"agent_id"`
	ListenerID        string    `json:"listener_id"`
	Arch              string    `json:"architecture"`
	Sleep             int       `json:"sleep"`
	Jitter            int       `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart int       `json:"working_hours_start"`
	WorkingHoursEnd   int       `json:"working_hours_end"`
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

type AgentTaskType string

const (
	LoadAssembly     AgentTaskType = "LoadAssembly"
	UnloadAssembly   AgentTaskType = "UnloadAssembly"
	AgentCommandType AgentTaskType = "AgentCommand"
	ShellCommand     AgentTaskType = "ShellCommand"
	HelpCommand      AgentTaskType = "HelpCommand"
	SetDelay         AgentTaskType = "SetDelay"
	SetJitter        AgentTaskType = "SetJitter"
	GetTasks         AgentTaskType = "GetTasks"
	KillTasks        AgentTaskType = "KillTasks"
	Exit             AgentTaskType = "Exit"
	Unknown          AgentTaskType = "Unknown"
)

type AgentTaskStatus string

const (
	Uninitialized AgentTaskStatus = "Uninitialized"
	Queued        AgentTaskStatus = "Queued"
	Sent          AgentTaskStatus = "Sent"
	Running       AgentTaskStatus = "Running"
	Complete      AgentTaskStatus = "Complete"
	Failed        AgentTaskStatus = "Failed"
	Aborted       AgentTaskStatus = "Aborted"
)
