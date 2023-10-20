package dynamo_test

import (
	"context"
	"errors"
	"testing"

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
	DOB       = "1979-12-09T00:00:00Z"
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

func (m *DynamodbMockClient) UpdateItem(input *dynamodb.UpdateItemInput) (*dynamodb.UpdateItemOutput, error) {
	// Ferry arguments into mock call
	args := m.Called(input)

	arg0 := args.Get(0)
	var resultOne *dynamodb.UpdateItemOutput
	if arg0 != nil {
		resultOne = arg0.(*dynamodb.UpdateItemOutput)
	} else {
		resultOne = &dynamodb.UpdateItemOutput{
			Attributes: nil,
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

func (m *DynamodbMockClient) Query(input *dynamodb.QueryInput) (*dynamodb.QueryOutput, error) {
	args := m.Called(input)

	arg0 := args.Get(0)
	var resultOne *dynamodb.QueryOutput
	if arg0 != nil {
		resultOne = arg0.(*dynamodb.QueryOutput)
	} else {
		resultOne = &dynamodb.QueryOutput{}
	}

	return resultOne, args.Error(1)
}

func TestGetUser(t *testing.T) {
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
					S: aws.String(DOB),
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
		assert.Equal(t, "1979-12-09T00:00:00Z", result.DOB)
	})
}

func TestGetUserByEmail(t *testing.T) {
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)

		client.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{}, errors.New("test error"))

		ctx := context.Background()
		_, err := repo.GetUserByEmail(ctx, userID)

		assert.Error(t, err)
	})

	t.Run("Returns nil, no error when user is not found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{},
		}, nil)
		ctx := context.Background()
		res, err := repo.GetUserByEmail(ctx, userID)

		assert.NoError(t, err)
		assert.Nil(t, res)
	})

	t.Run("Mashalls properties as expected when User is found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("Query", &dynamodb.QueryInput{
			ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
				":email": {
					S: aws.String(email),
				},
			},
			KeyConditionExpression: aws.String("Email = :email"),
			IndexName:              aws.String("email"),
			TableName:              aws.String("tableName"),
		}).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{{
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
					S: aws.String(DOB),
				},
			}},
		}, nil)
		ctx := context.Background()
		result, err := repo.GetUserByEmail(ctx, email)

		assert.NoError(t, err)

		assert.Equal(t, "userID", result.ID)
		assert.Equal(t, "firstName", result.FirstName)
		assert.Equal(t, "lastName", result.LastName)
		assert.Equal(t, "example@example.com", result.Email)
		assert.Equal(t, "1979-12-09T00:00:00Z", result.DOB)
	})
}

func TestCreateUser(t *testing.T) {
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{}, nil)
		client.On("PutItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		_, err := repo.CreateUser(ctx, user.User{})

		assert.Error(t, err)
	})

	t.Run("Properly marshalls passed in user into attribute struct", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{}, nil)
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
					S: aws.String(DOB),
				},
				"CreatedAt": {
					S: aws.String(DOB),
				},
				"LastModified": {
					S: aws.String(DOB),
				},
			},
			TableName: aws.String("tableName"),
		}).Return(nil, nil)
		ctx := context.Background()
		_, err := repo.CreateUser(ctx, user.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			DOB:          DOB,
			CreatedAt:    DOB,
			LastModified: DOB,
		})

		assert.NoError(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	t.Run("Returns error when error occurs", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)
		client.On("Query", mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{
					"LastModified": {
						S: aws.String("1234"),
					}},
				{
					"CreatedAt": {
						S: aws.String("1234"),
					}},
			},
		}, nil)
		client.On("UpdateItem", mock.Anything).Return(nil, errors.New("test error"))
		ctx := context.Background()
		_, err := repo.UpdateUser(ctx, user.User{})

		assert.Error(t, err)
	})

	t.Run("Properly marshalls passed in user into attribute struct", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)

		client.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{{
				"ID": {
					S: aws.String(userID),
				},
				"Email": {
					S: aws.String(email),
				},
				"CreatedAt": {
					S: aws.String(DOB),
				},
				"LastModified": {
					S: aws.String(DOB),
				},
			}},
		}, nil)
		putMock := client.On("UpdateItem", mock.Anything).Return(nil, nil)
		ctx := context.Background()
		_, err := repo.UpdateUser(ctx, user.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			DOB:          DOB,
			CreatedAt:    DOB,
			LastModified: DOB,
		})

		assert.NoError(t, err)
		call := putMock.Parent.Calls[1]
		arg := call.Arguments[0].(*dynamodb.UpdateItemInput)

		assert.Equal(t, firstName, *arg.ExpressionAttributeValues[":FirstName"].S)
		assert.Equal(t, lastName, *arg.ExpressionAttributeValues[":LastName"].S)
		assert.Equal(t, email, *arg.ExpressionAttributeValues[":Email"].S)
		assert.Equal(t, DOB, *arg.ExpressionAttributeValues[":DOB"].S)
	})

	t.Run("Properly detects when a user's email is taken", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)

		client.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{{
				"ID": {
					S: aws.String("DifferentID"),
				},
				"Email": {
					S: aws.String(email),
				},
			}},
		}, nil)
		ctx := context.Background()
		_, err := repo.UpdateUser(ctx, user.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			DOB:          DOB,
			CreatedAt:    DOB,
			LastModified: DOB,
		})

		assert.Error(t, err, "User email update failed. Attempted to change email")
	})
	t.Run("Does not false positive email dupe when same user is found", func(t *testing.T) {
		client := &DynamodbMockClient{}
		repo, _ := dynamo.New("tableName", client)

		client.On("Query", mock.Anything, mock.Anything).Return(&dynamodb.QueryOutput{
			Items: []map[string]*dynamodb.AttributeValue{
				{
					"ID": {
						S: aws.String(userID),
					}},

				{
					"LastModified": {
						S: aws.String("1234"),
					}},
				{
					"CreatedAt": {
						S: aws.String("1234"),
					}},
			},
		}, nil)
		client.On("UpdateItem", mock.Anything).Return(nil, nil)
		ctx := context.Background()
		_, err := repo.UpdateUser(ctx, user.User{
			ID:           userID,
			FirstName:    firstName,
			LastName:     lastName,
			Email:        email,
			DOB:          DOB,
			CreatedAt:    DOB,
			LastModified: DOB,
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
			TableName:           aws.String("tableName"),
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
