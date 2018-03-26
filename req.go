package delta

import (
	"context"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type contextKey string

func createRequest(ctx context.Context, e events.APIGatewayProxyRequest) (*http.Request, error) {
	var body io.Reader = strings.NewReader(e.Body)
	if e.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	headers := convertToHTTPHeader(e.Headers)
	host := headers.Get("host")
	length, err := strconv.ParseInt(headers.Get("content-length"), 10, 64)
	if err != nil {
		length = -1
	}
	var qs []string
	for key, val := range e.QueryStringParameters {
		qs = append(qs, url.QueryEscape(key)+"="+url.QueryEscape(val))
	}
	req := &http.Request{
		Method: e.HTTPMethod,
		URL: &url.URL{
			Path:     e.Path,
			RawPath:  e.Path,
			RawQuery: strings.Join(qs, "&"),
		},

		// just hardcode it
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,

		// content
		Header: headers,
		Body:   ioutil.NopCloser(body),

		// from header
		ContentLength:    length,
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
	}

	withEvent := context.WithValue(ctx, contextKey("lambda-event"), &e)

	return req.WithContext(withEvent), nil
}

// GetLambdaEvent from context
func GetLambdaEvent(ctx context.Context) (*events.APIGatewayProxyRequest, error) {
	if v := ctx.Value(contextKey("lambda-event")); v != nil {
		if event, ok := v.(*events.APIGatewayProxyRequest); ok {
			return event, nil
		}
	}
	return nil, errors.New("GetLambdaEvent: invalid context")
}
