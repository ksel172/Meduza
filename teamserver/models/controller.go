package models

const (
	// URL parameter constants
	ParamControllerID string = "controller_id"
)

type ControllerRegistration struct {
	ID       string `json:"id"`
	Endpoint string `json:"endpoint"`
}

type KeyPair struct {
	PublicKey  string `json:"PublicKey"`
	PrivateKey string `json:"PrivateKey"`
}

type HeartbeatRequest struct {
	Timestamp int64             `json:"timestamp"`
	Listeners map[string]string `json:"listeners"`
}
