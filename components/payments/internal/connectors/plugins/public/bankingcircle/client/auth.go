package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

func (c *Client) login(ctx context.Context) error {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "bankingcircle", "authorize")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.authorizationEndpoint+"/api/v1/authorizations/authorize", http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)

	//nolint:tagliatelle // allow for client-side structures
	type response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}
	type responseError struct {
		ErrorCode string `json:"errorCode"`
		ErrorText string `json:"errorText"`
	}

	var res response
	var errors []responseError
	statusCode, err := c.httpClient.Do(req, &res, &errors)
	switch err {
	case nil:
		// fallthrough
	case httpwrapper.ErrStatusCodeUnexpected:
		if len(errors) > 0 {
			return fmt.Errorf("failed to login: %s %s", errors[0].ErrorCode, errors[0].ErrorText)
		}
		return fmt.Errorf("failed to login: %d", statusCode)
	}

	c.accessToken = res.AccessToken
	expiresIn, err := strconv.Atoi(res.ExpiresIn)
	if err != nil {
		return fmt.Errorf("failed to convert expires_in to int: %w", err)
	}
	c.accessTokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
	return nil
}

func (c *Client) ensureAccessTokenIsValid(ctx context.Context) error {
	if c.accessToken == "" {
		return c.login(ctx)
	}

	if c.accessTokenExpiresAt.After(time.Now().Add(5 * time.Second)) {
		return nil
	}

	return c.login(ctx)
}
