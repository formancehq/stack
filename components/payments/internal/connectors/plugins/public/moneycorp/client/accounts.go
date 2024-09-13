package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type accountsResponse struct {
	Accounts []*Account `json:"data"`
}

type Account struct {
	ID         string `json:"id"`
	Attributes struct {
		AccountName string `json:"accountName"`
	} `json:"attributes"`
}

func (c *Client) GetAccounts(ctx context.Context, page int, pageSize int) ([]*Account, error) {
	// TODO(polo, crimson): metrics
	// metrics can also be embedded in wrapper
	// f := connectors.ClientMetrics(ctx, "moneycorp", "list_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/accounts", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	// TODO generic headers can be set in wrapper
	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("page[size]", strconv.Itoa(pageSize))
	q.Add("page[number]", fmt.Sprint(page))
	q.Add("sortBy", "id.asc")
	req.URL.RawQuery = q.Encode()

	accounts := accountsResponse{Accounts: make([]*Account, 0)}
	var errRes moneycorpError
	_, err = c.httpClient.Do(req, &accounts, &errRes)
	switch err {
	case nil:
		return accounts.Accounts, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get accounts: %w", err)
}
