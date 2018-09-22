package middleware

import (
	"net/http"

	"github.com/palantir/stacktrace"
)

// type Middleware struct {
// 	wrappedRoundTripper http.RoundTripper
// }

// type RoundTripper interface {
// 	RoundTrip(req *http.Request) (*http.Response, error)
// }

type MiddlewareFunc func(parent http.RoundTripper) http.RoundTripper

// func compose2(middlewarefunc1, middlewarefunc2 MiddlewareFunc) MiddlewareFunc {
// 	return func(parent http.RoundTripper) http.RoundTripper {
// 		return middlewarefunc1(middlewarefunc2(parent))
// 	}
// }

// Compose returns a middlewareFunc that successively applies all the middlewareFuncs passed in the parameters
// Pay attention, Compose() applies the middleware right to left: order matters when appyling your middlewareFuncs
func Compose(middlewareFuncs ...MiddlewareFunc) MiddlewareFunc {
	// fancy recursion, if you feel it's unclear it possible to rewrite it in a procedural way
	// check out http://nauvalatmaja.com/2016/04/15/function-composition-in-go/
	return func(parent http.RoundTripper) http.RoundTripper {
		// firstMiddlewareFunc := middlewareFuncs[0]
		// restMiddlewareFunc := middlewareFuncs[1:len(middlewareFuncs)]
		// if len(restMiddlewareFunc) == 1 {
		// 	return firstMiddlewareFunc(parent)
		// }
		// return firstMiddlewareFunc(Compose(restMiddlewareFunc...)(parent))

		// nono := sort.Reverse

		for i := len(middlewareFuncs) - 1; i >= 0; i-- {
			parent = middlewareFuncs[i](parent)
		}

		// for _, middlewareFunc := range sort.Reverse(middlewareFuncs) {
		// 	parent = middlewareFunc(parent)
		// }
		return parent
	}
}

///////////////////////////////////////
// PRE_REQUEST

type PreRequestFunc func(req *http.Request) (*http.Request, error)
type PreRequestMiddleware struct {
	preRequestFunc      PreRequestFunc
	wrappedRoundTripper http.RoundTripper
}

func (preRequestMiddleware PreRequestMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// apply the preRequestFunc to the request first
	req, err := preRequestMiddleware.preRequestFunc(req)
	if err != nil {
		return nil, stacktrace.Propagate(err, "preRequestMiddleware failed")
	}

	// then perform the request
	resp, err := preRequestMiddleware.wrappedRoundTripper.RoundTrip(req)
	if err != nil {
		return resp, stacktrace.Propagate(err, "request failed")
	}

	return resp, nil
}

func NewPreRequestMiddleware(preRequestFunc PreRequestFunc) MiddlewareFunc {
	return func(parent http.RoundTripper) http.RoundTripper {
		return PreRequestMiddleware{preRequestFunc: preRequestFunc, wrappedRoundTripper: parent}
	}
}

///////////////////////////////////////
// POST_REQUEST

type PostResponseFunc func(resp *http.Response) (*http.Response, error)
type PostResponseMiddleware struct {
	postResponseFunc    PostResponseFunc
	wrappedRoundTripper http.RoundTripper
}

func (postResponseMiddleware PostResponseMiddleware) RoundTrip(req *http.Request) (*http.Response, error) {
	// perform the request first
	resp, err := postResponseMiddleware.wrappedRoundTripper.RoundTrip(req)
	if err != nil {
		return nil, stacktrace.Propagate(err, "request failed")
	}

	// then apply the postResponseFunc to the response
	resp, err = postResponseMiddleware.postResponseFunc(resp)
	if err != nil {
		return resp, stacktrace.Propagate(err, "postRequestMiddleware failed")
	}

	return resp, nil
}

func NewPostResponseMiddleware(postResponseFunc PostResponseFunc) MiddlewareFunc {
	return func(parent http.RoundTripper) http.RoundTripper {
		return PostResponseMiddleware{postResponseFunc: postResponseFunc, wrappedRoundTripper: parent}
	}
}
