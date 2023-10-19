package handlers

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
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
