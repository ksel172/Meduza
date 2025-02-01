package models

import (
	"time"
)

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type ResUser struct {
	ID           string    `json:"id"`
	Username     string    `json:"username" validate:"alphanum,required,min=6,max=20"`
	PasswordHash string    `json:"password" validate:"required,min=6"`
	Role         string    `json:"role" validate:"oneof=admin user guest"`
	CreateBy     time.Time `json:"created_by omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserRole string

const (
	ADMIN UserRole = "admin"
	USER  UserRole = "user"
	GUEST UserRole = "guest"
)
