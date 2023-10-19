package validator_test

import (
	"testing"
	"time"

	"github.com/crestenstclair/crud/internal/validator"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	DOB            string `validate:"required,RFC3339Date"`
	NonRequiredDOB string `validate:"RFC3339Date"`
}

func TestValidator(t *testing.T) {
	t.Run("Does not do false positive on valid rfc3339 date", func(t *testing.T) {
		testTime := time.Now().Format(time.RFC3339)
		target := &testStruct{
			DOB: testTime,
		}

		err := validator.GetValidator().Struct(target)

		assert.NoError(t, err)
	})
	t.Run("Does not do false positive on non-required rfc3339 date", func(t *testing.T) {
		testTime := time.Now().Format(time.RFC3339)
		target := &testStruct{
			DOB: testTime,
		}

		err := validator.GetValidator().Struct(target)

		assert.NoError(t, err)
	})
}
