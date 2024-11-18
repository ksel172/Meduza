package models

// Just Making Life easy
type UserPResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}
