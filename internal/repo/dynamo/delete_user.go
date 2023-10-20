package dynamo

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

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
