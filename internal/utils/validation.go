package utils

import (
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	id_translations "github.com/go-playground/validator/v10/translations/id"
)

var (
	validate = validator.New()
	transEn  ut.Translator
	transId  ut.Translator
)

func init() {
	enLocale := en.New()
	idLocale := id.New()
	uni := ut.New(enLocale, enLocale, idLocale)

	var found bool
	transEn, found = uni.GetTranslator("en")
	if found {
		_ = en_translations.RegisterDefaultTranslations(validate, transEn)
	}

	transId, found = uni.GetTranslator("id")
	if found {
		_ = id_translations.RegisterDefaultTranslations(validate, transId)
	}
}

// ValidateStruct validates a struct's fields using the validator/v10 library tags.
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

// TranslateValidationErrors translates ValidationErrors into a slice of human-readable messages based on locale ("en" or "id").
func TranslateValidationErrors(errs validator.ValidationErrors, locale string) []string {
	trans := transEn
	if locale == "id" {
		trans = transId
	}

	var messages []string
	for _, err := range errs {
		messages = append(messages, err.Translate(trans))
	}
	return messages
}
