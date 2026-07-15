package utils

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct's fields using the validator/v10 library tags.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}
