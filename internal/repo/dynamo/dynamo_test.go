package dynamo_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/crestenstclair/crud/internal/repo/dynamo"
	"github.com/crestenstclair/crud/internal/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Set up some shared testing variables to avoid magic strings
var (
	userID    = "userID"
	firstName = "firstName"
	lastName  = "lastName"
	email     = "example@example.com"
)

type DynamodbMockClient struct {
	dynamodbiface.DynamoDBAPI
	mock.Mock
}

func (m *DynamodbMockClient) GetItem(input *dynamodb.GetItemInput) (*dynamodb.GetItemOutput, error) {
	// Ferry arguments into mock call
	args := m.Called(input)

	arg0 := args.Get(0)
	var resultOne *dynamodb.GetItemOutput
	if arg0 != nil {
		resultOne = arg0.(*dynamodb.GetItemOutput)
	} else {
		resultOne = &dynamodb.GetItemOutput{
			Item: nil,
		}
	}

	return resultOne, args.Error(1)
}

func (m *DynamodbMockClient) PutItem(input *dynamodb.PutItemInput) (*dynamodb.PutItemOutput, error) {
	args := m.Called(input)

	return &dynamodb.PutItemOutput{}, args.Error(1)
}

func (m *DynamodbMockClient) DeleteItem(input *dynamodb.DeleteItemInput) (*dynamodb.DeleteItemOutput, error) {
	args := m.Called(input)

	return nil, args.Error(1)
}

func TestGetUser(t *testing.T) {
	DOB, _ := time.Parse(time.RFC3339, "1979-12-09T00:00:00Z")
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("GetItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		_, err := repo.GetUser(ctx, userID)

		assert.Error(t, err)
	})

	t.Run("Returns nil, no error when user is not found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("GetItem", mock.Anything).Return(nil, nil)
		ctx := context.Background()
		res, err := repo.GetUser(ctx, userID)

		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("Mashalls properties as expected when User is found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("GetItem", &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(userID),
				},
			},
			TableName: aws.String("tableName"),
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(userID),
				},
				"FirstName": {
					S: aws.String(firstName),
				},
				"LastName": {
					S: aws.String(lastName),
				},
				"Email": {
					S: aws.String(email),
				},
				"DOB": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
			},
		}, nil)
		ctx := context.Background()
		result, err := repo.GetUser(ctx, userID)

		assert.NoError(t, err)

		assert.Equal(t, "userID", result.ID)
		assert.Equal(t, "firstName", result.FirstName)
		assert.Equal(t, "lastName", result.LastName)
		assert.Equal(t, "example@example.com", result.Email)
		assert.Equal(t, "1979-12-09T00:00:00Z", result.DOB.Format(time.RFC3339))
	})
}

func TestCreateUser(t *testing.T) {
	DOB, _ := time.Parse(time.RFC3339, "1979-12-09T00:00:00Z")
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("PutItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		_, err := repo.CreateUser(ctx, user.User{})

		assert.Error(t, err)
	})

	t.Run("Properly marshalls passed in user into attribute struct", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("PutItem", &dynamodb.PutItemInput{
			Item: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(userID),
				},
				"FirstName": {
					S: aws.String(firstName),
				},
				"LastName": {
					S: aws.String(lastName),
				},
				"Email": {
					S: aws.String(email),
				},
				"DOB": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
				"CreatedAt": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
				"LastModified": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
			},
			ConditionExpression: aws.String("attribute_not_exists(Email)"),
			TableName:           aws.String("tableName"),
		}).Return(nil, nil)
		ctx := context.Background()
		_, err := repo.CreateUser(ctx, user.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			DOB:          DOB,
			CreatedAt:    &DOB,
			LastModified: &DOB,
		})

		assert.NoError(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	DOB, _ := time.Parse(time.RFC3339, "1979-12-09T00:00:00Z")
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("PutItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		_, err := repo.UpdateUser(ctx, user.User{})

		assert.Error(t, err)
	})

	t.Run("Properly marshalls passed in user into attribute struct", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("PutItem", &dynamodb.PutItemInput{
			ConditionExpression: aws.String("attribute_exists(ID)"),
			Item: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(userID),
				},
				"FirstName": {
					S: aws.String(firstName),
				},
				"LastName": {
					S: aws.String(lastName),
				},
				"Email": {
					S: aws.String(email),
				},
				"DOB": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
				"CreatedAt": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
				"LastModified": {
					S: aws.String(DOB.Format(time.RFC3339)),
				},
			},
			TableName: aws.String("tableName"),
		}).Return(nil, nil)
		ctx := context.Background()
		_, err := repo.UpdateUser(ctx, user.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			DOB:          DOB,
			CreatedAt:    &DOB,
			LastModified: &DOB,
		})

		assert.NoError(t, err)
	})
}

func TestDeleteUser(t *testing.T) {
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("DeleteItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		err := repo.DeleteUser(ctx, userID)

		assert.Error(t, err)
	})

	t.Run("Returns nil when no error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("DeleteItem", &dynamodb.DeleteItemInput{
			ConditionExpression: aws.String("attribute_exists(ID)"),
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(userID),
				},
			},
		}).Return(nil, nil)
		ctx := context.Background()
		err := repo.DeleteUser(ctx, userID)

		assert.NoError(t, err)
	})
}
