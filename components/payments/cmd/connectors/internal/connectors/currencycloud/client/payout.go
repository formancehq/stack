package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type PayoutRequest struct {
	OnBehalfOf      string `json:"on_behalf_of"`
	BeneficiaryID   string `json:"beneficiary_id"`
	Currency        string `json:"currency"`
	Amount          string `json:"amount"`
	Reference       string `json:"reference"`
	UniqueRequestID string `json:"unique_request_id"`
}

func (pr *PayoutRequest) ToFormData() url.Values {
	form := url.Values{}
	form.Set("on_behalf_of", pr.OnBehalfOf)
	form.Set("beneficiary_id", pr.BeneficiaryID)
	form.Set("currency", pr.Currency)
	form.Set("amount", pr.Amount)
	form.Set("reference", pr.Reference)
	if pr.UniqueRequestID != "" {
		form.Set("unique_request_id", pr.UniqueRequestID)
	}

	return form
}

type PayoutResponse struct {
	ID               string `json:"id"`
	Amount           string `json:"amount"`
	BeneficiaryID    string `json:"beneficiary_id"`
	Currency         string `json:"currency"`
	Reference        string `json:"reference"`
	Status           string `json:"status"`
	Reason           string `json:"reason"`
	CreatorContactID string `json:"creator_contact_id"`
	PaymentType      string `json:"payment_type"`
	TransferredAt    string `json:"transferred_at"`
	PaymentDate      string `json:"payment_date"`
	FailureReason    string `json:"failure_reason"`
	CreatedAt        string `json:"created_at"`
	UpdatedAt        string `json:"updated_at"`
	UniqueRequestID  string `json:"unique_request_id"`
}

func (c *Client) InitiatePayout(ctx context.Context, payoutRequest *PayoutRequest) (*PayoutResponse, error) {
	form := payoutRequest.ToFormData()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("v2/payments/create"), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var payoutResponse PayoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&payoutResponse); err != nil {
		return nil, err
	}

	return &payoutResponse, nil
}

func (c *Client) GetPayout(ctx context.Context, payoutID string) (*PayoutResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.buildEndpoint("v2/payments/%s", payoutID), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var payoutResponse PayoutResponse
	if err = json.NewDecoder(resp.Body).Decode(&payoutResponse); err != nil {
		return nil, err
	}

	return &payoutResponse, nil
}
