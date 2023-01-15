package delta

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	json "github.com/json-iterator/go"
)

type LambdaURLTransformer struct{}

func WithLambdaURL() Options {
	return WithTransformer(LambdaURLTransformer{})
}

// Response implements Transformer
func (LambdaURLTransformer) Response(ctx context.Context, w *ResponseWriter) ([]byte, error) {
	res := &events.LambdaFunctionURLResponse{
		StatusCode:      w.status,
		Headers:         convertHttpHeader(w.header),
		IsBase64Encoded: w.encode,
		Body:            w.bodyString(),
	}
	return json.Marshal(res)
}

// Request implements Transformer
func (LambdaURLTransformer) Request(ctx context.Context, payload []byte) (*http.Request, error) {
	var e events.LambdaFunctionURLRequest
	json.Unmarshal(payload, &e)
	header := toHttpHeader(e.Headers)
	host := header.Get("host")
	var body []byte
	if e.IsBase64Encoded {
		var err error
		body, err = base64.StdEncoding.DecodeString(e.Body)
		if err != nil {
			return nil, err
		}
	} else {
		body = []byte(e.Body)
	}
	u := url.URL{
		Path:     e.RequestContext.HTTP.Path,
		RawPath:  e.RawPath,
		Host:     host,
		RawQuery: e.RawQueryString,
		Scheme:   "https",
	}
	req := http.Request{
		RequestURI: u.RequestURI(),
		Method:     e.RequestContext.HTTP.Method,
		URL:        &u,

		// just hardcode it
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,

		// content
		Header: header,
		Body:   io.NopCloser(bytes.NewReader(body)),

		// from header
		ContentLength:    int64(len(body)),
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       header.Get("x-forwarded-for"),
	}
	return withContextEvent(&req, ctx, &e), nil
}

func GetLambdaUrlEvent(ctx context.Context) (*events.LambdaFunctionURLRequest, error) {
	return getEvent[*events.LambdaFunctionURLRequest](ctx)
}
