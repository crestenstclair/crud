package user

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type User struct {
	FirstName string    `validate:"required"`
	LastName  string    `validate:"required"`
	Email     string    `validate:"required,email"`
	DOB       time.Time `validate:"required"`
}

func New(FirstName string, LastName string, Email string, DOB string) (*User, error) {
	dob, err := time.Parse(time.RFC3339, DOB)
	if err != nil {
		return nil, fmt.Errorf("Parsing DOB failed. DOB must be in RFC3339 format. %s", err)
	}

	result := &User{
		FirstName: FirstName,
		LastName:  LastName,
		Email:     Email,
		DOB:       dob,
	}

	err = validate.Struct(result)

	if err != nil {
		return nil, fmt.Errorf("User validation failed. %s", err)
	}

	return result, nil
}

func init() {
	// Singleton validation is the recommended way to do validation according to validator
	// https://pkg.go.dev/github.com/go-playground/validator/v10#hdr-Singleton
	validate = validator.New(validator.WithRequiredStructEnabled())
}
