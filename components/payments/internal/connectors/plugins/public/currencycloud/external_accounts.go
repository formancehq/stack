package currencycloud

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type externalAccountsState struct {
	LastPage      int       `json:"lastPage"`
	LastCreatedAt time.Time `json:"lastCreatedAt"`
}

func (p Plugin) fetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var oldState externalAccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	} else {
		oldState = externalAccountsState{
			// Currencycloud pages start at 1
			LastPage: 1,
		}
	}

	newState := externalAccountsState{
		LastPage:      oldState.LastPage,
		LastCreatedAt: oldState.LastCreatedAt,
	}

	var accounts []models.PSPAccount
	hasMore := false
	page := oldState.LastPage
	for {
		pagedBeneficiarise, nextPage, err := p.client.GetBeneficiaries(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		if len(pagedBeneficiarise) == 0 {
			break
		}

		for _, beneficiary := range pagedBeneficiarise {
			switch beneficiary.CreatedAt.Compare(oldState.LastCreatedAt) {
			case -1, 0:
				// Account already ingested, skip
				continue
			default:
			}

			raw, err := json.Marshal(beneficiary)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference:    beneficiary.ID,
				CreatedAt:    beneficiary.CreatedAt,
				Name:         &beneficiary.Name,
				DefaultAsset: &beneficiary.Currency,
				Raw:          raw,
			})

			newState.LastCreatedAt = beneficiary.CreatedAt

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
		return models.FetchNextExternalAccountsResponse{}, err
	}

	return models.FetchNextExternalAccountsResponse{
		ExternalAccounts: accounts,
		NewState:         payload,
		HasMore:          hasMore,
	}, nil
}
