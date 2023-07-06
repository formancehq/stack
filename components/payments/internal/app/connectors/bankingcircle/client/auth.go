package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

func (c *Client) login(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.authorizationEndpoint+"/api/v1/authorizations/authorize", http.NoBody)
	if err != nil {
		return fmt.Errorf("failed to create login request: %w", err)
	}

	req.SetBasicAuth(c.username, c.password)

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

	//nolint:tagliatelle // allow for client-side structures
	type response struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   string `json:"expires_in"`
	}

	var res response

	if err = json.Unmarshal(responseBody, &res); err != nil {
		return fmt.Errorf("failed to unmarshal login response: %w", err)
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
	if c.accessTokenExpiresAt.After(time.Now().Add(5 * time.Second)) {
		return nil
	}

	return c.login(ctx)
}
