package user_test

import (
	"testing"

	"github.com/crestenstclair/crud/internal/user"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	exampleFirstName := "Fred"
	exampleLastName := "Flintstone"
	exampleEmail := "fredflintstone@example.com"
	exampleDOB := "1970-12-09T00:00:00Z"

	t.Run("Errors when provided an invalid date", func(t *testing.T) {
		_, err := user.New(
			exampleFirstName,
			exampleLastName,
			exampleEmail,
			"INVALID",
			exampleDOB,
			exampleDOB,
		)
		assert.ErrorContains(t, err, "User validation failed. Key: 'User.DOB'")
	})
	t.Run("Errors when DOB is not provided", func(t *testing.T) {
		_, err := user.New(
			exampleFirstName,
			exampleLastName,
			exampleEmail,
			"",
			exampleDOB,
			exampleDOB,
		)
		assert.ErrorContains(t, err, "User validation failed. Key: 'User.DOB'")
	})
	t.Run("Errors when not provided required FirstName", func(t *testing.T) {
		_, err := user.New(
			"",
			exampleLastName,
			exampleEmail,
			exampleDOB,
			exampleDOB,
			exampleDOB,
		)

		assert.ErrorContains(t, err, "User validation failed. Key: 'User.FirstName'")
	})
	t.Run("Errors when not provided required LastName", func(t *testing.T) {
		_, err := user.New(
			exampleFirstName,
			"",
			exampleEmail,
			exampleDOB,
			exampleDOB,
			exampleDOB,
		)

		assert.ErrorContains(t, err, "User validation failed. Key: 'User.LastName'")
	})
	t.Run("Errors when not provided required Email", func(t *testing.T) {
		_, err := user.New(
			exampleFirstName,
			exampleLastName,
			"",
			exampleDOB,
			exampleDOB,
			exampleDOB,
		)

		assert.ErrorContains(t, err, "User validation failed. Key: 'User.Email'")
	})
	t.Run("Errors when Email is invalid", func(t *testing.T) {
		_, err := user.New(
			exampleFirstName,
			exampleLastName,
			"notAnEmail",
			exampleDOB,
			exampleDOB,
			exampleDOB,
		)

		assert.ErrorContains(t, err, "User validation failed. Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag")
	})
	t.Run("Successfully parses datetime", func(t *testing.T) {
		result, err := user.New(
			exampleFirstName,
			exampleLastName,
			exampleEmail,

			exampleDOB,
			exampleDOB,
			exampleDOB,
		)

		assert.NoError(t, err)

		assert.Equal(t, exampleFirstName, result.FirstName)
		assert.Equal(t, exampleLastName, result.LastName)
		assert.Equal(t, exampleEmail, result.Email)
		assert.Regexp(t, exampleDOB, result.DOB)
	})
	t.Run("Does not error when user is valid", func(t *testing.T) {})
}
