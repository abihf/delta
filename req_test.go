package delta

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"gotest.tools/assert"
)

func Test_NewRequest(t *testing.T) {
	e := &events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/hello",
		Body:       "body",

		Headers: map[string]string{
			"content-length": "4",
			"host":           "example.com",
		},

		QueryStringParameters: map[string]string{
			"a": "1",
			"b": "2",
		},
	}
	req, err := NewRequest(context.Background(), e)
	if err != nil {
		t.Errorf("Failed to create request object %+v", err)
	}

	assert.Equal(t, req.Header.Get("host"), "example.com")
	assert.Equal(t, req.URL.String(), "/hello?a=1&b=2")
	assert.Equal(t, req.URL.Query().Get("a"), "1")
}
