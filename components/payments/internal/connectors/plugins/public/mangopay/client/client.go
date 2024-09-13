package client

import (
	"strings"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
	"golang.org/x/oauth2/clientcredentials"
)

// TODO(polo): Fetch Client wallets (FEES, ...) in the future
type Client struct {
	httpClient httpwrapper.Client

	clientID string
	endpoint string
}

func New(clientID, apiKey, endpoint string) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	config := &httpwrapper.Config{
		OAuthConfig: &clientcredentials.Config{
			ClientID:     clientID,
			ClientSecret: apiKey,
			TokenURL:     endpoint + "/v2.01/oauth/token",
		},
	}
	httpClient, err := httpwrapper.NewClient(config)

	c := &Client{
		httpClient: httpClient,

		clientID: clientID,
		endpoint: endpoint,
	}
	return c, err
}
