package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/crestenstclair/crud/internal/config"
	"github.com/crestenstclair/crud/internal/crud"
	"github.com/crestenstclair/crud/internal/handlers"
	"github.com/crestenstclair/crud/internal/repo/mocks"
	"github.com/crestenstclair/crud/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func makeTestUser() user.User {
	testTime := time.Now().Format(time.RFC3339)
	return user.User{
		ID:           "id",
		FirstName:    "firstName",
		LastName:     "lastName",
		Email:        "example@example.com",
		DOB:          testTime,
		CreatedAt:    testTime,
		LastModified: testTime,
	}
}

func TestGetUser(t *testing.T) {
	t.Run("Returns 200 and user when fetch successful", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		res, err := handlers.GetUser(ctx, events.APIGatewayProxyRequest{}, &testCrud)
		assert.NoError(t, err)
		result := &user.User{}

		err = json.Unmarshal([]byte(res.Body), &result)
		assert.NoError(t, err)

		assert.Equal(t, testUser.FirstName, result.FirstName)
		assert.Equal(t, testUser.LastName, result.LastName)
		assert.Equal(t, testUser.Email, result.Email)
		assert.Equal(t, testUser.DOB.Format(time.RFC3339), result.DOB.Format(time.RFC3339))
		assert.Equal(t, testUser.CreatedAt.Format(time.RFC3339), result.CreatedAt.Format(time.RFC3339))
		assert.Equal(t, testUser.LastModified.Format(time.RFC3339), result.LastModified.Format(time.RFC3339))
	})
	t.Run("Returns 500 when an internal server error occurs", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(nil, errors.New("TestError"))

		res, err := handlers.GetUser(ctx, events.APIGatewayProxyRequest{}, &testCrud)

		assert.NoError(t, err)

		assert.Equal(t, 500, res.StatusCode)
	})
	t.Run("Returns 404 when user not found", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()

		mockRepo.On("GetUser", mock.Anything, mock.Anything).Return(nil, nil)

		res, err := handlers.GetUser(ctx, events.APIGatewayProxyRequest{}, &testCrud)

		assert.NoError(t, err)

		assert.Equal(t, 404, res.StatusCode)
	})
}
