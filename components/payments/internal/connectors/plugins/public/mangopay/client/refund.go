package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Refund struct {
	ID                     string `json:"Id"`
	Tag                    string `json:"Tag"`
	CreationDate           int64  `json:"CreationDate"`
	AuthorId               string `json:"AuthorId"`
	CreditedUserId         string `json:"CreditedUserId"`
	DebitedFunds           Funds  `json:"DebitedFunds"`
	CreditedFunds          Funds  `json:"CreditedFunds"`
	Fees                   Funds  `json:"Fees"`
	Status                 string `json:"Status"`
	ResultCode             string `json:"ResultCode"`
	ResultMessage          string `json:"ResultMessage"`
	ExecutionDate          int64  `json:"ExecutionDate"`
	Type                   string `json:"Type"`
	DebitedWalletId        string `json:"DebitedWalletId"`
	CreditedWalletId       string `json:"CreditedWalletId"`
	InitialTransactionID   string `json:"InitialTransactionId"`
	InitialTransactionType string `json:"InitialTransactionType"`
}

func (c *Client) GetRefund(ctx context.Context, refundID string) (*Refund, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "get_refund")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/refunds/%s", c.endpoint, c.clientID, refundID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get refund request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			// TODO(polo): log error
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var refund Refund
	if err := json.NewDecoder(resp.Body).Decode(&refund); err != nil {
		return nil, fmt.Errorf("failed to unmarshal refund response body: %w", err)
	}

	return &refund, nil
}
