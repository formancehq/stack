package client

import (
	"fmt"
	"net/http"

	"github.com/formancehq/payments/genericclient"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type apiTransport struct {
	APIKey     string
	underlying http.RoundTripper
}

func (t *apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t.APIKey))

	return t.underlying.RoundTrip(req)
}

type Client struct {
	apiClient *genericclient.APIClient
}

func New(apiKey, baseURL string) *Client {
	httpClient := &http.Client{
		Transport: &apiTransport{
			APIKey:     apiKey,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}

	configuration := genericclient.NewConfiguration()
	configuration.HTTPClient = httpClient
	configuration.Servers[0].URL = baseURL

	genericClient := genericclient.NewAPIClient(configuration)

	return &Client{
		apiClient: genericClient,
	}
}
