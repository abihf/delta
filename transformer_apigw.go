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
	res := &events.APIGatewayProxyResponse{
		StatusCode:        r.status,
		MultiValueHeaders: r.header,
		IsBase64Encoded:   r.encode,
		Body:              r.bodyString(),
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
	req := &http.Request{
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

	return req, nil
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

	var body io.Reader = strings.NewReader(e.Body)
	if e.IsBase64Encoded {
		body = base64.NewDecoder(base64.StdEncoding, body)
	}

	header := apigwConvertHeader(e.Headers)
	host := header.Get("host")
	length, err := strconv.ParseInt(header.Get("content-length"), 10, 64)
	if err != nil {
		length = -1
	}
	var qs []string
	for key, val := range e.QueryStringParameters {
		qs = append(qs, url.QueryEscape(key)+"="+url.QueryEscape(val))
	}
	u := &url.URL{
		Scheme:   e.RequestContext.Protocol,
		Path:     e.Path,
		RawPath:  e.Path,
		RawQuery: strings.Join(qs, "&"),
	}
	req := &http.Request{
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
		ContentLength:    length,
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       e.RequestContext.Identity.SourceIP,
	}

	return req, nil
}

// apigwConvertHeader creates new Header from APIGWProxyHeader
func apigwConvertHeader(ph map[string]string) http.Header {
	header := make(http.Header)
	for name, value := range ph {
		header.Set(name, value)
	}
	return header
}
