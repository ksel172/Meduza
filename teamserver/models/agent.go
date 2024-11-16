package models

import (
	"time"
)

// Contains all information required for controlling an agent.
type Agent struct {
	ID            string      `json:"id"`
	Name          string      `json:"name"`
	Note          string      `json:"note"`
	Status        string      `json:"status"`
	Config        AgentConfig `json:"config"`
	Info          AgentInfo   `json:"agent_info"`
	FirstCallback time.Time   `json:"first_callback"`
	LastCallback  time.Time   `json:"last_callback"`
	ModifiedAt    time.Time   `json:"modified_at"`
}

// AgentInfo contains information about the agent computer
type AgentInfo struct {
	UUMOID     string `json:"uumo_id"`
	HostName   string `json:"host_name"`
	IPAddr     string `json:"ip_addr"`
	Username   string `json:"username"`
	SystemInfo string `json:"system_info"`
	OSInfo     string `json:"os_info"`
}

// AgentConfig controls how the agent operates
type AgentConfig struct {
	ID              string            `json:"id"`
	Name            string            `json:"name"`
	CallbackURLs    []string          `json:"callback_urls"`
	RotationType    string            `json:"rotation_type"`
	RotationRetries int               `json:"rotation_retries"`
	Sleep           time.Duration     `json:"sleep"`
	Jitter          int               `json:"jitter"` // Jitter as a percentage
	StartDate       time.Time         `json:"start_date"`
	KillDate        time.Time         `json:"kill_date"`
	WorkingHours    [2]int            `json:"working_hours"`
	Headers         map[string]string `json:"headers"` // Custom headers
}

// AgentTask represents the information of a task sent to an Agent
type AgentTask struct {
	ID       string    `json:"id"`
	AgentID  string    `json:"agent_id"`
	Type     string    `json:"type"`
	Status   string    `json:"status"`
	Module   string    `json:"module"`
	Commmand string    `json:"commmand"`
	Created  time.Time `json:"created"`
	Started  time.Time `json:"started"`
	Finished time.Time `json:"finished"`
}

// AgentCommand represents the information of a command sent to an Agent
type AgentCommand struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	Started    time.Time `json:"started"`
	Completed  time.Time `json:"completed"`
	Parameters []string  `json:"parameters"`
	Output     string    `json:"output"`
}
