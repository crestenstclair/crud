package handlers

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/crestenstclair/crud/internal/crud"
	"github.com/crestenstclair/crud/internal/user"
	"go.uber.org/zap"
)

func CreateUser(ctx context.Context, request events.APIGatewayProxyRequest, crud *crud.Crud) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(crud.Config.RequestTimeoutMS)*time.Millisecond)

	defer cancel()

	crud.Logger.Info(request.Body)

	var body map[string]string
	err := json.Unmarshal([]byte(request.Body), &body)
  if err != nil {
    crud.Logger.Error("Failed to create user", zap.Error(err))
		return makeResponse(map[string]string{
			"error": "An internal error occured",
		}, 500), nil
  }
	usr, err := user.New(
		body["firstName"],
		body["lastName"],
		body["email"],
		body["DOB"],
	)
	if err != nil {
		crud.Logger.Error("Invalid user provided", zap.Error(err))

		return makeResponse(map[string]string{
			"error": err.Error(),
		}, 400), nil
	}

	_, err = crud.Repo.CreateUser(ctx, *usr)

	switch err.(type) {
	case nil:
		return makeResponse(usr, 200), nil
	case *dynamodb.ConditionalCheckFailedException:
		crud.Logger.Error("Failed to create user, duplicate user provided", zap.Error(err))
		return makeResponse(map[string]string{
			"error": "Email already in use",
		}, 400), nil
	default:
		crud.Logger.Error("Failed to create user", zap.Error(err))
		return makeResponse(map[string]string{
			"error": "An internal error occured",
		}, 500), nil
	}
}
