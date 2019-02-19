package delta

import (
	"gotest.tools/assert"
	"bytes"
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func TestResponseWriter_ToAPIGWProxyResponse(t *testing.T) {
  w := NewResponseWriter()
  w.Header().Set("X-Powered-By", "Go")
  fmt.Fprintf(w, "Hello %s!", "world")

  agres := w.ToAPIGWProxyResponse()
  
  // by default status code is 200
  assert.Equal(t, agres.StatusCode, 200)

  // 
  assert.Equal(t, agres.)
}
