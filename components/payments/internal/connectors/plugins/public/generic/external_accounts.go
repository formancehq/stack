package generic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type externalAccountsState struct {
	LastCreatedAtFrom time.Time `json:"lastCreatedAtFrom"`
}

func (p Plugin) fetchExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	}

	newState := externalAccountsState{
		LastCreatedAtFrom: oldState.LastCreatedAtFrom,
	}

	accounts := make([]models.PSPAccount, 0, req.PageSize)
	hasMore := false
	for page := 0; ; page++ {
		pagedExternalAccounts, err := p.client.ListBeneficiaries(ctx, int64(page), int64(req.PageSize), oldState.LastCreatedAtFrom)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		if len(pagedExternalAccounts) == 0 {
			break
		}

		for _, account := range pagedExternalAccounts {
			switch account.CreatedAt.Compare(oldState.LastCreatedAtFrom) {
			case -1, 0:
				// Account already ingested, skip
				continue
			default:
			}

			raw, err := json.Marshal(account)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference: account.Id,
				CreatedAt: account.CreatedAt,
				Name:      &account.OwnerName,
				Metadata:  account.Metadata,
				Raw:       raw,
			})

			newState.LastCreatedAtFrom = account.CreatedAt

			if len(accounts) >= req.PageSize {
				break
			}
		}

		if len(pagedExternalAccounts) < req.PageSize {
			break
		}

		if len(accounts) >= req.PageSize {
			hasMore = true
			break
		}
	}

	payload, err := json.Marshal(newState)
	if err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	return models.FetchNextExternalAccountsResponse{
		ExternalAccounts: accounts,
		NewState:         payload,
		HasMore:          hasMore,
	}, nil
}
