package statusvalidator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/palantir/stacktrace"
	"github.com/teamwork/synthesis/utility/httpClient/middleware"
)

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	var outputSlice []string
	// Format the URL
	outputSlice = append(outputSlice,
		fmt.Sprintf("%v %v", r.Method, r.URL), // URL
	)
	// Loop through headers
	for name, headers := range r.Header {
		for _, h := range headers {
			outputSlice = append(outputSlice, fmt.Sprintf("%v: %v", name, h))
		}
	}

	// If this is a POST, add post data
	if r.Method == "POST" {
		r.ParseForm()
		outputSlice = append(outputSlice, "\n")
		outputSlice = append(outputSlice, r.Form.Encode())
	}
	// Return the request as a string
	return strings.Join(outputSlice, "\n")
}

func FormatResponse(resp *http.Response) (string, error) {
	var outputSlice []string
	// format the URL
	outputSlice = append(outputSlice, fmt.Sprintf("%v %v", resp.Request.Method, resp.Request.URL))

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", stacktrace.Propagate(err, "failed reading response body")
	}
	// re-write body into response if needed in next middlewares
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	// format the body bytes
	outputSlice = append(outputSlice, string(bodyBytes))

	return strings.Join(outputSlice, "\n"), nil
}

func DumpResponse(resp *http.Response) error {
	formatted, err := FormatResponse(resp)
	if err != nil {
		return err
	}
	fmt.Println(formatted)

	return nil
}

func Printfln(format string, a ...interface{}) (int, error) {
	return fmt.Println(fmt.Sprintf(format, a...))
}

func New() middleware.MiddlewareFunc {
	return middleware.NewPostResponseMiddleware(func(resp *http.Response) (*http.Response, error) {
		if resp.StatusCode >= 200 && resp.StatusCode <= 399 {
			Printfln("HTTP %s: %s", resp.Status, resp.Request.URL)
		} else {
			err := DumpResponse(resp)
			if err != nil {
				return resp, stacktrace.Propagate(err, "failed dumping response")
			}
			return resp, stacktrace.NewError("HTTP status code is outside the success (2xx-3xx) range")
		}
		return resp, nil
	})
}
