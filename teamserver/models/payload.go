package models

import "time"

type PayloadRequest struct {
	PayloadName       string    `json:"payload_name" validate:"required"`
	ListenerID        string    `json:"listener_id" validate:"required"`
	Arch              string    `json:"architecture" validate:"required,oneof=win-x64 win-x86 linux-x64 linux-x86"`
	SelfContained     string    `json:"self_contained" validate:"required,oneof=true false"`
	Sleep             int       `json:"sleep" validate:"required"`
	Jitter            int       `json:"jitter" validate:"required"`
	StartDate         time.Time `json:"start_date" validate:"required"`
	KillDate          time.Time `json:"kill_date" validate:"required"`
	WorkingHoursStart int       `json:"working_hours_start" validate:"required"`
	WorkingHoursEnd   int       `json:"working_hours_end" validate:"required"`
}

type PayloadConfig struct {
	PayloadID         string    `json:"payload_id"`
	PayloadName       string    `json:"payload_name"`
	ConfigID          string    `json:"config_id"`
	ListenerID        string    `json:"listener_id"`
	Token             string    `json:"token"`
	Arch              string    `json:"architecture"`
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
		PayloadName:       payloadRequest.PayloadName,
		ConfigID:          "",
		ListenerID:        payloadRequest.ListenerID,
		Arch:              payloadRequest.Arch,
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
		ConfigID:          payloadConfig.ConfigID,
		ListenerID:        payloadConfig.ListenerID,
		Arch:              payloadConfig.Arch,
		Sleep:             payloadConfig.Sleep,
		Jitter:            payloadConfig.Jitter,
		StartDate:         payloadConfig.StartDate,
		KillDate:          payloadConfig.KillDate,
		WorkingHoursStart: payloadConfig.WorkingHoursStart,
		WorkingHoursEnd:   payloadConfig.WorkingHoursEnd,
	}
}
