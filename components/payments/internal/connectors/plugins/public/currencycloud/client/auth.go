package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) authenticate(ctx context.Context) (string, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "currencycloud", "authenticate")
	// now := time.Now()
	// defer f(ctx, now)

	form := make(url.Values)

	form.Add("login_id", c.loginID)
	form.Add("api_key", c.apiKey)

	req, err := http.NewRequest(http.MethodPost,
		c.buildEndpoint("v2/authenticate/api"), strings.NewReader(form.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to do get request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	//nolint:tagliatelle // allow for client code
	type response struct {
		AuthToken string `json:"auth_token"`
	}

	var res response

	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", fmt.Errorf("failed to decode response body: %w", err)
	}

	return res.AuthToken, nil
}
