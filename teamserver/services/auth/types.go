package auth

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/ksel172/Meduza/teamserver/models"
)

// AuthRequest is the auth request
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse is the auth response
type AuthResponse struct {
	Token string `json:"token"`
}

// UserClaim is the custom claim for token
type UserClaim struct {
	Id   uuid.UUID
	Role models.UserRole
	*jwt.RegisteredClaims
}
