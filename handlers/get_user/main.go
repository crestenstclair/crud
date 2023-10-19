package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/crestenstclair/crud/internal/crud"
	"github.com/crestenstclair/crud/internal/handlers"
)

var inst *crud.Crud

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return handlers.GetUser(ctx, request, inst)
}

func main() {
	lambda.Start(Handler)
}

func init() {
	tmp, err := crud.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	inst = tmp
}
