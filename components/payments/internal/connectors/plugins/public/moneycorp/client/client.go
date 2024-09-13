package client

import (
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client
	endpoint   string
}

func newHTTPClient(clientID, apiKey, endpoint string) *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &apiTransport{
			clientID:   clientID,
			apiKey:     apiKey,
			endpoint:   endpoint,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func New(clientID, apiKey, endpoint string) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	c := &Client{
		httpClient: newHTTPClient(clientID, apiKey, endpoint),
		endpoint:   endpoint,
	}

	return c, nil
}
