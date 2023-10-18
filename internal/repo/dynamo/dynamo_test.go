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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
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

func TestGetUser(t *testing.T) {
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("GetItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		_, err := repo.GetUser(ctx, "userID")

		assert.Error(t, err)
	})

	t.Run("Returns nil, no error when user is not found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("GetItem", mock.Anything).Return(nil, nil)
		ctx := context.Background()
		res, err := repo.GetUser(ctx, "userID")

		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("Mashalls properties as expected when User is found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("GetItem", &dynamodb.GetItemInput{
			Key: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String("userID"),
				},
			},
			TableName: aws.String("tableName"),
		}).Return(&dynamodb.GetItemOutput{
			Item: map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String("userID"),
				},
				"FirstName": {
					S: aws.String("firstName"),
				},
				"LastName": {
					S: aws.String("lastName"),
				},
				"Email": {
					S: aws.String("example@example.com"),
				},
				"DOB": {
					S: aws.String("1970-12-09T00:00:00Z"),
				},
			},
		}, nil)
		ctx := context.Background()
		result, err := repo.GetUser(ctx, "userID")

		assert.NoError(t, err)

		assert.Equal(t, "userID", result.ID)
		assert.Equal(t, "firstName", result.FirstName)
		assert.Equal(t, "lastName", result.LastName)
		assert.Equal(t, "example@example.com", result.Email)
		assert.Equal(t, "1970-12-09T00:00:00Z", result.DOB.Format(time.RFC3339))
	})
}
