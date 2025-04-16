package registry

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
)

// Server is a structure that holds the URL of the registry and the HTTP client
type Server struct {
	URL    string
	Client *http.Client
	Logger LoggerInterface
}

// NewServer Create a new Registry with the given URL and credentials, then Ping()s it
// before returning it to verify that the registry is available.
//
// You can, alternately, construct a Registry manually by populating the fields.
// This passes http.DefaultTransport to WrapTransport when creating the
// http.Client.
func NewServer(registryURL, username, password string, logger ...LoggerInterface) (*Server, error) {
	return newServerFromTransport(registryURL, username, password, http.DefaultTransport, logger...)
}

// NewInsecureServer Create a new Registry, as with New, using an http.Transport that disables
// SSL certificate verification.
func NewInsecureServer(registryURL, username, password string, logger ...LoggerInterface) (*Server, error) {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // nolint:gosec
		},
	}

	return newServerFromTransport(registryURL, username, password, transport, logger...)
}

func newServerFromTransport(
	registryURL,
	username,
	password string,
	transport http.RoundTripper,
	logger ...LoggerInterface,
) (*Server, error) {
	url := strings.TrimSuffix(registryURL, "/")
	transport = WrapTransport(transport, url, username, password)

	srv := &Server{
		URL: url,
		Client: &http.Client{
			Transport: transport,
		},
	}

	// Set the logger if provided
	srv.Logger = &defaultLogger{}
	if len(logger) > 0 && logger[0] != nil {
		srv.Logger = logger[0]
	}

	if err := srv.Ping(); err != nil {
		srv.Logger.Error("registry.ping error=%s", err)
		return nil, err
	}

	return srv, nil
}

// WrapTransport Given an existing http.RoundTripper such as http.DefaultTransport, build the
// transport stack necessary to authenticate to the Docker registry API. This
// adds in support for OAuth bearer tokens and HTTP Basic auth, and sets up
// error handling this library relies on.
func WrapTransport(transport http.RoundTripper, url, username, password string) http.RoundTripper {
	tokenTransport := &TokenTransport{
		Transport: transport,
		Username:  username,
		Password:  password,
	}
	basicAuthTransport := &BasicTransport{
		Transport: tokenTransport,
		URL:       url,
		Username:  username,
		Password:  password,
	}
	errorTransport := &ErrorTransport{
		Transport: basicAuthTransport,
	}
	return errorTransport
}

func (r *Server) url(pathTemplate string, args ...any) string {
	pathSuffix := fmt.Sprintf(pathTemplate, args...)
	url := fmt.Sprintf("%s%s", r.URL, pathSuffix)
	return url
}

// Ping the registry to verify that it is available
func (srv *Server) Ping() error {
	url := srv.url("/v2/")
	srv.Logger.Debug("registry.ping url=%s", url)
	resp, err := srv.Client.Get(url)
	if resp != nil {
		// nolint: errcheck
		defer resp.Body.Close()
	}
	return err
}

func addRegistryV2Header(req *http.Request) {
	// Docker Registry API v2
	req.Header.Set("Docker-Distribution-Api-Version", "registry/2.0")
}
