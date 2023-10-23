package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Account struct {
	AccountID          string `json:"accountId"`
	AccountDescription string `json:"accountDescription"`
	AccountIdentifiers []struct {
		Account              string `json:"account"`
		FinancialInstitution string `json:"financialInstitution"`
		Country              string `json:"country"`
	} `json:"accountIdentifiers"`
	Status           string `json:"status"`
	Currency         string `json:"currency"`
	OpeningDate      string `json:"openingDate"`
	ClosingDate      string `json:"closingDate"`
	OwnedByCompanyID string `json:"ownedByCompanyId"`
	ProtectionType   string `json:"protectionType"`
	Balances         []struct {
		Type                     string      `json:"type"`
		Currency                 string      `json:"currency"`
		BeginOfDayAmount         json.Number `json:"beginOfDayAmount"`
		FinancialDate            string      `json:"financialDate"`
		IntraDayAmount           json.Number `json:"intraDayAmount"`
		LastTransactionTimestamp string      `json:"lastTransactionTimestamp"`
	} `json:"balances"`
}

func (c *Client) GetAccounts(ctx context.Context, page int) ([]*Account, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/api/v1/accounts", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create account request: %w", err)
	}

	q := req.URL.Query()
	q.Add("PageSize", "100")
	q.Add("PageNumber", fmt.Sprint(page))

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

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

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read accounts response body: %w", err)
	}

	type response struct {
		Result   []*Account `json:"result"`
		PageInfo struct {
			CurrentPage int `json:"currentPage"`
			PageSize    int `json:"pageSize"`
		} `json:"pageInfo"`
	}

	var res response

	if err = json.Unmarshal(responseBody, &res); err != nil {
		return nil, fmt.Errorf("failed to unmarshal accounts response: %w", err)
	}

	return res.Result, nil
}
