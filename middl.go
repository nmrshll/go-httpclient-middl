package middl

import (
	"net/http"

	"github.com/nmrshll/go-httpclient-middl/middleware"
	"github.com/palantir/stacktrace"
)

// Client is a wrapper around http.Client that you can add middleware to
type Client struct {
	*http.Client
}

// NewClient creates a new HttpClient
func NewClient(clients ...*http.Client) (*Client, error) {
	// validate parameters
	{
		if len(clients) > 1 {
			return nil, stacktrace.NewError("can't instantiate Client with more than one http.Client")
		}
	}

	// if a client is given then use this one
	if len(clients) == 1 {
		client := clients[0]
		if client.Transport == nil {
			client.Transport = http.DefaultTransport
		}
		return &Client{client}, nil
	}

	// if no client is given create a client using http.DefaultClient
	client := &Client{Client: &http.Client{Transport: http.DefaultTransport}}
	return client, nil
}

type MiddlewareFunc func(parent http.RoundTripper) http.RoundTripper

func (c *Client) UseMiddleware(middleware ...middleware.MiddlewareFunc) {
	// validate / modify parameters
	{
		// client should have transport
		if c.Transport == nil {
			c.Transport = http.DefaultTransport
		}
	}

	for _, middlewareFunc := range middleware {
		c.Transport = middlewareFunc(c.Transport)
	}
}
