package adyen

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins"
	"github.com/formancehq/payments/internal/models"
)

type accountsState struct {
	LastPage int `json:"lastPage"`

	// Adyen API sort the accounts by ID which is the same as the name
	// and we cannot sort by other things. It means that when we fetched
	// everything, we will need to return an empty state in order to
	// refetch everything at the next polling iteration...
	// It should not change anything in the database, but it will generate
	// duplicates in events, but with the same IdempotencyKey.
	LastID string `json:"lastId"`
}

func (p Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	} else {
		oldState = accountsState{
			LastPage: 1,
		}
	}

	newState := accountsState{
		LastPage: oldState.LastPage,
		LastID:   oldState.LastID,
	}

	var accounts []models.PSPAccount
	hasMore := false
	page := oldState.LastPage
	for {
		pagedAccount, err := p.client.GetMerchantAccounts(ctx, int32(page), int32(req.PageSize))
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		if len(pagedAccount) == 0 {
			break
		}

		for _, account := range pagedAccount {
			raw, err := json.Marshal(account)
			if err != nil {
				return models.FetchNextAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference: *account.Id,
				CreatedAt: time.Now().UTC(),
				Name:      account.Name,
				Raw:       raw,
			})

			newState.LastID = *account.Id

			if len(accounts) >= req.PageSize {
				break
			}
		}

		if len(pagedAccount) < req.PageSize {
			break
		}

		if len(accounts) >= req.PageSize {
			hasMore = true
			break
		}

		page++
	}

	newState.LastPage = page

	if !hasMore {
		// Since the merchant accounts sorting is done by ID, if a new one is
		// created with and ID lower than the last one we fetched, we will not
		// fetch it. So we need to reset the state to fetch everything again
		// when we have fetched eveything.
		// It will not create duplicates inside the database since we're based
		// on the ID of the account, but it will create duplicates in the events
		// but with the same IdempotencyKey, so should be fine.
		newState = accountsState{}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	return models.FetchNextAccountsResponse{
		Accounts: accounts,
		NewState: payload,
		HasMore:  hasMore,
	}, plugins.ErrNotImplemented
}
