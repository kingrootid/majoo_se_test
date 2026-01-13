package utils

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var trans ut.Translator

// SetupValidatorWithTranslations sets up validator with English translations
func SetupValidatorWithTranslations() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// Register translator
		enLoc := en.New()
		uni := ut.New(enLoc, enLoc)
		trans, _ = uni.GetTranslator("en")
		
		// Register English translations
		_ = en_translations.RegisterDefaultTranslations(v, trans)
	}
}

// GetValidator returns the validator instance
func GetValidator() *validator.Validate {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		return v
	}
	return nil
}

// GetTranslator returns the translator instance
func GetTranslator() ut.Translator {
	return trans
}