package dynamo

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/crestenstclair/crud/internal/user"
)

func (d DynamoRepo) UpdateUser(ctx context.Context, u user.User) (*user.User, error) {
	existingUser, err := d.GetUserByEmail(ctx, u.Email)
	if err != nil {
		return nil, err
	}

	if existingUser != nil && existingUser.ID != u.ID {
		return nil, &UniqueConstraintViolation{
			Message: "User email update failed. Attempted to change email to existing users email.",
		}
	}

	// Update last modified to right now
	u.LastModified = time.Now().Format(time.RFC3339)

	// Marshal into map to avoid a bunch of boilerplate
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, err
	}

	// Remove unwanted fields
	delete(av, "CreatedAt")
	delete(av, "ID")

	// Initialize update expression in order to ensure CreatedAt is preserved between updates
	updateExpression := "set CreatedAt = CreatedAt"
	expressionValues := map[string]*dynamodb.AttributeValue{}

	// Populate expression and value map with escaped "expressionvalue" keys
	for k, v := range av {
		key := ":" + k
		updateExpression += fmt.Sprintf(", %s = :%s", k, k)
		expressionValues[key] = v
	}

	response, err := d.client.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(u.ID),
			},
		},
		TableName:                 &d.tableName,
		ConditionExpression:       aws.String("attribute_exists(ID)"),
		UpdateExpression:          aws.String(updateExpression),
		ExpressionAttributeValues: expressionValues,
		ReturnValues:              aws.String("ALL_NEW"),
	})

	if err != nil {
		return nil, err
	}

	var result *user.User

	err = dynamodbattribute.UnmarshalMap(response.Attributes, &result)

	if err != nil {
		return nil, err
	}

	return result, nil
}
