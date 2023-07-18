package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type RecipientAccount struct {
	ID      uint64 `json:"id"`
	Profile uint64 `json:"profile"`
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
