package handlers

import (
	"context"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/crestenstclair/crud/internal/crud"
	"go.uber.org/zap"
)

func DeleteUser(ctx context.Context, request events.APIGatewayProxyRequest, crud *crud.Crud) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(crud.Config.RequestTimeoutMS)*time.Millisecond)
	defer cancel()

	id := request.PathParameters["id"]
	err := crud.Repo.DeleteUser(ctx, id)

	switch err.(type) {
	case nil:
		return makeResponse(map[string]string{
			id: id,
		}, 200), nil
	case *dynamodb.ConditionalCheckFailedException:
		crud.Logger.Error("Failed to delete user, ID not found", zap.Error(err))
		return makeResponse(map[string]string{
			"error": "Failed to delete user, ID not found",
			id:      id,
		}, 404), nil
	default:
		crud.Logger.Error("Failed to delete user", zap.Error(err))
		return makeResponse(map[string]string{
			"error": "An internal error occured",
		}, 500), nil
	}
}
