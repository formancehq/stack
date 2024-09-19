package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
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

func (c *Client) GetAccounts(ctx context.Context, page int, pageSize int, fromOpeningDate time.Time) ([]Account, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "bankingcircle", "list_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.endpoint+"/api/v1/accounts", http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create account request: %w", err)
	}

	q := req.URL.Query()
	q.Add("PageSize", fmt.Sprint(pageSize))
	q.Add("PageNumber", fmt.Sprint(page))
	if !fromOpeningDate.IsZero() {
		q.Add("OpeningDateFrom", fromOpeningDate.Format(time.DateOnly))
	}
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	type response struct {
		Result   []Account `json:"result"`
		PageInfo struct {
			CurrentPage int `json:"currentPage"`
			PageSize    int `json:"pageSize"`
		} `json:"pageInfo"`
	}

	res := response{Result: make([]Account, 0)}
	statusCode, err := c.httpClient.Do(req, &res, nil)
	switch err {
	case nil:
		return res.Result, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, fmt.Errorf("received status code %d for get accounts", statusCode)
	}
	return nil, fmt.Errorf("failed to get accounts: %w", err)
}

func (c *Client) GetAccount(ctx context.Context, accountID string) (*Account, error) {
	if err := c.ensureAccessTokenIsValid(ctx); err != nil {
		return nil, err
	}

	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "bankingcircle", "get_account")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/v1/accounts/%s", c.endpoint, accountID), http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create account request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+c.accessToken)

	var account Account
	statusCode, err := c.httpClient.Do(req, &account, nil)
	switch err {
	case nil:
		return &account, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, fmt.Errorf("received status code %d for get account", statusCode)
	}
	return nil, fmt.Errorf("failed to get account: %w", err)
}
