package delta

import (
	"net/http"
)

func convertToHTTPHeader(headers map[string]string) http.Header {
	result := make(http.Header)
	for name, value := range headers {
		result.Set(name, value)
	}
	return result
}

func convertToLambdaHeader(headers http.Header) map[string]string {
	res := make(map[string]string)
	for name := range headers {
		res[name] = headers.Get(name)
	}
	return res
}
