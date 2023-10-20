package handlers_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/crestenstclair/crud/internal/config"
	"github.com/crestenstclair/crud/internal/crud"
	"github.com/crestenstclair/crud/internal/handlers"
	"github.com/crestenstclair/crud/internal/repo/mocks"
	"github.com/crestenstclair/crud/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap/zaptest"
)

func TestCreateUser(t *testing.T) {
	t.Run("Returns 200 and user when create successful", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(getUserMap()),
		}, &testCrud)
		assert.NoError(t, err)
		result := &user.User{}

		err = json.Unmarshal([]byte(res.Body), &result)
		assert.NoError(t, err)

		assert.Equal(t, testUser.FirstName, result.FirstName)
		assert.Equal(t, testUser.LastName, result.LastName)
		assert.Equal(t, testUser.Email, result.Email)
		assert.Equal(t, testUser.DOB, result.DOB)
	})
	t.Run("errors when firstname is not provided", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		userMap := getUserMap()

		userMap["firstName"] = ""

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("errors when lastname is not provided", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		userMap := getUserMap()

		userMap["lastName"] = ""

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("errors when email is not provided", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		userMap := getUserMap()

		userMap["email"] = ""

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("errors when email is invalid", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		userMap := getUserMap()

		userMap["email"] = "not a vaild email"

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("errors when dob is not provided", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		userMap := getUserMap()

		userMap["DOB"] = ""

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("errors when dob is invalid", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		testUser := makeTestUser()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(&testUser, nil)

		userMap := getUserMap()

		userMap["DOB"] = "not a vaild DOB"

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
	t.Run("returns 500 when random error returns", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()

		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, errors.New("something went wrong"))

		userMap := getUserMap()

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 500, res.StatusCode)
	})
	t.Run("returns 400 when conditional check fails", func(t *testing.T) {
		mockRepo := mocks.Repo{}
		testCrud := crud.Crud{
			Repo:   &mockRepo,
			Logger: zaptest.NewLogger(t),
			Config: &config.Config{},
		}

		ctx := context.Background()
		mockRepo.On("CreateUser", mock.Anything, mock.Anything).Return(nil, &dynamodb.ConditionalCheckFailedException{})

		userMap := getUserMap()

		res, err := handlers.CreateUser(ctx, events.APIGatewayProxyRequest{
			Body: toJsonEscapedString(userMap),
		}, &testCrud)
		assert.NoError(t, err)

		assert.Equal(t, 400, res.StatusCode)
	})
}
