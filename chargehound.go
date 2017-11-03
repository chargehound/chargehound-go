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
	version    = "6.0.0"

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

// Chargehound client params
type ClientParams struct {
	// The Chargehound API key used to interact with the API (Required).
	APIKey string
	// The API version (Optional).
	APIVersion string
}

// Creates a new chargehound client with the specified api key and the default configuration.
func New(params *ClientParams) *Client {

	ch := Client{
		APIKey:     params.APIKey,
		Basepath:   basepath,
		Host:       host,
		HTTPClient: &http.Client{Timeout: defaultHTTPTimeout},
		Protocol:   protocol,
		Version:    version,
		APIVersion: params.APIVersion,
	}

	ch.Disputes = &Disputes{client: &ch}

	return &ch
}
