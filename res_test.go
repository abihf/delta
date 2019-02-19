package delta

import (
	"fmt"
	"testing"

	"gotest.tools/assert"
)

func TestResponseWriter_ToAPIGWProxyResponse(t *testing.T) {
	w := NewResponseWriter()
	w.Header().Set("X-Powered-By", "Go")
	fmt.Fprintf(w, "Hello %s!", "world")

	agres := w.ToAPIGWProxyResponse()

	// by default status code is 200
	assert.Equal(t, agres.StatusCode, 200)
}
