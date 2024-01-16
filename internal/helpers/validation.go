package helpers

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate *validator.Validate = validator.New(validator.WithRequiredStructEnabled())

func MapValidationErrors(err error) []map[string]string {
	var errors []map[string]string

	for _, err := range err.(validator.ValidationErrors) {
		errors = append(errors, map[string]string{
			"field":  strings.ToLower(err.Field()),
			"reason": err.Tag(),
		})
	}

	return errors
}
