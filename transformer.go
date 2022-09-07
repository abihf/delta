package delta

import (
	"context"
	"net/http"
	"strings"
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

func convertHttpHeader(header http.Header) map[string]string {
	res := make(map[string]string, len(header))
	for k, v := range header {
		res[k] = strings.Join(v, ", ")
	}
	return res
}

func toHttpHeader(header map[string]string) http.Header {
	res := make(http.Header, len(header))
	for k, v := range header {
		res.Set(k, v)
	}
	return res
}
