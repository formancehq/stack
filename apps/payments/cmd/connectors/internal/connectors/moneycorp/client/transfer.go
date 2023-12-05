package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type transferRequest struct {
	Transfer struct {
		Attributes *TransferRequest `json:"attributes"`
	} `json:"data"`
}

type TransferRequest struct {
	SourceAccountID    string  `json:"-"`
	IdempotencyKey     string  `json:"-"`
	ReceivingAccountID string  `json:"receivingAccountId"`
	TransferAmount     float64 `json:"transferAmount"`
	TransferCurrency   string  `json:"transferCurrency"`
	TransferReference  string  `json:"transferReference,omitempty"`
	ClientReference    string  `json:"clientReference,omitempty"`
}

type transferResponse struct {
	Transfer *TransferResponse `json:"data"`
}

type TransferResponse struct {
	ID         string `json:"id"`
	Attributes struct {
		SendingAccountID     int64   `json:"sendingAccountId"`
		SendingAccountName   string  `json:"sendingAccountName"`
		ReceivingAccountID   int64   `json:"receivingAccountId"`
		ReceivingAccountName string  `json:"receivingAccountName"`
		CreatedAt            string  `json:"createdAt"`
		CreatedBy            string  `json:"createdBy"`
		UpdatedAt            string  `json:"updatedAt"`
		TransferReference    string  `json:"transferReference"`
		ClientReference      string  `json:"clientReference"`
		TransferDate         string  `json:"transferDate"`
		TransferAmount       float64 `json:"transferAmount"`
		TransferCurrency     string  `json:"transferCurrency"`
		TransferStatus       string  `json:"transferStatus"`
	} `json:"attributes"`
}

func (c *Client) InitiateTransfer(ctx context.Context, tr *TransferRequest) (*TransferResponse, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/transfers", c.endpoint, tr.SourceAccountID)

	reqBody := &transferRequest{}
	reqBody.Transfer.Attributes = tr
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Idempotency-Key", tr.IdempotencyKey)

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

	if resp.StatusCode != http.StatusCreated {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var transferResponse transferResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return transferResponse.Transfer, nil
}

func (c *Client) GetTransfer(ctx context.Context, accountID string, transferID string) (*TransferResponse, error) {
	endpoint := fmt.Sprintf("%s/accounts/%s/transfers/%s", c.endpoint, accountID, transferID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create get transfer request: %w", err)
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

	var transferResponse transferResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return transferResponse.Transfer, nil
}
