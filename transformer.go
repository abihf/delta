package delta

import (
	"context"
	"net/http"
)

type Transformer interface {
	ToReq(context.Context, []byte) (*http.Request, error)
	FromRes(context.Context, *ResponseWriter) ([]byte, error)
}

func WithTransformer(t Transformer) Options {
	return func(c *config) {
		c.transformer = t
	}
}
