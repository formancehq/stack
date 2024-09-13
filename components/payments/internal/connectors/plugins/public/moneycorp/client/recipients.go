package client

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type recipientsResponse struct {
	Recipients []*Recipient `json:"data"`
}

type Recipient struct {
	ID         string `json:"id"`
	Attributes struct {
		BankAccountCurrency string `json:"bankAccountCurrency"`
		CreatedAt           string `json:"createdAt"`
		BankAccountName     string `json:"bankAccountName"`
	} `json:"attributes"`
}

func (c *Client) GetRecipients(ctx context.Context, accountID string, page int, pageSize int) ([]*Recipient, error) {
	// TODO(polo): add metrics
	// f := connectors.ClientMetrics(ctx, "moneycorp", "list_recipients")
	// now := time.Now()
	// defer f(ctx, now)

	endpoint := fmt.Sprintf("%s/accounts/%s/recipients", c.endpoint, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipients request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("page[size]", strconv.Itoa(pageSize))
	q.Add("page[number]", fmt.Sprint(page))
	q.Add("sortBy", "createdAt.asc")
	req.URL.RawQuery = q.Encode()

	recipients := recipientsResponse{Recipients: make([]*Recipient, 0)}
	var errRes moneycorpError
	_, err = c.httpClient.Do(req, &recipients, &errRes)
	switch err {
	case nil:
		return recipients.Recipients, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error()
	}
	return nil, fmt.Errorf("failed to get recipients %w", err)
}
