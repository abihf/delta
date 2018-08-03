package delta

import (
	"net/http"
)

// APIGWProxyHeader is format that used by API Gateway to store http header
type APIGWProxyHeader map[string]string

// Header wrap normal http.Header and add some converter
type Header struct {
	http.Header
}

// ToAPIGWProxyHeader convert http.Header to LambdaHeader
func (h *Header) ToAPIGWProxyHeader() APIGWProxyHeader {
	res := make(APIGWProxyHeader)

	for name := range h.Header {
		res[name] = h.Get(name)
	}

	return res
}

// HeaderFromAPIGWProxyHeader creates new Header from APIGWProxyHeader
func HeaderFromAPIGWProxyHeader(ph APIGWProxyHeader) *Header {
	header := &Header{make(http.Header)}
	for name, value := range ph {
		header.Set(name, value)
	}
	return header
}
