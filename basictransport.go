package registry

import (
	"net/http"
	"strings"
)

// BasicTransport is a custom HTTP transport that adds basic authentication to requests.
type BasicTransport struct {
	Transport http.RoundTripper
	URL       string
	Username  string
	Password  string
}

// RoundTrip executes a single HTTP transaction and adds basic authentication
// if the request URL matches the specified URL.
func (t *BasicTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if strings.HasPrefix(req.URL.String(), t.URL) {
		if t.Username != "" || t.Password != "" {
			req.SetBasicAuth(t.Username, t.Password)
		}
	}
	resp, err := t.Transport.RoundTrip(req)
	return resp, err
}
