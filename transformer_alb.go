package delta

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

type alb struct{}

func WithALB() Options {
	return WithTransformer(&alb{})
}

// FromRes implements Transformer
func (*alb) FromRes(_ context.Context, r *ResponseWriter) ([]byte, error) {
	res := &events.ALBTargetGroupResponse{
		StatusCode:        r.status,
		MultiValueHeaders: r.header,
		IsBase64Encoded:   r.encode,
		Body:              r.bodyString(),
	}
	return json.Marshal(res)
}

// ToReq implements Transformer
func (*alb) ToReq(ctx context.Context, payload []byte) (*http.Request, error) {
	var e events.ALBTargetGroupRequest
	json.Unmarshal(payload, &e)
	header := http.Header(e.MultiValueHeaders)
	host := header.Get("host")
	qs := url.Values(e.MultiValueQueryStringParameters)
	u := url.URL{
		Path:     e.Path,
		Host:     host,
		RawQuery: qs.Encode(),
		Scheme:   "https",
	}
	req := http.Request{
		RequestURI: u.RequestURI(),
		Method:     e.HTTPMethod,
		URL:        &u,

		// just hardcode it
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,

		// content
		Header: header,
		Body:   ioutil.NopCloser(strings.NewReader(e.Body)),

		// from header
		ContentLength:    int64(len(e.Body)),
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       header.Get("x-forwarded-for"),
	}
	return &req, nil
}
