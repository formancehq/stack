package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
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

	var res responseWrapper[[]Account]
	var errRes modulrError
	_, err = c.httpClient.Do(req, &res, &errRes)
	switch err {
	case nil:
		return &res, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get accounts: %w", err)
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

	var res Account
	var errRes modulrError
	_, err = c.httpClient.Do(req, &res, &errRes)
	switch err {
	case nil:
		return &res, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get account: %w", err)
}
