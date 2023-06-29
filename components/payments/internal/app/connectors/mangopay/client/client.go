package client

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2/clientcredentials"
)

type Client struct {
	httpClient *http.Client

	clientID string
	endpoint string

	logger logging.Logger
}

func newHTTPClient(clientID, apiKey, endpoint string) *http.Client {
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: apiKey,
		TokenURL:     endpoint + "/v2.01/oauth/token",
	}

	httpClient := config.Client(context.Background())

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: otelhttp.NewTransport(httpClient.Transport),
	}
}

func NewClient(clientID, apiKey, endpoint string, logger logging.Logger) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	c := &Client{
		httpClient: newHTTPClient(clientID, apiKey, endpoint),

		clientID: clientID,
		endpoint: endpoint,

		logger: logger,
	}

	return c, nil
}
