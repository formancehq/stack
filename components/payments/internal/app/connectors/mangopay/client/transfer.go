package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Funds struct {
	Currency string      `json:"Currency"`
	Amount   json.Number `json:"Amount"`
}

type TransferRequest struct {
	AuthorID         string `json:"AuthorId"`
	CreditedUserID   string `json:"CreditedUserId,omitempty"`
	DebitedFunds     Funds  `json:"DebitedFunds"`
	Fees             Funds  `json:"Fees"`
	DebitedWalletID  string `json:"DebitedWalletId"`
	CreditedWalletID string `json:"CreditedWalletId"`
}

type TransferResponse struct {
	ID               string `json:"Id"`
	CreationDate     int64  `json:"CreationDate"`
	AuthorID         string `json:"AuthorId"`
	CreditedUserID   string `json:"CreditedUserId"`
	DebitedFunds     Funds  `json:"DebitedFunds"`
	Fees             Funds  `json:"Fees"`
	CreditedFunds    Funds  `json:"CreditedFunds"`
	Status           string `json:"Status"`
	ResultCode       string `json:"ResultCode"`
	ResultMessage    string `json:"ResultMessage"`
	Type             string `json:"Type"`
	ExecutionDate    int64  `json:"ExecutionDate"`
	Nature           string `json:"Nature"`
	DebitedWalletID  string `json:"DebitedWalletId"`
	CreditedWalletID string `json:"CreditedWalletId"`
}

func (c *Client) InitiateWalletTransfer(ctx context.Context, transferRequest *TransferRequest) (*TransferResponse, error) {
	endpoint := fmt.Sprintf("%s/v2.01/%s/transfers", c.endpoint, c.clientID)

	body, err := json.Marshal(transferRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

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

	var transferResponse TransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return &transferResponse, nil
}

func (c *Client) GetWalletTransfer(ctx context.Context, transferID string) (*TransferResponse, error) {
	endpoint := fmt.Sprintf("%s/v2.01/%s/transfers/%s", c.endpoint, c.clientID, transferID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

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

	var transfer TransferResponse
	if err := json.NewDecoder(resp.Body).Decode(&transfer); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return &transfer, nil
}
