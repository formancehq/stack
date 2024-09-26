package generic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type accountsState struct {
	LastCreatedAtFrom time.Time `json:"lastCreatedAtFrom"`
}

func (p Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	}

	newState := accountsState{
		LastCreatedAtFrom: oldState.LastCreatedAtFrom,
	}

	accounts := make([]models.PSPAccount, 0, req.PageSize)
	hasMore := false
	for page := 0; ; page++ {
		pagedAccounts, err := p.client.ListAccounts(ctx, int64(page), int64(req.PageSize), oldState.LastCreatedAtFrom)
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		for _, account := range pagedAccounts {
			switch account.CreatedAt.Compare(oldState.LastCreatedAtFrom) {
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
				Reference: account.Id,
				CreatedAt: account.CreatedAt,
				Name:      &account.AccountName,
				Metadata:  account.Metadata,
				Raw:       raw,
			})

			newState.LastCreatedAtFrom = account.CreatedAt

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
