package delta

import (
	"context"
	"encoding/base64"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	json "github.com/json-iterator/go"
)

type apigwCommon struct{}

// Response implements Transformer
func (apigwCommon) Response(ctx context.Context, r *ResponseWriter) ([]byte, error) {
	headers := make(map[string]string, len(r.header))
	for name, value := range r.header {
		headers[name] = value[0]
	}
	res := &events.APIGatewayProxyResponse{
		StatusCode:      r.status,
		Headers:         headers,
		IsBase64Encoded: r.encode,
		Body:            r.bodyString(),
	}
	return json.Marshal(res)
}

type ApiGatewayV2Transformer struct {
	apigwCommon
}

func WithAPIGatewayV2() Options {
	return WithTransformer(ApiGatewayV2Transformer{})
}

// Request implements Transformer
func (ApiGatewayV2Transformer) Request(ctx context.Context, payload []byte) (*http.Request, error) {
	var e events.APIGatewayV2HTTPRequest
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return nil, err
	}
	var body io.Reader = strings.NewReader(e.Body)
	if e.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	header := apigwConvertHeader(e.Headers)
	header["Cookie"] = e.Cookies
	host := e.RequestContext.DomainName
	length, _ := strconv.ParseInt(header.Get("content-length"), 10, 64)
	u := &url.URL{
		Path:     e.RequestContext.HTTP.Path,
		Host:     host,
		RawPath:  e.RawPath,
		RawQuery: e.RawQueryString,
		Scheme:   "https",
	}
	req := http.Request{
		RequestURI: u.RequestURI(),
		Method:     e.RequestContext.HTTP.Method,
		URL:        u,

		// just hardcode it
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,

		// content
		Header: header,
		Body:   io.NopCloser(body),

		// from header
		ContentLength:    length,
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       e.RequestContext.HTTP.SourceIP,
	}

	return withContextEvent(&req, ctx, &e), nil
}

func GetApiGatewayV2Event(ctx context.Context) (*events.APIGatewayV2HTTPRequest, error) {
	return getEvent[*events.APIGatewayV2HTTPRequest](ctx)
}

type ApiGatewayV1Transformer struct {
	apigwCommon
}

func WithAPIGatewayV1() Options {
	return WithTransformer(ApiGatewayV1Transformer{})
}

// Request implements Transformer
func (ApiGatewayV1Transformer) Request(ctx context.Context, payload []byte) (*http.Request, error) {
	var e events.APIGatewayProxyRequest
	err := json.Unmarshal(payload, &e)
	if err != nil {
		return nil, err
	}

	var body io.Reader = strings.NewReader(e.Body)
	if e.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	header := http.Header(e.MultiValueHeaders)
	host := header.Get("host")
	u := &url.URL{
		Scheme:   e.RequestContext.Protocol,
		Path:     e.Path,
		RawQuery: url.Values(e.MultiValueQueryStringParameters).Encode(),
	}

	req := http.Request{
		Method:     e.HTTPMethod,
		URL:        u,
		RequestURI: u.RequestURI(),

		// just hardcode it
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,

		// content
		Header: header,
		Body:   io.NopCloser(body),

		// from header
		ContentLength:    int64(len(e.Body)),
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       e.RequestContext.Identity.SourceIP,
	}

	return withContextEvent(&req, ctx, &e), nil
}

func GetApiGatewayV1Event(ctx context.Context) (*events.APIGatewayProxyRequest, error) {
	return getEvent[*events.APIGatewayProxyRequest](ctx)
}

// apigwConvertHeader creates new Header from APIGWProxyHeader
func apigwConvertHeader(ph map[string]string) http.Header {
	header := make(http.Header)
	for name, value := range ph {
		header.Set(name, value)
	}
	return header
}
