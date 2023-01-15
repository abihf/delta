package delta

import (
	"context"
	"errors"
	"net/http"
)

type contextKeyType string
var contextKey = contextKeyType("lambda-event")

func getEvent[T any](ctx context.Context) (T, error) {
	if v := ctx.Value(contextKey); v != nil {
		if event, ok := v.(T); ok {
			return event, nil
		}
	}
	var zero T
	return zero, errors.New("can not get event: invalid context")
}

func withContextEvent(r *http.Request, ctx context.Context, event interface{}) *http.Request {
	return r.WithContext(context.WithValue(ctx, contextKey, event))
}
