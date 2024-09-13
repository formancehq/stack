package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type PayoutRequest struct {
	AuthorID            string `json:"AuthorId"`
	DebitedFunds        Funds  `json:"DebitedFunds"`
	Fees                Funds  `json:"Fees"`
	DebitedWalletID     string `json:"DebitedWalletId"`
	BankAccountID       string `json:"BankAccountId"`
	BankWireRef         string `json:"BankWireRef,omitempty"`
	PayoutModeRequested string `json:"PayoutModeRequested,omitempty"`
}

type PayoutResponse struct {
	ID              string `json:"Id"`
	ModeRequest     string `json:"ModeRequested"`
	ModeApplied     string `json:"ModeApplied"`
	FallbackReason  string `json:"FallbackReason"`
	CreationDate    int64  `json:"CreationDate"`
	AuthorID        string `json:"AuthorId"`
	DebitedFunds    Funds  `json:"DebitedFunds"`
	Fees            Funds  `json:"Fees"`
	CreditedFunds   Funds  `json:"CreditedFunds"`
	Status          string `json:"Status"`
	ResultCode      string `json:"ResultCode"`
	ResultMessage   string `json:"ResultMessage"`
	Type            string `json:"Type"`
	Nature          string `json:"Nature"`
	ExecutionDate   int64  `json:"ExecutionDate"`
	BankAccountID   string `json:"BankAccountId"`
	DebitedWalletID string `json:"DebitedWalletId"`
	PaymentType     string `json:"PaymentType"`
	BankWireRef     string `json:"BankWireRef"`
}

func (c *Client) InitiatePayout(ctx context.Context, payoutRequest *PayoutRequest) (*PayoutResponse, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "initiate_payout")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/payouts/bankwire", c.endpoint, c.clientID)

	body, err := json.Marshal(payoutRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get wallets: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			// TODO: log error
			_ = err
		}
	}()

	if resp.StatusCode != http.StatusOK {
		// Never retry payout initiation
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var payoutResponse PayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&payoutResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return &payoutResponse, nil
}

func (c *Client) GetPayout(ctx context.Context, payoutID string) (*PayoutResponse, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "mangopay", "get_payout")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/v2.01/%s/payouts/%s", c.endpoint, c.clientID, payoutID)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create get payout request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get payout: %w", err)
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

	var payoutResponse PayoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&payoutResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal payout response body: %w", err)
	}

	return &payoutResponse, nil
}
