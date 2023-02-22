package api

import (
	"github.com/go-playground/validator/v10"
)

// newValidator func for create a new validator for api requests.
func newValidator() *validator.Validate {
	validate := validator.New()
	return validate
}
