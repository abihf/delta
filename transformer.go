package delta

import (
	"context"
	"net/http"
)

type Transformer interface {
	Request(context.Context, []byte) (*http.Request, error)
	Response(context.Context, *ResponseWriter) ([]byte, error)
}

func WithTransformer(t Transformer) Options {
	return func(c *config) {
		c.transformer = t
	}
}

func convertFromHttpHeader(header http.Header) map[string]string {
	res := make(map[string]string, len(header))
	for k, v := range header {
		res[k] = v[0]
	}
	return res
}

func convertToHttpHeader(header map[string]string) http.Header {
	res := make(http.Header, len(header))
	for k, v := range header {
		res.Set(k, v)
	}
	return res
}
