package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
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

	var res TransferResponse
	var errRes modulrError
	_, err = c.httpClient.Do(req, &res, &errRes)
	switch err {
	case nil:
		return &res, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to initiate transfer: %w", err)
}

func (c *Client) GetTransfer(ctx context.Context, transferID string) (TransferResponse, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "modulr", "get_transfer")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildEndpoint("payments?id=%s", transferID), nil)
	if err != nil {
		return TransferResponse{}, fmt.Errorf("failed to create get transfer request: %w", err)
	}

	var res getTransferResponse
	var errRes modulrError
	_, err = c.httpClient.Do(req, &res, &errRes)
	switch err {
	case nil:
		if len(res.Content) == 0 {
			return TransferResponse{}, fmt.Errorf("transfer not found")
		}
		return res.Content[0], nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return TransferResponse{}, errRes.Error()
	}
	return TransferResponse{}, fmt.Errorf("failed to get transactions %w", err)
}
