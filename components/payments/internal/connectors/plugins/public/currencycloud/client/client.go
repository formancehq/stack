package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type apiTransport struct {
	authToken  string
	underlying *otelhttp.Transport
}

func (t *apiTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("X-Auth-Token", t.authToken)

	return t.underlying.RoundTrip(req)
}

type Client struct {
	httpClient httpwrapper.Client
	endpoint   string
	loginID    string
	apiKey     string
}

func (c *Client) buildEndpoint(path string, args ...interface{}) string {
	return fmt.Sprintf("%s/%s", c.endpoint, fmt.Sprintf(path, args...))
}

const DevAPIEndpoint = "https://devapi.currencycloud.com"

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

// New creates a new client for the CurrencyCloud API.
func New(ctx context.Context, loginID, apiKey, endpoint string) (*Client, error) {
	if endpoint == "" {
		endpoint = DevAPIEndpoint
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c := &Client{
		endpoint: endpoint,
		loginID:  loginID,
		apiKey:   apiKey,
	}

	// Tokens expire after 30 minutes of inactivity which should not be the case
	// for us since we're polling the API frequently.
	// TODO(polo): add refreh
	authToken, err := c.authenticate(ctx, newHTTPClient())
	if err != nil {
		return nil, err
	}

	config := &httpwrapper.Config{
		Transport: &apiTransport{
			authToken:  authToken,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
	httpClient, err := httpwrapper.NewClient(config)
	if err != nil {
		return nil, err
	}
	c.httpClient = httpClient
	return c, nil
}
