package dynamo

import (
	"context"

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

func New(tableName string, db dynamodbiface.DynamoDBAPI) (*DynamoRepo, error) {
	// TODO: add test to ensure client is working
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

	dynamodbattribute.UnmarshalMap(response.Item, &result)

	return result, nil
}
