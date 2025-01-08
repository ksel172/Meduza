package models

import "time"

type PayloadRequest struct {
	ListenerID        string    `json:"listenerID"`
	Sleep             int       `json:"sleep"`
	Jitter            int       `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart int       `json:"working_hours_start"`
	WorkingHoursEnd   int       `json:"working_hours_end"`
	//ListenerType   string        `json:"listenerType"`
}

type PayloadConfig struct {
	AgentID           string    `json:"id"`
	ListenerID        string    `json:"listenerID"`
	ListenerConfig    any       `json:"config"`
	Sleep             int       `json:"sleep"`
	Jitter            int       `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart int       `json:"working_hours_start"`
	WorkingHoursEnd   int       `json:"working_hours_end"`
	//ListenerType   string        `json:"listenerType"`
}

func IntoPayloadConfig(payloadRequest PayloadRequest) PayloadConfig {
	return PayloadConfig{
		AgentID:           "",
		ListenerID:        payloadRequest.ListenerID,
		ListenerConfig:    nil,
		Sleep:             payloadRequest.Sleep,
		Jitter:            payloadRequest.Jitter,
		StartDate:         payloadRequest.StartDate,
		KillDate:          payloadRequest.KillDate,
		WorkingHoursStart: payloadRequest.WorkingHoursStart,
		WorkingHoursEnd:   payloadRequest.WorkingHoursEnd,
	}
}

func IntoAgentConfig(payloadConfig PayloadConfig) AgentConfig {
	return AgentConfig{
		AgentID:           payloadConfig.AgentID,
		ListenerID:        payloadConfig.ListenerID,
		Sleep:             payloadConfig.Sleep,
		Jitter:            payloadConfig.Jitter,
		StartDate:         payloadConfig.StartDate,
		KillDate:          payloadConfig.KillDate,
		WorkingHoursStart: payloadConfig.WorkingHoursStart,
		WorkingHoursEnd:   payloadConfig.WorkingHoursEnd,
	}
}
