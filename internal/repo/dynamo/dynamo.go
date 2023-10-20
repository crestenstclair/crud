package dynamo

import (
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

type DynamoRepo struct {
	client    dynamodbiface.DynamoDBAPI
	tableName string
}

type UniqueConstraintViolation struct {
	Message string
}

func (u UniqueConstraintViolation) Error() string {
	return u.Message
}

func New(tableName string, db dynamodbiface.DynamoDBAPI) (*DynamoRepo, error) {
	return &DynamoRepo{
		client:    db,
		tableName: tableName,
	}, nil
}
