package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crestenstclair/crud/internal/user"
)

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

func (d DynamoRepo) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	response, err := d.client.Query(&dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":email": {
				S: aws.String(email),
			},
		},
		KeyConditionExpression: aws.String("Email = :email"),
		IndexName:              aws.String("email"),
		TableName:              &d.tableName,
	})
	
  if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, nil
	}

	var result *user.User

	err = dynamodbattribute.UnmarshalMap(response.Items[0], &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
