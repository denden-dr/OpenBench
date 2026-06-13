package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct validates a struct using the go-playground/validator package.
// Returns a formatted error string listing the failed validation rules.
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err == nil {
		return nil
	}

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		var errMsgs []string
		for _, e := range validationErrors {
			errMsgs = append(errMsgs, fmt.Sprintf("field '%s' failed on tag '%s'", e.Field(), e.Tag()))
		}
		return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, ", "))
	}

	return err
}
