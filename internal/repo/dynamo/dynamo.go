package dynamo

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/crestenstclair/crud/internal/user"
)

type DynamoRepo struct {
	client    dynamodbiface.DynamoDBAPI
	tableName string
}

type UniqueConstraintViolation struct {
	Message string
}

func (u UniqueConstraintViolation) Error() string {
	return fmt.Sprintf("%s", u.Message)
}

func New(tableName string, db dynamodbiface.DynamoDBAPI) (*DynamoRepo, error) {
	return &DynamoRepo{
		client:    db,
		tableName: tableName,
	}, nil
}

func (d DynamoRepo) GetUser(ctx context.Context, userID string) (*user.User, error) {
	response, err := d.client.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
		TableName: &d.tableName,
	})
	if err != nil {
		return nil, err
	}

	if response.Item == nil {
		// We don't return an error in this case b/c it is not
		// an error specific to querying DynamoDB.
		// That will be handled at a higher level.
		return nil, nil
	}

	var result *user.User

	err = dynamodbattribute.UnmarshalMap(response.Item, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d DynamoRepo) CreateUser(ctx context.Context, u user.User) (*user.User, error) {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}
	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		Item:                av,
		TableName:           &d.tableName,
		ConditionExpression: aws.String("attribute_not_exists(Email)"),
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (d DynamoRepo) UpdateUser(ctx context.Context, u user.User) (*user.User, error) {
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}

	existingUser, err := d.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil && existingUser.ID != u.ID {
		return nil, &UniqueConstraintViolation{
			Message: "User email update failed. Attempted to change email to existing users email.",
		}
	}

	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		Item:                av,
		TableName:           &d.tableName,
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (d DynamoRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	response, err := d.client.GetItem(&dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Email": {
				S: aws.String(email),
			},
		},
		TableName: &d.tableName,
	})
	if err != nil {
		return nil, err
	}

	if response.Item == nil {
		return nil, nil
	}

	var result *user.User

	err = dynamodbattribute.UnmarshalMap(response.Item, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d DynamoRepo) DeleteUser(ctx context.Context, userID string) error {
	_, err := d.client.DeleteItem(&dynamodb.DeleteItemInput{
		TableName:           &d.tableName,
		ConditionExpression: aws.String("attribute_exists(ID)"),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(userID),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}
