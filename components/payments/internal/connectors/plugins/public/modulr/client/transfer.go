package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type DestinationType string

const (
	DestinationTypeAccount     DestinationType = "ACCOUNT"
	DestinationTypeBeneficiary DestinationType = "BENEFICIARY"
)

type Destination struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

type TransferRequest struct {
	SourceAccountID   string      `json:"sourceAccountId"`
	Destination       Destination `json:"destination"`
	Currency          string      `json:"currency"`
	Amount            json.Number `json:"amount"`
	Reference         string      `json:"reference"`
	ExternalReference string      `json:"externalReference"`
	PaymentDate       string      `json:"paymentDate"`
}

type getTransferResponse struct {
	Content []TransferResponse `json:"content"`
}

type TransferResponse struct {
	ID                string `json:"id"`
	Status            string `json:"status"`
	CreatedDate       string `json:"createdDate"`
	ExternalReference string `json:"externalReference"`
	ApprovalStatus    string `json:"approvalStatus"`
	Message           string `json:"message"`
}

func (c *Client) InitiateTransfer(ctx context.Context, transferRequest *TransferRequest) (*TransferResponse, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "modulr", "initiate_transfer")
	// now := time.Now()
	// defer f(ctx, now)

	body, err := json.Marshal(transferRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, c.buildEndpoint("payments"), bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create transfer request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate transfer: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		// TODO(polo): retryable errors
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res TransferResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTransfer(ctx context.Context, transferID string) (TransferResponse, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "modulr", "get_transfer")
	// now := time.Now()
	// defer f(ctx, now)

	resp, err := c.httpClient.Get(c.buildEndpoint("payments?id=%s", transferID))
	if err != nil {
		return TransferResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO(polo): retryable errors
		return TransferResponse{}, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res getTransferResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return TransferResponse{}, err
	}

	if len(res.Content) == 0 {
		return TransferResponse{}, fmt.Errorf("transfer not found")
	}

	return res.Content[0], nil
}
