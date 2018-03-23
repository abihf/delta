# delta
Use golang http handler in AWS Lambda

## usage
```go
func apiHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("x-powered-by", "delta")
  fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func main() {
  mux := http.NewServeMux()
  mux.HandleFunc("/api/", apiHandler)

  // start
  delta.Start(mux)
}
```

## License
[MIT](LICENSE)
