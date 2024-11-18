package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated-at"`
}

type ResUser struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username" validate:"required,min=6,max=20"`
	PasswordHash string    `json:"password" validate:"required,min=6,max=20"`
	Role         UserRole  `json:"role" validate:"required,oneof=admin moderator client visitor"`
	CreateBy     time.Time `json:"created_by omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRole string

const (
	ADMIN     UserRole = "admin"
	MODERATOR UserRole = "moderator"
	CLIENT    UserRole = "client"
	VISITOR   UserRole = "visitor"
)
