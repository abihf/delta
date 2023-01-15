package delta

import (
	"bytes"
	"encoding/base64"
	"errors"
	"io"
	"net/http"
)

var _ http.ResponseWriter = &ResponseWriter{}
var _ io.StringWriter = &ResponseWriter{}
var _ io.ReaderFrom = &ResponseWriter{}

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

	return res
}

// Header returns http.Header. You can modify it to send response header
func (rw *ResponseWriter) Header() http.Header {
	return rw.header
}

// Write appends chunk to response body
func (rw *ResponseWriter) Write(chunk []byte) (int, error) {
	return rw.buffer.Write(chunk)
}

// WriteHeader set status code of current request
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.status = statusCode
}

func (rw *ResponseWriter) bodyString() string {
	if !rw.encode {
		return rw.buffer.String()
	}
	return base64.RawStdEncoding.EncodeToString(rw.buffer.Bytes())
}

// ReadFrom implements
func (rw *ResponseWriter) ReadFrom(r io.Reader) (n int64, err error) {
	return rw.buffer.ReadFrom(r)
}

// WriteString implements io.StringWriter
func (rw *ResponseWriter) WriteString(s string) (n int, err error) {
	return rw.buffer.WriteString(s)
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
