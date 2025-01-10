package models

import "time"

const (
	// Windows architectures
	ArchWinX64   = "win-x64"
	ArchWinX86   = "win-x86"
	ArchWinArm64 = "win-arm64"

	// Linux architectures
	ArchLinuxX64   = "linux-x64"
	ArchLinuxArm   = "linux-arm"
	ArchLinuxArm64 = "linux-arm64"
)

// AllArchs returns a slice of all supported architectures
func AllArchs() []string {
	return []string{
		ArchWinX64,
		ArchWinX86,
		ArchWinArm64,
		ArchLinuxX64,
		ArchLinuxArm,
		ArchLinuxArm64,
	}
}

type PayloadRequest struct {
	PayloadName       string    `json:"payload_name"`
	ListenerID        string    `json:"listener_id"`
	Arch              string    `json:"arch"`
	Sleep             int       `json:"sleep"`
	Jitter            int       `json:"jitter"` // Jitter as a percentage
	StartDate         time.Time `json:"start_date"`
	KillDate          time.Time `json:"kill_date"`
	WorkingHoursStart int       `json:"working_hours_start"`
	WorkingHoursEnd   int       `json:"working_hours_end"`
	//ListenerType   string        `json:"listenerType"`
}

type PayloadConfig struct {
	PayloadID         string    `json:"payload_id"`
	PayloadName       string    `json:"payload_name"`
	ConfigID          string    `json:"config_id"`
	ListenerID        string    `json:"listener_id"`
	Arch              string    `json:"arch"`
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
		Sleep:             payloadConfig.Sleep,
		Jitter:            payloadConfig.Jitter,
		StartDate:         payloadConfig.StartDate,
		KillDate:          payloadConfig.KillDate,
		WorkingHoursStart: payloadConfig.WorkingHoursStart,
		WorkingHoursEnd:   payloadConfig.WorkingHoursEnd,
	}
}
