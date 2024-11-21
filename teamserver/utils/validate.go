package utils

import (
	"github.com/go-playground/validator/v10"
)

// Validate struct holds a validator instance
type Validate struct {
	validate *validator.Validate
}

// NewValidatorService creates and returns a new Validate instance
func NewValidatorService() *Validate {
	return &Validate{
		validate: validator.New(),
	}
}

// ValidateStruct validates a struct and returns an error if validation fails
func (v *Validate) ValidateStruct(s interface{}) error {
	return v.validate.Struct(s)
}
