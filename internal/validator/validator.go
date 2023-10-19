package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func GetValidator() *validator.Validate {
	return validate
}

func IsRFC3339Date(fl validator.FieldLevel) bool {
	target := fl.Field().String()
	// Handle case where field is note required
	if target == "" {
		return true
	}
	_, err := time.Parse(time.RFC3339, target)
	return err == nil
}

func init() {
	// Singleton validation is the recommended way to do validation according to validator
	// https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Singleton
	validate = validator.New(validator.WithRequiredStructEnabled())

	_ = validate.RegisterValidation("RFC3339Date", IsRFC3339Date)
}
