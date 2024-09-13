package client

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	httpClient *http.Client
	endpoint   string
	loginID    string
	apiKey     string
}

func (c *Client) buildEndpoint(path string, args ...interface{}) string {
	return fmt.Sprintf("%s/%s", c.endpoint, fmt.Sprintf(path, args...))
}

const DevAPIEndpoint = "https://devapi.currencycloud.com"

func newAuthenticatedHTTPClient(authToken string) *http.Client {
	return &http.Client{
		Transport: &apiTransport{
			authToken:  authToken,
			underlying: otelhttp.NewTransport(http.DefaultTransport),
		},
	}
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

// New creates a new client for the CurrencyCloud API.
func New(loginID, apiKey, endpoint string) (*Client, error) {
	if endpoint == "" {
		endpoint = DevAPIEndpoint
	}

	c := &Client{
		httpClient: newHTTPClient(),
		endpoint:   endpoint,
		loginID:    loginID,
		apiKey:     apiKey,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Tokens expire after 30 minutes of inactivity which should not be the case
	// for us since we're polling the API frequently.
	// TODO(polo): add refreh
	authToken, err := c.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	c.httpClient = newAuthenticatedHTTPClient(authToken)

	return c, nil
}
