package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type TransferRequest struct {
	SourceAccountID      string  `json:"source_account_id"`
	DestinationAccountID string  `json:"destination_account_id"`
	Currency             string  `json:"currency"`
	Amount               float64 `json:"amount"`
	Reason               string  `json:"reason,omitempty"`
	UniqueRequestID      string  `json:"unique_request_id,omitempty"`
}

func (tr *TransferRequest) ToFormData() url.Values {
	form := url.Values{}
	form.Set("source_account_id", tr.SourceAccountID)
	form.Set("destination_account_id", tr.DestinationAccountID)
	form.Set("currency", tr.Currency)
	form.Set("amount", fmt.Sprintf("%v", tr.Amount))
	if tr.Reason != "" {
		form.Set("reason", tr.Reason)
	}
	if tr.UniqueRequestID != "" {
		form.Set("unique_request_id", tr.UniqueRequestID)
	}

	return form
}

type TransferResponse struct {
	ID                   string `json:"id"`
	ShortReference       string `json:"short_reference"`
	SourceAccountID      string `json:"source_account_id"`
	DestinationAccountID string `json:"destination_account_id"`
	Currency             string `json:"currency"`
	Amount               string `json:"amount"`
	Status               string `json:"status"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	CompletedAt          string `json:"completed_at"`
	CreatorAccountID     string `json:"creator_account_id"`
	CreatorContactID     string `json:"creator_contact_id"`
	Reason               string `json:"reason"`
	UniqueRequestID      string `json:"unique_request_id"`
}

func (c *Client) InitiateTransfer(ctx context.Context, transferRequest *TransferRequest) (*TransferResponse, error) {
	form := transferRequest.ToFormData()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		c.buildEndpoint("v2/transfers/create"), strings.NewReader(form.Encode()))
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

	var res TransferResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetTransfer(ctx context.Context, transferID string) (*TransferResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.buildEndpoint("v2/transfers/%s", transferID), http.NoBody)
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

	var res TransferResponse
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
