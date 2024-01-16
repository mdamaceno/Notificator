package helpers

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestValidationHelper(t *testing.T) {
	t.Run("MapValidationErrors", func(t *testing.T) {
		t.Run("should return an array of ValidationErrors in a map", func(t *testing.T) {
			structToValidate := struct {
				Name string `validate:"required"`
				Age  int    `validate:"required"`
			}{
				Name: "",
				Age:  0,
			}

			err := validator.New().Struct(structToValidate)
			result := MapValidationErrors(err)

			assert.Equal(t, 2, len(result))
		})
	})
}
