package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//nolint:tagliatelle // allow for clients
type Account struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Balance     string `json:"balance"`
	Currency    string `json:"currency"`
	CustomerID  string `json:"customerId"`
	Identifiers []struct {
		AccountNumber string `json:"accountNumber"`
		SortCode      string `json:"sortCode"`
		Type          string `json:"type"`
	} `json:"identifiers"`
	DirectDebit bool   `json:"directDebit"`
	CreatedDate string `json:"createdDate"`
}

func (c *Client) GetAccounts(ctx context.Context, page, pageSize int, fromCreatedAt time.Time) (*responseWrapper[[]Account], error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "modulr", "list_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildEndpoint("accounts"), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	q.Add("sortField", "createdDate")
	q.Add("sortOrder", "asc")
	if !fromCreatedAt.IsZero() {
		q.Add("fromCreatedDate", fromCreatedAt.Format("2006-01-02T15:04:05-0700"))
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO(polo): retryable errors
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res responseWrapper[[]Account]
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "modulr", "list_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.buildEndpoint("accounts/%s", accountID), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// TODO(polo): retryable errors
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res Account
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
