package delta

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// ResponseWriter implements http.ResponseWriter used for buffering response
type ResponseWriter struct {
	header Header
	buffer bytes.Buffer
	status int
	encode bool
}

// NewResponseWriter creates new empty ResponseWriter
func NewResponseWriter() *ResponseWriter {
	res := &ResponseWriter{
		header: Header{make(http.Header)},
		status: 200,
	}

	// set default content-type
	res.header.Set("content-type", "application/json")

	return res
}

// Header returns http.Header. You can modify it to send response header
func (r *ResponseWriter) Header() http.Header {
	return r.header.Header
}

// Write appends chunk to response body
func (r *ResponseWriter) Write(chunk []byte) (int, error) {
	return r.buffer.Write(chunk)
}

// WriteHeader set status code of current request
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

// ToAPIGWProxyResponse convert it to events.APIGatewayProxyResponse
func (r *ResponseWriter) ToAPIGWProxyResponse() *events.APIGatewayProxyResponse {
	var body string
	if r.encode {
		body = base64.StdEncoding.EncodeToString(r.buffer.Bytes())
	} else {
		body = r.buffer.String()
	}
	r.header.Set("content-length", strconv.Itoa(r.buffer.Len()))
	return &events.APIGatewayProxyResponse{
		StatusCode:      r.status,
		Headers:         r.header.ToAPIGWProxyHeader(),
		Body:            body,
		IsBase64Encoded: r.encode,
	}
}

// NewErrorResponse create API Gateway Proxy Response contains error message
func NewErrorResponse(err error) *events.APIGatewayProxyResponse {
	return &events.APIGatewayProxyResponse{
		StatusCode:      500,
		Headers:         map[string]string{},
		Body:            err.Error(),
		IsBase64Encoded: false,
	}
}

// SetBase64Encoding overrides base64 encoding for this response
// see Configuration.SetBase64Encoding
func SetBase64Encoding(w http.ResponseWriter, enabled bool) error {
	if rw, ok := w.(*ResponseWriter); ok {
		rw.encode = enabled
		return nil
	}
	return errors.New("SetBase64Encoding: invalid ResponseWriter object")
}
