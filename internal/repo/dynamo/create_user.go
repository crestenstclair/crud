package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crestenstclair/crud/internal/user"
)

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
