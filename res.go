package delta

import (
	"bytes"
	"encoding/base64"
	"errors"
	"net/http"
)

// ResponseWriter implements http.ResponseWriter used for buffering response
type ResponseWriter struct {
	header http.Header
	buffer bytes.Buffer
	status int
	encode bool
}

// NewResponseWriter creates new empty ResponseWriter
func NewResponseWriter() *ResponseWriter {
	res := &ResponseWriter{
		header: make(http.Header),
		status: 200,
	}

	// set default content-type
	res.header.Set("content-type", "application/json")

	return res
}

// Header returns http.Header. You can modify it to send response header
func (r *ResponseWriter) Header() http.Header {
	return r.header
}

// Write appends chunk to response body
func (r *ResponseWriter) Write(chunk []byte) (int, error) {
	return r.buffer.Write(chunk)
}

// WriteHeader set status code of current request
func (r *ResponseWriter) WriteHeader(statusCode int) {
	r.status = statusCode
}

func (r *ResponseWriter) bodyString() string {
	if !r.encode {
		return r.buffer.String()
	}
	return base64.RawStdEncoding.EncodeToString(r.buffer.Bytes())
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
