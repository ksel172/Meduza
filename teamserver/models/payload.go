package models

import (
	"time"

	"github.com/google/uuid"
)

const (
	// URL parameter constants
	ParamPayloadID    string = "id"
	ParamPayloadToken string = "token"
)

// TODO: add payload supported protocol kinds to validate listener

// PayloadConfig represents a payload
type PayloadConfig struct {
	ID         string `json:"id"`
	ListenerID string `json:"listener_id"` // FK to listeners.ID
	ConfigID   string `json:"config_id"`   // FK to agent_config.ID
	Name       string `json:"name"`
	Arch       string `json:"architecture"`

	// Fields used by agents for checkin
	PublicKey  []byte `json:"-"`
	PrivateKey []byte `json:"-"`
	Token      string `json:"-"`

	CreatedAt time.Time `json:"created_at"`
}

// PayloadRequest represents the data the user sends to create a PayloadConfig
type PayloadRequest struct {
	ListenerID  string `json:"listener_id" validate:"required"` // The listener that is running this payload
	ConfigID    string `json:"config_id" validate:"required"`   // The configuration to use for the payload agents initially
	PayloadName string `json:"name" validate:"required"`
	Arch        string `json:"architecture" validate:"required,oneof=win-x64 win-x86 linux-x64 linux-x86"`
}

// IntoPayloadConfig is the function to convert a PayloadRequest into a PayloadConfig
func IntoPayloadConfig(payloadRequest PayloadRequest) PayloadConfig {
	return PayloadConfig{
		ID:         uuid.New().String(),
		Name:       payloadRequest.PayloadName,
		ConfigID:   payloadRequest.ConfigID,
		ListenerID: payloadRequest.ListenerID,
		Arch:       payloadRequest.Arch,
		CreatedAt:  time.Now(),
	}
}
