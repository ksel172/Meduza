package models

import "time"

const (
	ParamCertificateID = "certificateID"
)

// Certificate represents an SSL/TLS certificate or key
type Certificate struct {
	ID        string    `json:"id" db:"id"`
	Type      string    `json:"type" db:"type"`         // "cert" or "key"
	Path      string    `json:"path" db:"path"`         // File path on server
	Filename  string    `json:"filename" db:"filename"` // Original filename
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
