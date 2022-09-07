package delta

import (
	"bytes"
	"context"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	json "github.com/json-iterator/go"
)

type lambdaUrl struct{}

func WithLambdaURL() Options {
	return WithTransformer(&lambdaUrl{})
}

// FromRes implements Transformer
func (*lambdaUrl) FromRes(ctx context.Context, w *ResponseWriter) ([]byte, error) {
	res := &events.LambdaFunctionURLResponse{
		StatusCode:      w.status,
		Headers:         convertHttpHeader(w.header),
		IsBase64Encoded: w.encode,
		Body:            w.bodyString(),
	}
	return json.Marshal(res)
}

// ToReq implements Transformer
func (*lambdaUrl) ToReq(ctx context.Context, payload []byte) (*http.Request, error) {
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
		Body:   ioutil.NopCloser(bytes.NewReader(body)),

		// from header
		ContentLength:    int64(len(body)),
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       header.Get("x-forwarded-for"),
	}
	return &req, nil
}
