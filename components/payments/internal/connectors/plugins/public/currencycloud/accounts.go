package currencycloud

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type accountsState struct {
	LastPage      int       `json:"lastPage"`
	LastCreatedAt time.Time `json:"lastCreatedAt"`
}

func (p Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	} else {
		oldState = accountsState{
			// Currencycloud pages start at 1
			LastPage: 1,
		}
	}

	newState := accountsState{
		LastPage:      oldState.LastPage,
		LastCreatedAt: oldState.LastCreatedAt,
	}

	var accounts []models.PSPAccount
	hasMore := false
	page := oldState.LastPage
	for {
		pagedAccounts, nextPage, err := p.client.GetAccounts(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		for _, account := range pagedAccounts {
			switch account.CreatedAt.Compare(oldState.LastCreatedAt) {
			case -1, 0:
				// Account already ingested, skip
				continue
			default:
			}

			raw, err := json.Marshal(account)
			if err != nil {
				return models.FetchNextAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference: account.ID,
				CreatedAt: account.CreatedAt,
				Name:      &account.AccountName,
				Raw:       raw,
			})

			newState.LastCreatedAt = account.CreatedAt

			if len(accounts) >= req.PageSize {
				break
			}
		}

		if len(accounts) >= req.PageSize {
			hasMore = true
			break
		}

		if nextPage == -1 {
			break
		}

		page = nextPage
	}

	newState.LastPage = page

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
