package client

import (
	"fmt"
	"net/http"

	lru "github.com/hashicorp/golang-lru/v2"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

const apiEndpoint = "https://api.wise.com"

type apiTransport struct {
	APIKey     string
	underlying http.RoundTripper
}

func (t *apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.APIKey))

	return t.underlying.RoundTrip(req)
}

type Client struct {
	httpClient *http.Client

	recipientAccountsCache *lru.Cache[uint64, *RecipientAccount]
}

func (w *Client) endpoint(path string) string {
	return fmt.Sprintf("%s/%s", apiEndpoint, path)
}

func NewClient(apiKey string) *Client {
	recipientsCache, _ := lru.New[uint64, *RecipientAccount](2048)
	httpClient := &http.Client{
		Transport: &apiTransport{
			APIKey:     apiKey,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}

	return &Client{
		httpClient:             httpClient,
		recipientAccountsCache: recipientsCache,
	}
}
