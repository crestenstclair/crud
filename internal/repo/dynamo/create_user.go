package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crestenstclair/crud/internal/user"
)

func (d DynamoRepo) CreateUser(ctx context.Context, u user.User) (*user.User, error) {
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
			Message: "User creation failed. Email already in use by existing user.",
		}
	}

	_, err = d.client.PutItem(&dynamodb.PutItemInput{
		Item:      av,
		TableName: &d.tableName,
	})

	if err != nil {
		return nil, err
	}

	return &u, nil
}
