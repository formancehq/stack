package client

import (
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client

	endpoint string

	logger logging.Logger
}

func newHTTPClient(clientID, apiKey, endpoint string, logger logging.Logger) *http.Client {
	return &http.Client{
		Timeout: 10 * time.Second,
		Transport: &apiTransport{
			logger:     logger,
			clientID:   clientID,
			apiKey:     apiKey,
			endpoint:   endpoint,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func NewClient(clientID, apiKey, endpoint string, logger logging.Logger) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	c := &Client{
		httpClient: newHTTPClient(clientID, apiKey, endpoint, logger),

		endpoint: endpoint,

		logger: logger,
	}

	return c, nil
}
