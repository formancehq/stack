package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"time"
)

type PaymentAccount struct {
	Account              string `json:"account"`
	FinancialInstitution string `json:"financialInstitution"`
	Country              string `json:"country"`
}

type PaymentRequest struct {
	IdempotencyKey         string         `json:"idempotencyKey"`
	RequestedExecutionDate time.Time      `json:"requestedExecutionDate"`
	DebtorAccount          PaymentAccount `json:"debtorAccount"`
	DebtorReference        string         `json:"debtorReference"`
	CurrencyOfTransfer     string         `json:"currencyOfTransfer"`
	Amount                 struct {
		Currency string     `json:"currency"`
		Amount   *big.Float `json:"amount"`
	} `json:"amount"`
	ChargeBearer    string          `json:"chargeBearer"`
	CreditorAccount *PaymentAccount `json:"creditorAccount"`
}

type PaymentResponse struct {
	PaymentID string `json:"paymentId"`
	Status    string `json:"status"`
}

func (c *Client) InitiateTransferOrPayouts(ctx context.Context, transferRequest *PaymentRequest) (*PaymentResponse, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	body, err := json.Marshal(transferRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint+"/api/v1/payments/singles", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create payments request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make transfer: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("failed to make transfer: %w", err)
	}

	var transferResponse PaymentResponse
	if err := json.NewDecoder(resp.Body).Decode(&transferResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal wallets response body: %w", err)
	}

	return &transferResponse, nil
}
