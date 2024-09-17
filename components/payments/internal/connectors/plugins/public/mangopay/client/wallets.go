package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
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

func (c *Client) GetWallets(ctx context.Context, userID string, page, pageSize int) ([]Wallet, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "list_wallets")
	// now := time.Now()
	// defer f(ctx, now)

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

	var wallets []Wallet
	var errRes mangopayError
	_, err = c.httpClient.Do(req, &wallets, errRes)
	switch err {
	case nil:
		return wallets, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get wallets %w", err)
}

func (c *Client) GetWallet(ctx context.Context, walletID string) (*Wallet, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "get_wallets")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/wallets/%s", c.endpoint, c.clientID, walletID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet request: %w", err)
	}

	var wallet Wallet
	var errRes mangopayError
	_, err = c.httpClient.Do(req, &wallet, errRes)
	switch err {
	case nil:
		return &wallet, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get wallet %w", err)
}
