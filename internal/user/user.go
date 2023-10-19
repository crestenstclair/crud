package user

import (
	"fmt"

	"github.com/crestenstclair/crud/internal/validator"
)

type User struct {
	ID           string
	FirstName    string `validate:"required"`
	LastName     string `validate:"required"`
	Email        string `validate:"required,email"`
	DOB          string `validate:"required,RFC3339Date"`
	CreatedAt    string `validate:"RFC3339Date"`
	LastModified string `validate:"RFC3339Date"`
}

func New(FirstName string, LastName string, Email string, DOB string, createdAt string, lastModified string) (*User, error) {
	result := &User{
		FirstName:    FirstName,
		LastName:     LastName,
		Email:        Email,
		DOB:          DOB,
		CreatedAt:    createdAt,
		LastModified: lastModified,
	}

	err := validator.GetValidator().Struct(result)
	if err != nil {
		return nil, fmt.Errorf("User validation failed. %s", err)
	}

	return result, nil
}
