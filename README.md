# δelta
Use golang http handler in AWS Lambda

## usage
```go
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("x-powered-by", "delta")
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello/", helloHandler)

	// start lambda handling if it runs on lambda
	// otherwise start http server on port 3000
	delta.ServeOrStartLambda(":3000", mux)
}

```

## License
[MIT](LICENSE)
