package delta

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
)

type Handler struct {
	c *config
	h http.Handler
}

func NewHandler(h http.Handler, opts ...Options) *Handler {
	var c config
	for _, o := range opts {
		o(&c)
	}
	if c.transformer == nil {
		c.transformer = &ApiGatewayV1Transformer{}
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
	lh.h.ServeHTTP(res, req)

	// set content length
	if res.header.Get("content-length") == "" && res.header.Get("content-encoding") == "" {
		res.header.Set("content-length", strconv.Itoa(res.buffer.Len()))
	}

	b, err := lh.c.transformer.Response(ctx, res)
	if err != nil {
		return nil, fmt.Errorf("can not transform response: %w", err)
	}
	return b, nil
}

func (lh *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lh.h.ServeHTTP(w, r)
}
