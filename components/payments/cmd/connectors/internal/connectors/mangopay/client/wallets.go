package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type wallet struct {
	ID           string `json:"Id"`
	Description  string `json:"Description"`
	CreationDate int64  `json:"CreationDate"`
	Currency     string `json:"Currency"`
	Balance      struct {
		Currency string      `json:"Currency"`
		Amount   json.Number `json:"Amount"`
	} `json:"Balance"`
}

func (c *Client) GetWallets(ctx context.Context, userID string, page int) ([]*wallet, error) {
	f := connectors.ClientMetrics(ctx, "mangopay", "list_wallets")
	now := time.Now()
	defer f(ctx, now)

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
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var wallets []*wallet
	if err := json.NewDecoder(resp.Body).Decode(&wallets); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return wallets, nil
}
