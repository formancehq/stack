package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type payoutRequest struct {
	Payout struct {
		Attributes *PayoutRequest `json:"attributes"`
	} `json:"data"`
}

type PayoutRequest struct {
	SourceAccountID  string  `json:"-"`
	IdempotencyKey   string  `json:"-"`
	RecipientID      string  `json:"recipientId"`
	PaymentDate      string  `json:"paymentDate"`
	PaymentAmount    float64 `json:"paymentAmount"`
	PaymentCurrency  string  `json:"paymentCurrency"`
	PaymentMethgod   string  `json:"paymentMethod"`
	PaymentReference string  `json:"paymentReference"`
	ClientReference  string  `json:"clientReference"`
	PaymentPurpose   string  `json:"paymentPurpose"`
}

type payoutResponse struct {
	Payout *PayoutResponse `json:"data"`
}

type PayoutResponse struct {
	ID         string `json:"id"`
	Attributes struct {
		AccountID        string  `json:"accountId"`
		PaymentAmount    float64 `json:"paymentAmount"`
		PaymentCurrency  string  `json:"paymentCurrency"`
		PaymentApproved  bool    `json:"paymentApproved"`
		PaymentStatus    string  `json:"paymentStatus"`
		PaymentMethod    string  `json:"paymentMethod"`
		PaymentDate      string  `json:"paymentDate"`
		PaymentValueDate string  `json:"paymentValueDate"`
		RecipientDetails struct {
			RecipientID int32 `json:"recipientId"`
		} `json:"recipientDetails"`
		PaymentReference string `json:"paymentReference"`
		ClientReference  string `json:"clientReference"`
		CreatedAt        string `json:"createdAt"`
		CreatedBy        string `json:"createdBy"`
		UpdatedAt        string `json:"updatedAt"`
		PaymentPurpose   string `json:"paymentPurpose"`
	} `json:"attributes"`
}

func (c *Client) InitiatePayout(ctx context.Context, pr *PayoutRequest) (*PayoutResponse, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/payments", c.endpoint, pr.SourceAccountID)

	reqBody := &payoutRequest{}
	reqBody.Payout.Attributes = pr
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal payout request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", pr.IdempotencyKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res payoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res.Payout, nil
}

func (c *Client) GetPayout(ctx context.Context, accountID string, payoutID string) (*PayoutResponse, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/payments/%s", c.endpoint, accountID, payoutID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create get payout request request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get account balances: %w", err)
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

	var payoutResponse payoutResponse
	if err := json.NewDecoder(resp.Body).Decode(&payoutResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return payoutResponse.Payout, nil
}
