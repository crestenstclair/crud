package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

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
		assert.Equal(t, testUser.DOB, result.DOB)
		assert.Equal(t, testUser.CreatedAt, result.CreatedAt)
		assert.Equal(t, testUser.LastModified, result.LastModified)
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
