package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
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

func (m *Client) GetAccounts(ctx context.Context, page, pageSize int, fromCreatedDate string) (*responseWrapper[[]*Account], error) {
	f := connectors.ClientMetrics(ctx, "modulr", "list_accounts")
	now := time.Now()
	defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, m.buildEndpoint("accounts"), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create accounts request: %w", err)
	}

	q := req.URL.Query()
	q.Add("page", strconv.Itoa(page))
	q.Add("size", strconv.Itoa(pageSize))
	q.Add("sortField", "createdDate")
	q.Add("sortOrder", "asc")
	if fromCreatedDate != "" {
		q.Add("fromCreatedDate", fromCreatedDate)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var res responseWrapper[[]*Account]
	if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return nil, err
	}

	return &res, nil
}
