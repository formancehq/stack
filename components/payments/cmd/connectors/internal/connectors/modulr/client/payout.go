package client

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
)

type PayoutRequest struct {
	SourceAccountID string `json:"sourceAccountId"`
	Destination     struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"destination"`
	Currency          string      `json:"currency"`
	Amount            json.Number `json:"amount"`
	Reference         string      `json:"reference"`
	ExternalReference string      `json:"externalReference"`
}

type PayoutResponse struct {
	ID                string `json:"id"`
	Status            string `json:"status"`
	CreatedDate       string `json:"createdDate"`
	ExternalReference string `json:"externalReference"`
	ApprovalStatus    string `json:"approvalStatus"`
	Message           string `json:"message"`
}

func (c *Client) InitiatePayout(ctx context.Context, payoutRequest *PayoutRequest) (*PayoutResponse, error) {
	f := connectors.ClientMetrics(ctx, "modulr", "initiate_payout")
	now := time.Now()
	defer f(ctx, now)

	body, err := json.Marshal(payoutRequest)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Post(c.buildEndpoint("payments"), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, unmarshalErrorWithoutRetry(resp.StatusCode, resp.Body).Error()
	}

	var res PayoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetPayout(ctx context.Context, payoutID string) (*PayoutResponse, error) {
	f := connectors.ClientMetrics(ctx, "modulr", "get_payout")
	now := time.Now()
	defer f(ctx, now)

	resp, err := c.httpClient.Get(c.buildEndpoint("payments?id=%s", payoutID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalErrorWithRetry(resp.StatusCode, resp.Body).Error()
	}

	var res PayoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
