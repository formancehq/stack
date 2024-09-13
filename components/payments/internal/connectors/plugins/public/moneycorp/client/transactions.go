package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type transactionsResponse struct {
	Transactions []*Transaction `json:"data"`
}

type fetchTransactionRequest struct {
	Data struct {
		Attributes struct {
			TransactionDateTimeFrom string `json:"transactionDateTimeFrom"`
		} `json:"attributes"`
	} `json:"data"`
}

type Transaction struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		AccountID            int32       `json:"accountId"`
		CreatedAt            string      `json:"createdAt"`
		Currency             string      `json:"transactionCurrency"`
		Amount               json.Number `json:"transactionAmount"`
		Direction            string      `json:"transactionDirection"`
		Type                 string      `json:"transactionType"`
		ClientReference      string      `json:"clientReference"`
		TransactionReference string      `json:"transactionReference"`
	} `json:"attributes"`
}

func (c *Client) GetTransactions(ctx context.Context, accountID string, page, pageSize int, lastCreatedAt time.Time) ([]*Transaction, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "moneycorp", "list_transactions")
	// now := time.Now()
	// defer f(ctx, now)

	var body io.Reader
	if !lastCreatedAt.IsZero() {
		reqBody := fetchTransactionRequest{
			Data: struct {
				Attributes struct {
					TransactionDateTimeFrom string "json:\"transactionDateTimeFrom\""
				} "json:\"attributes\""
			}{
				Attributes: struct {
					TransactionDateTimeFrom string "json:\"transactionDateTimeFrom\""
				}{
					TransactionDateTimeFrom: lastCreatedAt.Format("2006-01-02T15:04:05.999999999"),
				},
			},
		}

		raw, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal transfer request: %w", err)
		}

		body = bytes.NewBuffer(raw)
	} else {
		body = http.NoBody
	}

	endpoint := fmt.Sprintf("%s/accounts/%s/transactions/find", c.endpoint, accountID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create transactions request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	q := req.URL.Query()
	q.Add("page[size]", strconv.Itoa(pageSize))
	q.Add("page[number]", fmt.Sprint(page))
	q.Add("sortBy", "createdAt.asc")
	req.URL.RawQuery = q.Encode()

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get transactions: %w", err)
	}

	defer func() {
		err = resp.Body.Close()
		if err != nil {
			// TODO(polo): log error
			// c.logger.Error(err)
			_ = err
		}
	}()

	if resp.StatusCode == http.StatusNotFound {
		return []*Transaction{}, nil
	}

	if resp.StatusCode != http.StatusOK {
		// TODO(polo): retryable errors
		return nil, unmarshalError(resp.StatusCode, resp.Body).Error()
	}

	var transactions transactionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal transactions response body: %w", err)
	}

	return transactions.Transactions, nil
}
