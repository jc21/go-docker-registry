package registry

import (
	"fmt"
	"io"
	"net/http"
)

// HTTPStatusError is returned when a non-successful HTTP response is received.
type HTTPStatusError struct {
	Response *http.Response
	// Copied from `Response.Body` to avoid problems with unclosed bodies later.
	// Nobody calls `err.Response.Body.Close()`, ever.
	Body []byte
}

// Error implements the error interface for HTTPStatusError.
func (err *HTTPStatusError) Error() string {
	return fmt.Sprintf("http: non-successful response (status=%v body=%q)", err.Response.StatusCode, err.Body)
}

var _ error = &HTTPStatusError{}

// ErrorTransport is a custom HTTP transport that checks for non-successful HTTP responses.
type ErrorTransport struct {
	Transport http.RoundTripper
}

// RoundTrip executes a single HTTP transaction and checks for non-successful responses.
func (t *ErrorTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	resp, err := t.Transport.RoundTrip(request)
	if err != nil {
		return resp, err
	}

	if resp.StatusCode >= 400 {
		// nolint: errcheck
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("http: failed to read response body (status=%v, err=%q)", resp.StatusCode, err)
		}

		return nil, &HTTPStatusError{
			Response: resp,
			Body:     body,
		}
	}

	return resp, err
}
