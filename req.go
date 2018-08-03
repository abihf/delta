package delta

import (
	"context"
	"encoding/base64"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// NewRequest create http.Request based on API Gateway Proxy request
func NewRequest(ctx context.Context, e *events.APIGatewayProxyRequest) (*http.Request, error) {
	var body io.Reader = strings.NewReader(e.Body)
	if e.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	header := HeaderFromAPIGWProxyHeader(e.Headers)
	host := header.Get("host")
	length, err := strconv.ParseInt(header.Get("content-length"), 10, 64)
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
		Header: header.Header,
		Body:   ioutil.NopCloser(body),

		// from header
		ContentLength:    length,
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
	}

	withEvent := attachLambdaEvent(ctx, e)
	return req.WithContext(withEvent), nil
}
