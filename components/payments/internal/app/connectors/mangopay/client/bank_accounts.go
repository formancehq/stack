package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type bankAccount struct {
	ID           string `json:"Id"`
	OwnerName    string `json:"OwnerName"`
	CreationDate int64  `json:"CreationDate"`
}

func (c *Client) GetBankAccounts(ctx context.Context, userID string, page int) ([]*bankAccount, error) {
	endpoint := fmt.Sprintf("%s/v2.01/%s/users/%s/bankaccounts", c.endpoint, c.clientID, userID)
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
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var bankAccounts []*bankAccount
	if err := json.NewDecoder(resp.Body).Decode(&bankAccounts); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return bankAccounts, nil
}
