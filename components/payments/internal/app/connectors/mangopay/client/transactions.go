package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type payment struct {
	Id             string `json:"Id"`
	Tag            string `json:"Tag"`
	CreationDate   int64  `json:"CreationDate"`
	AuthorId       string `json:"AuthorId"`
	CreditedUserId string `json:"CreditedUserId"`
	DebitedFunds   struct {
		Currency string `json:"Currency"`
		Amount   int64  `json:"Amount"`
	} `json:"DebitedFunds"`
	CreditedFunds struct {
		Currency string `json:"Currency"`
		Amount   int64  `json:"Amount"`
	} `json:"CreditedFunds"`
	Fees struct {
		Currency string `json:"Currency"`
		Amount   int64  `json:"Amount"`
	} `json:"Fees"`
	Status           string `json:"Status"`
	ResultCode       string `json:"ResultCode"`
	ResultMessage    string `json:"ResultMessage"`
	ExecutionDate    int64  `json:"ExecutionDate"`
	Type             string `json:"Type"`
	Nature           string `json:"Nature"`
	CreditedWalletID string `json:"CreditedWalletId"`
	DebitedWalletID  string `json:"DebitedWalletId"`
}

func (c *Client) GetAllTransactions(ctx context.Context, userID string) ([]*payment, error) {
	var payments []*payment

	for page := 1; ; page++ {
		pagedPayments, err := c.getTransactions(ctx, userID, page)
		if err != nil {
			return nil, err
		}

		if len(pagedPayments) == 0 {
			break
		}

		payments = append(payments, pagedPayments...)
	}

	return payments, nil
}

func (c *Client) getTransactions(ctx context.Context, userID string, page int) ([]*payment, error) {
	endpoint := fmt.Sprintf("%s/v2.01/%s/users/%s/transactions", c.endpoint, c.clientID, userID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create login request: %w", err)
	}

	q := req.URL.Query()
	q.Add("per_page", "100")
	q.Add("page", fmt.Sprint(page))
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			c.logger.Error(err)
		}
	}()

	var payments []*payment
	if err := json.NewDecoder(resp.Body).Decode(&payments); err != nil {
		return nil, fmt.Errorf("failed to unmarshal login response body: %w", err)
	}

	return payments, nil
}
