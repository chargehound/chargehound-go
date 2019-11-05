// Package chargehound contains Go bindings for the Chargehound API.
package chargehound

import (
	"net/http"
	"time"
)

const (
	APIVersion = ""
	basepath   = "/v1/"
	host       = "api.chargehound.com"
	protocol   = "https://"
	version    = "8.3.1"

	defaultHTTPTimeout = 60 * time.Second
)

type Client struct {
	// The Chargehound API key used to interact with the API.
	APIKey string
	// The API host.
	Host string
	// The API scheme.
	Protocol string
	// The API base path.
	Basepath string
	// The client version.
	Version string
	// The API version.
	APIVersion string
	// The client http timeout.
	HTTPClient *http.Client
	// The disputes resource.
	Disputes *Disputes
}

// Chargehound client optional params.
type ClientParams struct {
	// The API version
	APIVersion string
}

// Creates a new chargehound client with the specified api key and the default configuration.
func New(key string, params *ClientParams) *Client {

	var apiVersion string
	if params != nil {
		apiVersion = params.APIVersion
	} else {
		apiVersion = APIVersion
	}

	ch := Client{
		APIKey:     key,
		Basepath:   basepath,
		Host:       host,
		HTTPClient: &http.Client{Timeout: defaultHTTPTimeout},
		Protocol:   protocol,
		Version:    version,
		APIVersion: apiVersion,
	}

	ch.Disputes = &Disputes{client: &ch}

	return &ch
}
