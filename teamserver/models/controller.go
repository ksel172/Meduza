package models

import (
	"time"
)

type ListenerController struct {
	ID             string    `json:"id"`
	CallbackIP     string    `json:"callback_ip"`
	CallbackPort   string    `json:"callback_port"`
	ListenerKind   string    `json:"listener_kind"`
	ListenerConfig string    `json:"listener_config"`
	CreatedAt      time.Time `json:"created_at"`
	LastSeen       time.Time `json:"last_seen"`
	//Status         string `json:"status"`
}
