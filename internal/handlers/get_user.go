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

func GetUser(ctx context.Context, request events.APIGatewayProxyRequest, crud *crud.Crud) (events.APIGatewayProxyResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(crud.Config.RequestTimeoutMS)*time.Millisecond)

	defer cancel()

	id := request.PathParameters["id"]
	user, err := crud.Repo.GetUser(ctx, id)
	if err != nil {
		crud.Logger.Error("Failed to get user", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            "An internal error has occured",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil

	}

	if user == nil {
		return events.APIGatewayProxyResponse{
			StatusCode:      404,
			IsBase64Encoded: false,
			Body:            fmt.Sprintf("User not found. ID: %s", id),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil
	}

	body, err := json.Marshal(user)
	if err != nil {
		crud.Logger.Error("Failed to get user", zap.Error(err))
		return events.APIGatewayProxyResponse{
			StatusCode:      500,
			IsBase64Encoded: false,
			Body:            "An internal error has occured",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}, nil

	}

	var buf bytes.Buffer
	json.HTMLEscape(&buf, body)

	return events.APIGatewayProxyResponse{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            buf.String(),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
