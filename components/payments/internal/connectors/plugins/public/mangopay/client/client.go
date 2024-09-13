package client

import (
	"context"
	"net/http"
	"strings"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"golang.org/x/oauth2/clientcredentials"
)

// TODO(polo): Fetch Client wallets (FEES, ...) in the future
type Client struct {
	httpClient *http.Client

	clientID string
	endpoint string
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

func New(clientID, apiKey, endpoint string) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	c := &Client{
		httpClient: newHTTPClient(clientID, apiKey, endpoint),

		clientID: clientID,
		endpoint: endpoint,
	}

	return c, nil
}
