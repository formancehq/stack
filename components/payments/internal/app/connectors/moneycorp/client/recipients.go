package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type recipientsResponse struct {
	Recipients []*Recipient `json:"data"`
}

type Recipient struct {
	AccountID           string `json:"accountId"`
	BankAccountCurrency string `json:"bankAccountCurrency"`
	CreatedAt           string `json:"createdAt"`
	TemplateReference   string `json:"templateReference"`
}

func (c *Client) GetRecipients(ctx context.Context, page int, pageSize int) ([]*Recipient, error) {
	endpoint := fmt.Sprintf("%s/recipients", c.endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create recipients request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("pagesize", strconv.Itoa(pageSize))
	q.Add("pagenumber", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var recipients recipientsResponse
	if err := json.NewDecoder(resp.Body).Decode(&recipients); err != nil {
		return nil, fmt.Errorf("failed to unmarshal recipients response body: %w", err)
	}

	return recipients.Recipients, nil
}
