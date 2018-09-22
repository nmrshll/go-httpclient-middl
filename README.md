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
	fmt.Println("Hello world !")
}
```
