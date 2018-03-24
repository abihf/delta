package delta

import (
	"bytes"
	"encoding/base64"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

type responseWriter struct {
	header http.Header
	buffer bytes.Buffer
	status int
}

func newResponseWriter() *responseWriter {
	res := &responseWriter{
		header: make(http.Header),
		status: 200,
	}

	// set default content-type
	res.header.Set("content-type", "application/json")

	return res
}

func (r *responseWriter) Header() http.Header {
	return r.header
}

func (r *responseWriter) Write(body []byte) (int, error) {
	return r.buffer.Write(body)
}

func (r *responseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *responseWriter) toLambdaResponse(encode bool) (events.APIGatewayProxyResponse, error) {
	var body string
	if encode {
		body = base64.StdEncoding.EncodeToString(r.buffer.Bytes())
	} else {
		body = r.buffer.String()
	}
	r.header.Set("content-length", strconv.Itoa(r.buffer.Len()))
	return events.APIGatewayProxyResponse{
		StatusCode:      r.status,
		Headers:         convertToLambdaHeader(r.header),
		Body:            body,
		IsBase64Encoded: encode,
	}, nil
}

func newErrorResponse(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:      500,
		Headers:         map[string]string{},
		Body:            err.Error(),
		IsBase64Encoded: false,
	}
}
