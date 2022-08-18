package delta

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

const contextKey = "lambda-event"

// GetLambdaEvent from context
func GetLambdaEvent(ctx context.Context) (*events.APIGatewayProxyRequest, error) {
	if v := ctx.Value(contextKey); v != nil {
		if event, ok := v.(*events.APIGatewayProxyRequest); ok {
			return event, nil
		}
	}
	return nil, errors.New("GetLambdaEvent: invalid context")
}

// GetLambdaEventV2 from context
func GetLambdaEventV2(ctx context.Context) (*events.APIGatewayV2HTTPRequest, error) {
	if v := ctx.Value(contextKey); v != nil {
		if event, ok := v.(*events.APIGatewayV2HTTPRequest); ok {
			return event, nil
		}
	}
	return nil, errors.New("GetLambdaEvent: invalid context")
}

func withLambdaEvent(ctx context.Context, event interface{}) context.Context {
	return context.WithValue(ctx, contextKey, event)
}
