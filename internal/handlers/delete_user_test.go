package handlers_test

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/crestenstclair/crud/internal/config"
	"github.com/crestenstclair/crud/internal/crud"
	"github.com/crestenstclair/crud/internal/handlers"
	"github.com/crestenstclair/crud/internal/repo/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestDeleteUser(t *testing.T) {
	t.Run("Returns 200 when delete successful", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()

		mockRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(nil)

		res, err := handlers.DeleteUser(ctx, events.APIGatewayProxyRequest{}, &testCrud)
		assert.NoError(t, err)

    assert.Equal(t, 200, res.StatusCode)

	})
	t.Run("Returns 500 when an internal server error occurs", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()

		mockRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(errors.New("TestError"))

		res, err := handlers.DeleteUser(ctx, events.APIGatewayProxyRequest{}, &testCrud)

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

		mockRepo.On("DeleteUser", mock.Anything, mock.Anything).Return(&dynamodb.ConditionalCheckFailedException{})

		res, err := handlers.DeleteUser(ctx, events.APIGatewayProxyRequest{}, &testCrud)

		assert.NoError(t, err)

		assert.Equal(t, 404, res.StatusCode)
	})
}
