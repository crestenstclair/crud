package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/crestenstclair/crud/internal/crud"
	"go.uber.org/zap"
)

func makeResponse[T any](respBody T, statusCode int) events.APIGatewayProxyResponse {
	var buf bytes.Buffer

	body, _ := json.Marshal(respBody)

	json.HTMLEscape(&buf, body)
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}
}

func GetUser(ctx context.Context, request events.APIGatewayProxyRequest, crud *crud.Crud) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(crud.Config.RequestTimeoutMS)*time.Millisecond)

	defer cancel()

	id := request.PathParameters["id"]
	user, err := crud.Repo.GetUser(ctx, id)
	if err != nil {
		crud.Logger.Error("Failed to get user", zap.Error(err))

		return makeResponse(map[string]string{
			"error": "An internal error occured",
		}, 500), nil
	}

	if user == nil {
		crud.Logger.Error("User not found", zap.String("id", id))
		return makeResponse(map[string]string{
			"error": fmt.Sprintf("User not found. ID: %s", id),
		}, 404), nil
	}

	if err != nil {
		crud.Logger.Error("Failed to parse user after fetching from database", zap.Error(err))
		return makeResponse(map[string]string{
			"error": "An internal error occured",
		}, 500), nil
	}

	return makeResponse(user, 200), nil
}
