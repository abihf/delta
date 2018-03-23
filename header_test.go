package delta

import (
	"net/http"
	"testing"

	"github.com/gotestyourself/gotestyourself/assert"
)

func Test_convertToHTTPHeader(t *testing.T) {
	orig := map[string]string{
		"A": "1",
		"B": "2",
	}
	coverted := convertToHTTPHeader(orig)

	assert.Equal(t, coverted.Get("A"), "1")
	assert.Equal(t, coverted.Get("b"), "2")
	assert.Equal(t, coverted.Get("c"), "")

}

func Test_convertToLambdaHeader(t *testing.T) {
	orig := make(http.Header)
	orig.Set("X", "1")
	orig.Add("y-z", "2")

	coverted := convertToLambdaHeader(orig)

	assert.Equal(t, coverted["X"], "1")
	assert.Equal(t, coverted["Y-Z"], "2")
}
