package models

import "time"

type PayloadRequest struct {
	ListenerID string `json:"listenerID"`
	//ListenerType   string        `json:"listenerType"`
	Sleep        int       `json:"sleep"`
	Jitter       int       `json:"jitter"` // Jitter as a percentage
	StartDate    time.Time `json:"start_date"`
	KillDate     time.Time `json:"kill_date"`
	WorkingHours [2]int    `json:"working_hours"`
}

type PayloadConfig struct {
	ID         string `json:"id"`
	ListenerID string `json:"listenerID"`
	//ListenerType   string        `json:"listenerType"`
	ListenerConfig any       `json:"config"`
	Sleep          int       `json:"sleep"`
	Jitter         int       `json:"jitter"` // Jitter as a percentage
	StartDate      time.Time `json:"start_date"`
	KillDate       time.Time `json:"kill_date"`
	WorkingHours   [2]int    `json:"working_hours"`
}

func IntoPayloadConfig(payloadRequest PayloadRequest) PayloadConfig {
	return PayloadConfig{
		ID:             "",
		ListenerID:     payloadRequest.ListenerID,
		ListenerConfig: nil,
		Sleep:          payloadRequest.Sleep,
		Jitter:         payloadRequest.Jitter,
		StartDate:      payloadRequest.StartDate,
		KillDate:       payloadRequest.KillDate,
		WorkingHours:   payloadRequest.WorkingHours,
	}
}
