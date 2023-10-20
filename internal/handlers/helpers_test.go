package handlers_test

import (
	"bytes"
	"encoding/json"

	"github.com/crestenstclair/crud/internal/user"
	"github.com/google/uuid"
)

func getUserMap() map[string]string {
	return map[string]string{
		"firstName": "firstName",
		"lastName":  "lastName",
		"email":     "example@example.com",
		"DOB":       "1979-12-09T00:00:00Z",
	}
}

func toUserMap(u *user.User) map[string]string {
	return map[string]string{
		"firstName": u.FirstName,
		"lastName":  u.LastName,
		"email":     u.Email,
		"DOB":       u.DOB,
    "ID": u.ID,
	}
}

func toJsonEscapedString(target map[string]string) string {
	result, _ := json.Marshal(target)

	var buf bytes.Buffer
	json.HTMLEscape(&buf, result)

	return buf.String()
}

func makeTestUser() user.User {
	testTime := "1979-12-09T00:00:00Z"
	return user.User{
		ID:           uuid.NewString(),
		FirstName:    "firstName",
		LastName:     "lastName",
		Email:        "example@example.com",
		DOB:          testTime,
		CreatedAt:    testTime,
		LastModified: testTime,
	}
}
