package models

import "time"

type PayloadConfig struct {
	ID             string        `json:"id"`
	ListenerID     string        `json:"listenerID"`
	ListenerType   string        `json:"listenerType"`
	ListenerConfig any           `json:"config"`
	Sleep          time.Duration `json:"sleep"`
	Jitter         int           `json:"jitter"` // Jitter as a percentage
	StartDate      time.Time     `json:"start_date"`
	KillDate       time.Time     `json:"kill_date"`
	WorkingHours   [2]int        `json:"working_hours"`
}

func IntoPayloadConfig(agentConfig AgentConfig) PayloadConfig {
	return PayloadConfig{
		ID:             agentConfig.ID,
		ListenerID:     agentConfig.ListenerID,
		ListenerConfig: nil, // Assuming ListenerConfig is not present in AgentConfig
		Sleep:          agentConfig.Sleep,
		Jitter:         agentConfig.Jitter,
		StartDate:      agentConfig.StartDate,
		KillDate:       agentConfig.KillDate,
		WorkingHours:   agentConfig.WorkingHours,
	}
}
