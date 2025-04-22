package validator

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type Validator struct {
	Errors []ValidationError
}

func New() *Validator {
	return &Validator{
		Errors: []ValidationError{},
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(field, message string) {
	v.Errors = append(v.Errors, ValidationError{Field: field, Message: message})
}

func (v *Validator) Check(condition bool, field, message string) {
	if !condition {
		v.AddError(field, message)
	}
}

func (v *Validator) Required(field, value string) {
	if value == "" {
		v.AddError(field, field+" is required")
	}
}

func (v *Validator) MinValue(field string, value, min float64) {
	if value < min {
		v.AddError(field, fmt.Sprintf("%s must be greater than %f", field, min))
	}
}

func (v *Validator) ValidUUID(field string, value string) (uuid.UUID, bool) {
	if value == "" {
		v.AddError(field, field+" is required")
		return uuid.Nil, false
	}

	id, err := uuid.Parse(value)

	if err != nil {
		v.AddError(field, field+" is not a valid uuid")
		return uuid.Nil, false
	}

	return id, true
}

func (v *Validator) ErrorMessage() string {
	errorMessages := make([]string, len(v.Errors))
	for i, e := range v.Errors {
		errorMessages[i] = e.Message
	}
	return strings.Join(errorMessages, ", ")
}
