package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/formancehq/payments/internal/connectors/httpwrapper"
)

type RecipientAccountsResponse struct {
	Content                []*RecipientAccount `json:"content"`
	SeekPositionForCurrent uint64              `json:"seekPositionForCurrent"`
	SeekPositionForNext    uint64              `json:"seekPositionForNext"`
	Size                   int                 `json:"size"`
}

type RecipientAccount struct {
	ID       uint64 `json:"id"`
	Profile  uint64 `json:"profileId"`
	Currency string `json:"currency"`
	Name     struct {
		FullName string `json:"fullName"`
	} `json:"name"`
}

func (c *Client) GetRecipientAccounts(ctx context.Context, profileID uint64, pageSize int, seekPositionForNext uint64) (*RecipientAccountsResponse, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "wise", "list_recipient_accounts")
	// now := time.Now()
	// defer f(ctx, now)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, c.endpoint("v2/accounts"), http.NoBody)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("profile", fmt.Sprintf("%d", profileID))
	q.Add("size", fmt.Sprintf("%d", pageSize))
	q.Add("sort", "id,asc")
	if seekPositionForNext > 0 {
		q.Add("seekPosition", fmt.Sprintf("%d", seekPositionForNext))
	}
	req.URL.RawQuery = q.Encode()

	var accounts RecipientAccountsResponse
	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &accounts, &errRes)
	switch err {
	case nil:
		return &accounts, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		return nil, errRes.Error(statusCode).Error()
	}
	return nil, fmt.Errorf("failed to get recipient accounts: %w", err)
}

func (c *Client) GetRecipientAccount(ctx context.Context, accountID uint64) (*RecipientAccount, error) {
	// TODO(polo): metrics
	// f := connectors.ClientMetrics(ctx, "wise", "get_recipient_account")
	// now := time.Now()
	// defer f(ctx, now)

	if rc, ok := c.recipientAccountsCache.Get(accountID); ok {
		return rc, nil
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet, c.endpoint(fmt.Sprintf("v1/accounts/%d", accountID)), http.NoBody)
	if err != nil {
		return nil, err
	}

	var res RecipientAccount
	var errRes wiseErrors
	statusCode, err := c.httpClient.Do(req, &res, &errRes)
	switch err {
	case nil:
		c.recipientAccountsCache.Add(accountID, &res)
		return &res, nil
	case httpwrapper.ErrStatusCodeUnexpected:
		// TODO(polo): retryable errors
		e := errRes.Error(statusCode)
		if e.Code == "RECIPIENT_MISSING" {
			// This is a valid response, we just don't have the account amoungs
			// our recipients.
			return &RecipientAccount{}, nil
		}
		return nil, e.Error()
	}
	return nil, fmt.Errorf("failed to get recipient account: %w", err)
}
