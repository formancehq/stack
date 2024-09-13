package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

//nolint:tagliatelle // allow different styled tags in client
type Transaction struct {
	ID        string    `json:"id"`
	AccountID string    `json:"account_id"`
	Currency  string    `json:"currency"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Action    string    `json:"action"`

	Amount json.Number `json:"amount"`
}

func (c *Client) GetTransactions(ctx context.Context, page int, pageSize int, updatedAtFrom time.Time) ([]Transaction, int, error) {
	if page < 1 {
		return nil, 0, fmt.Errorf("page must be greater than 0")
	}

	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "currencycloud", "list_transactions")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		c.buildEndpoint("v2/transactions/find"), http.NoBody)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", fmt.Sprint(page))
	q.Add("per_page", fmt.Sprint(pageSize))
	if !updatedAtFrom.IsZero() {
		q.Add("updated_at_from", updatedAtFrom.Format(time.DateOnly))
	}
	q.Add("order", "updated_at")
	q.Add("order_asc_desc", "asc")
	req.URL.RawQuery = q.Encode()

	req.Header.Add("Accept", "application/json")

	//nolint:tagliatelle // allow for client code
	type response struct {
		Transactions []Transaction `json:"transactions"`
		Pagination   struct {
			NextPage int `json:"next_page"`
		} `json:"pagination"`
	}

	res := response{Transactions: make([]Transaction, 0)}
	var errRes currencyCloudError
	_, err = c.httpClient.Do(req, &res, nil)
	switch err {
	case nil:
		return res.Transactions, res.Pagination.NextPage, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, 0, errRes.Error()
	}
	return nil, 0, fmt.Errorf("failed to get transactions %w", err)
}
