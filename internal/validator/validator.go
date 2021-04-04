package validator

import (
	"github.com/go-playground/validator/v10"
)

// Validator describes the Validator.
type Validator struct {
	validator *validator.Validate
}

// nolint
// Validate satisfies the echo.Validator interface.
func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return err
	}

	return nil
}

// New creates a new Validator.
func New() *Validator {
	return &Validator{
		validator: validator.New(),
	}
}
