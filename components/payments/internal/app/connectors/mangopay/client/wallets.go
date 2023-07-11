package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type wallet struct {
	ID           string `json:"Id"`
	Description  string `json:"Description"`
	CreationDate int64  `json:"CreationDate"`
	Currency     string `json:"Currency"`
}

func (c *Client) GetWallets(ctx context.Context, userID string, page int) ([]*wallet, error) {
	endpoint := fmt.Sprintf("%s/v2.01/%s/users/%s/wallets", c.endpoint, c.clientID, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("per_page", "100")
	q.Add("page", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var wallets []*wallet
	if err := json.NewDecoder(resp.Body).Decode(&wallets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login response body: %w", err)
	}

	return wallets, nil
}
