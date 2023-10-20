package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crestenstclair/crud/internal/user"
)

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
