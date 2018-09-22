package logger

import (
	"fmt"
	"net/http"

	"github.com/teamwork/synthesis/utility/httpClient/middleware"
)

func New() middleware.MiddlewareFunc {
	return middleware.Compose(
		NewPreRequestLogger(),
		NewPostResponseLogger(),
	)
}

func NewPreRequestLogger() middleware.MiddlewareFunc {
	return middleware.NewPreRequestMiddleware(func(req *http.Request) (*http.Request, error) {
		fmt.Printf("Sending request to %v\n", req.URL)
		return req, nil
	})
}

func NewPostResponseLogger() middleware.MiddlewareFunc {
	return middleware.NewPostResponseMiddleware(func(resp *http.Response) (*http.Response, error) {
		fmt.Printf("Received %v response\n", resp.Status)
		return resp, nil
	})
}
