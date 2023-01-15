package delta

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	json "github.com/json-iterator/go"
)

type AlbTransformer struct{}

func WithALB() Options {
	return WithTransformer(AlbTransformer{})
}

// Response implements Transformer
func (AlbTransformer) Response(_ context.Context, r *ResponseWriter) ([]byte, error) {
	res := &events.ALBTargetGroupResponse{
		StatusCode:        r.status,
		MultiValueHeaders: r.header,
		IsBase64Encoded:   r.encode,
		Body:              r.bodyString(),
	}
	return json.Marshal(res)
}

// Request implements Transformer
func (AlbTransformer) Request(ctx context.Context, payload []byte) (*http.Request, error) {
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
		Body:   io.NopCloser(strings.NewReader(e.Body)),

		// from header
		ContentLength:    int64(len(e.Body)),
		TransferEncoding: []string{},
		Close:            true,
		Host:             host,
		RemoteAddr:       header.Get("x-forwarded-for"),
	}
	return withContextEvent(&req, ctx, &e), nil
}

func GetAlbEvent(ctx context.Context) (*events.ALBTargetGroupRequest, error) {
	return getEvent[*events.ALBTargetGroupRequest](ctx)
}
