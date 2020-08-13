package delta

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"gotest.tools/assert"
)

func Test_NewRequestV2(t *testing.T) {
	figure := `{
			"version": "2.0",
			"routeKey": "$default",
			"rawPath": "/my/path",
			"rawQueryString": "parameter1=value1&parameter1=value2&parameter2=value",
			"cookies": [
				"cookie1",
				"cookie2"
			],
			"headers": {
				"Header1": "value1",
				"Header2": "value1,value2"
			},
			"queryStringParameters": {
				"parameter1": "value1,value2",
				"parameter2": "value"
			},
			"requestContext": {
				"accountId": "123456789012",
				"apiId": "api-id",
				"authorizer": {
					"jwt": {
						"claims": {
							"claim1": "value1",
							"claim2": "value2"
						},
						"scopes": [
							"scope1",
							"scope2"
						]
					}
				},
				"domainName": "id.execute-api.us-east-1.amazonaws.com",
				"domainPrefix": "id",
				"http": {
					"method": "POST",
					"path": "/my/path",
					"protocol": "HTTP/1.1",
					"sourceIp": "IP",
					"userAgent": "agent"
				},
				"requestId": "id",
				"routeKey": "$default",
				"stage": "$default",
				"time": "12/Mar/2020:19:03:58 +0000",
				"timeEpoch": 1583348638390
			},
			"body": "Hello from Lambda",
			"pathParameters": {
				"parameter1": "value1"
			},
			"isBase64Encoded": false,
			"stageVariables": {
				"stageVariable1": "value1",
				"stageVariable2": "value2"
			}
		}`

	var e events.APIGatewayV2HTTPRequest
	json.Unmarshal([]byte(figure), &e)
	req, err := NewRequestV2(context.Background(), &e)
	if err != nil {
		t.Errorf("Failed to create request object %+v", err)
	}

	assert.Equal(t, req.Header.Get("host"), "example.com")
	assert.Equal(t, req.URL.String(), "/hello?a=1&b=2")
	assert.Equal(t, req.URL.Query().Get("a"), "1")
}

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
