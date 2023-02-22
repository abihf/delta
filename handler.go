package delta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

type Handler struct {
	c *config
	h http.Handler
}

// make sure it implements lambda's & http's handler
var _ lambda.Handler = &Handler{}
var _ http.Handler = &Handler{}

func NewHandler(h http.Handler, opts ...Options) *Handler {
	var c config
	for _, o := range opts {
		o(&c)
	}
	if c.transformer == nil {
		c.transformer = GetDefaultTransformer()
	}
	return &Handler{h: h, c: &c}
}

// Invoke implements lambda.Handler
func (lh *Handler) Invoke(ctx context.Context, payload []byte) ([]byte, error) {
	req, err := lh.c.transformer.Request(ctx, payload)
	if err != nil {
		return nil, fmt.Errorf("can not transform request payload: %w", err)
	}

	res := NewResponseWriter()
	res.encode = lh.c.encodeResponse
	lh.h.ServeHTTP(res, req)

	b, err := lh.c.transformer.Response(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("can not transform response: %w", err)
	}
	return b, nil
}

func (lh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lh.h.ServeHTTP(w, r)
}
