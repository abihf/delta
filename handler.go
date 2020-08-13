package delta

import (
	"context"
	"net/http"
	"os"

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
	var handler interface{}
	switch os.Getenv("LAMBDA_PAYLOAD_FORMAT") {
	case "2.0":
		handler = CreateLambdaHandlerV2(globalConfig, h)
	default:
		handler = CreateLambdaHandler(globalConfig, h)
	}
	lambda.Start(handler)
}

// ServeOrStartLambda will start http server if it's not in lambda environment,
// otherwise it starts handling lambda
func ServeOrStartLambda(addr string, h http.Handler) error {
	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		Start(h)
		return nil
	}

	return http.ListenAndServe(addr, h)
}

// LambdaHandler func type for lambda.Start()
type LambdaHandler func(context.Context, *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)

// CreateLambdaHandler create lambda handler
func CreateLambdaHandler(conf *Configuration, h http.Handler) LambdaHandler {
	return func(ctx context.Context, e *events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		req, err := NewRequest(ctx, e)
		if err != nil {
			return NewErrorResponse(err), err
		}
		res := NewResponseWriter()
		SetBase64Encoding(res, conf != nil && conf.EncodeResponse)

		h.ServeHTTP(res, req)
		lambdaResponse := res.ToAPIGWProxyResponse()
		return lambdaResponse, nil
	}
}

// LambdaHandlerV2 func type for lambda.Start()
type LambdaHandlerV2 func(context.Context, *events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error)

// CreateLambdaHandlerV2 create lambda handler
func CreateLambdaHandlerV2(conf *Configuration, h http.Handler) LambdaHandlerV2 {
	return func(ctx context.Context, e *events.APIGatewayV2HTTPRequest) (*events.APIGatewayProxyResponse, error) {
		req, err := NewRequestV2(ctx, e)
		if err != nil {
			return NewErrorResponse(err), err
		}
		res := NewResponseWriter()
		SetBase64Encoding(res, conf != nil && conf.EncodeResponse)

		h.ServeHTTP(res, req)
		lambdaResponse := res.ToAPIGWProxyResponse()
		return lambdaResponse, nil
	}
}
