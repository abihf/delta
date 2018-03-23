package delta

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// Start lambda server
//
// Example:
// mux := http.NewServeMux()
// mux.Handle("/", handeHandler)
// delta.Start(nil, mux)
func Start(h http.Handler) {
	lambda.Start(CreateLambdaHandler(globalConfig, h))
}

// LambdaHandler func type for lambda.Start()
type LambdaHandler func(context.Context, events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)

// CreateLambdaHandler create lambda handler
func CreateLambdaHandler(conf *Configuration, h http.Handler) LambdaHandler {
	return func(ctx context.Context, e events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		req, err := createRequest(ctx, e)
		if err != nil {
			return newErrorResponse(err), err
		}
		res := newResponseWriter()
		h.ServeHTTP(res, req)
		return res.toLambdaResponse(conf != nil && conf.EncodeResponse)
	}
}
