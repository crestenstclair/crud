package user

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/crestenstclair/crud/internal/validator"
	"github.com/google/uuid"
)

type User struct {
  ID           string `validate:"uuid"`
	FirstName    string `validate:"required"`
	LastName     string `validate:"required"`
	Email        string `validate:"required,email"`
	DOB          string `validate:"required,RFC3339Date"`
	CreatedAt    string `validate:"RFC3339Date"`
	LastModified string `validate:"RFC3339Date"`
}

func Parse(jsonString string) (*User, error) {
	result := &User{}
	err := json.Unmarshal([]byte(jsonString), &result)
	if err != nil {
		return nil, err
	}

	err = validator.GetValidator().Struct(result)
	if err != nil {
		return nil, fmt.Errorf("User validation failed. %s", err)
	}

	return result, nil
}

func New(FirstName string, LastName string, Email string, DOB string) (*User, error) {
	result := &User{
		FirstName:    FirstName,
		LastName:     LastName,
		Email:        Email,
		DOB:          DOB,
	}

  result.LastModified = time.Now().Format(time.RFC3339)
  result.CreatedAt = time.Now().Format(time.RFC3339)
  result.ID = uuid.NewString()

	err := validator.GetValidator().Struct(result)
	if err != nil {
		return nil, fmt.Errorf("User validation failed. %s", err)
	}

	return result, nil
}
