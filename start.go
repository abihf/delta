package delta

import (
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

// ServeOrStartLambda will start http server if it's not in lambda environment,
// otherwise it starts handling lambda
func ServeOrStartLambda(addr string, h http.Handler, opts ...Options) error {
	if _, ok := os.LookupEnv("LAMBDA_TASK_ROOT"); ok {
		Start(h, opts...)
		return nil
	}

	return http.ListenAndServe(addr, h)
}

// Start lambda server
//
// Example:
// mux := http.NewServeMux()
// mux.Handle("/", handeHandler)
// delta.Start(nil, mux)
func Start(h http.Handler, opts ...Options) {
	lambda.Start(NewHandler(h, opts...))
}
