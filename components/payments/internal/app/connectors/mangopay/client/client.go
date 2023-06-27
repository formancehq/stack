package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/zitadel/oidc/pkg/oidc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client struct {
	httpClient *http.Client

	clientID string
	apiKey   string
	endpoint string

	logger logging.Logger

	accessToken          string
	accessTokenExpiresAt time.Time
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func NewClient(clientID, apiKey, endpoint string, logger logging.Logger) (*Client, error) {
	endpoint = strings.TrimSuffix(endpoint, "/")

	c := &Client{
		httpClient: newHTTPClient(),

		clientID: clientID,
		apiKey:   apiKey,
		endpoint: endpoint,

		logger: logger,
	}

	if err := c.login(context.Background()); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Client) login(ctx context.Context) error {
	form := url.Values{
		"grant_type": []string{string(oidc.GrantTypeClientCredentials)},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.endpoint+"/v2.01/oauth/token", bytes.NewBufferString(form.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(c.clientID, c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read login response body: %w", err)
	}

	type response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	var res response

	if err = json.Unmarshal(responseBody, &res); err != nil {
		return fmt.Errorf("failed to unmarshal login response: %w", err)
	}

	c.accessToken = res.AccessToken
	c.accessTokenExpiresAt = time.Now().Add(time.Duration(res.ExpiresIn) * time.Second)

	return nil
}

func (c *Client) ensureAccessTokenIsValid(ctx context.Context) error {
	if c.accessTokenExpiresAt.After(time.Now()) {
		return nil
	}

	return c.login(ctx)
}
