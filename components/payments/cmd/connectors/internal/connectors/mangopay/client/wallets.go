package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type Wallet struct {
	ID           string   `json:"Id"`
	Owners       []string `json:"Owners"`
	Description  string   `json:"Description"`
	CreationDate int64    `json:"CreationDate"`
	Currency     string   `json:"Currency"`
	Balance      struct {
		Currency string      `json:"Currency"`
		Amount   json.Number `json:"Amount"`
	} `json:"Balance"`
}

func (c *Client) GetWallets(ctx context.Context, userID string, page, pageSize int) ([]*Wallet, error) {
	f := connectors.ClientMetrics(ctx, "mangopay", "list_wallets")
	now := time.Now()
	defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/users/%s/wallets", c.endpoint, c.clientID, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("per_page", strconv.Itoa(pageSize))
	q.Add("page", fmt.Sprint(page))
	q.Add("Sort", "CreationDate:ASC")
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

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalErrorWithRetry(resp.StatusCode, resp.Body).Error()
	}

	var wallets []*Wallet
	if err := json.NewDecoder(resp.Body).Decode(&wallets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return wallets, nil
}

func (c *Client) GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	f := connectors.ClientMetrics(ctx, "mangopay", "get_wallets")
	now := time.Now()
	defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/wallets/%s", c.endpoint, c.clientID, walletID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalErrorWithRetry(resp.StatusCode, resp.Body).Error()
	}

	var wallet Wallet
	if err := json.NewDecoder(resp.Body).Decode(&wallet); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallet response body: %w", err)
	}

	return &wallet, nil
}
