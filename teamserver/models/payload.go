package models

import "time"

const (
	// URL parameter constants
	ParamPayloadID string = "payload_id"
)

type PayloadRequest struct {
	PayloadName       string    `json:"payload_name" validate:"required"`
	ListenerID        string    `json:"listener_id" validate:"required"`
	Arch              string    `json:"architecture" validate:"required,oneof=win-x64 win-x86 linux-x64 linux-x86"`
	SelfContained     bool      `json:"self_contained" validate:"required,oneof=true false"`
	Sleep             uint      `json:"sleep" validate:"required"`
	Jitter            uint      `json:"jitter" validate:"required"`
	StartDate         time.Time `json:"start_date" validate:"required"`
	KillDate          time.Time `json:"kill_date" validate:"required"`
	WorkingHoursStart uint8     `json:"working_hours_start" validate:"required"`
	WorkingHoursEnd   uint8     `json:"working_hours_end" validate:"required"`
}

type PayloadConfig struct {
	PayloadID         string    `json:"payload_id"`
	PayloadName       string    `json:"payload_name"`
	ConfigID          string    `json:"config_id"`
	ListenerID        string    `json:"listener_id"`
	PublicKey         []byte    `json:"-"`
	PrivateKey        []byte    `json:"-"`
	Token             string    `json:"token"`
	Arch              string    `json:"architecture"`
	ListenerConfig    any       `json:"config"`
	Sleep             uint      `json:"sleep"`
	Jitter            uint      `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart uint8     `json:"working_hours_start"`
	WorkingHoursEnd   uint8     `json:"working_hours_end"`
	CreatedAt         time.Time `json:"created_at"`
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
		CreatedAt:         time.Now(),
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
