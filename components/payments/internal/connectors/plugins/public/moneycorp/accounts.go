package moneycorp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type accountsState struct {
	LastPage int `json:"lastPage"`
	// Moneycorp does not send the creation date for accounts, but we can still
	// sort by ID created (which is incremental when creating accounts).
	LastIDCreated string `json:"lastIDCreated"`
}

func (p Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	}

	newState := accountsState{
		LastPage:      oldState.LastPage,
		LastIDCreated: oldState.LastIDCreated,
	}

	var accounts []models.PSPAccount
	hasMore := false
	for page := oldState.LastPage; ; page++ {
		newState.LastPage = page

		pagedAccounts, err := p.client.GetAccounts(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		for _, account := range pagedAccounts {
			if account.ID <= oldState.LastIDCreated {
				continue
			}

			raw, err := json.Marshal(account)
			if err != nil {
				return models.FetchNextAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference: account.ID,
				// Moneycorp does not send the opening date of the account
				CreatedAt: time.Now().UTC(),
				Name:      &account.Attributes.AccountName,
				Raw:       raw,
			})

			newState.LastIDCreated = account.ID

			if len(accounts) >= req.PageSize {
				break
			}
		}

		if len(pagedAccounts) < req.PageSize {
			break
		}

		if len(accounts) >= req.PageSize {
			hasMore = true
			break
		}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	return models.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: payload,
		HasMore:  hasMore,
	}, nil
}
