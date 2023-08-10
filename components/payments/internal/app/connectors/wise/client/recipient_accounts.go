package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RecipientAccount struct {
	ID         uint64 `json:"id"`
	Profile    uint64 `json:"profile"`
	Currency   string `json:"currency"`
	HolderName string `json:"accountHolderName"`
}

func (w *Client) GetRecipientAccounts(ctx context.Context, profileID uint64) ([]*RecipientAccount, error) {
	var recipientAccounts []*RecipientAccount

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, w.endpoint("v1/accounts"), http.NoBody)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("profile", fmt.Sprintf("%d", profileID))
	req.URL.RawQuery = q.Encode()

	res, err := w.httpClient.Do(req)
	if err != nil {
		return recipientAccounts, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		res.Body.Close()

		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if err = res.Body.Close(); err != nil {
		return nil, fmt.Errorf("failed to close response body: %w", err)
	}

	err = json.Unmarshal(body, &recipientAccounts)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal transfers: %w", err)
	}

	return recipientAccounts, nil
}

func (w *Client) GetRecipientAccount(ctx context.Context, accountID uint64) (*RecipientAccount, error) {
	if rc, ok := w.recipientAccountsCache.Get(accountID); ok {
		return rc, nil
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, w.endpoint(fmt.Sprintf("v1/accounts/%d", accountID)), http.NoBody)
	if err != nil {
		return nil, err
	}

	res, err := w.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		type errorResponse struct {
			Errors []struct {
				Code    string `json:"code"`
				Message string `json:"message"`
			}
		}

		var e errorResponse
		err = json.NewDecoder(res.Body).Decode(&e)
		if err != nil {
			return nil, fmt.Errorf("failed to decode error response: %w", err)
		}

		if len(e.Errors) == 0 {
			return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
		}

		switch e.Errors[0].Code {
		case "RECIPIENT_MISSING":
			// This is a valid response, we just don't have the account amoungs
			// our recipients.
			return &RecipientAccount{}, nil
		}

		return nil, fmt.Errorf("unexpected status code: %d with err: %v", res.StatusCode, e)
	}

	var account RecipientAccount
	err = json.NewDecoder(res.Body).Decode(&account)
	if err != nil {
		return nil, fmt.Errorf("failed to decode account: %w", err)
	}

	w.recipientAccountsCache.Add(accountID, &account)

	return &account, nil
}
