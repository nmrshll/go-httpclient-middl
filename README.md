# go-httpclient-middl
Add middleware to your HTTP clients in Go

## Why ?
If you need an HTTP client that automatically logs requests, tracks metrics, joins a token with each request, validates HTTP status codes, or does any other custom action for every request/reponse, this library might be for you.

## How ?
Download the library using `go get -u github.com/nmrshll/go-httpclient-middl`

Then use it this way:
[embedmd]:# (.docs/examples/quickstart.go /func main/ $)
```go
func main() {
	httpClient := http.Client{Timeout: 30 * time.Second}
	client, err := middl.NewClient(&httpClient)
	if err != nil {
		log.Fatal(err)
	}

	// add middleware to you client (classic examples provided in this library or custom)
	client.UseMiddleware(logger.New())
	client.UseMiddleware(statusvalidator.New())

	// then do your requests as usual
	resp, err := client.Get("https://google.com")
	if err != nil {
		log.Fatal(err)
	}
	if resp == nil {
		log.Fatalf("no response from server")
	} // else
	defer resp.Body.Close()

	// do something with the response here
}
```
