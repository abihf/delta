package delta

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

type contextKeyType string

var contextKey contextKeyType = "lambda-event"

// GetLambdaEvent from context
func GetLambdaEvent(ctx context.Context) (*events.APIGatewayProxyRequest, error) {
	if v := ctx.Value(contextKey); v != nil {
		if event, ok := v.(*events.APIGatewayProxyRequest); ok {
			return event, nil
		}
	}
	return nil, errors.New("GetLambdaEvent: invalid context")
}

func attachLambdaEvent(ctx context.Context, event *events.APIGatewayProxyRequest) context.Context {
	return context.WithValue(ctx, contextKey, event)
}
