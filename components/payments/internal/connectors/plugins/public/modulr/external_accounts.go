package modulr

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/models"
)

type externalAccountsState struct {
	LastCreatedAt time.Time `json:"lastCreatedAt"`
}

func (p Plugin) fetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var oldState externalAccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	}

	newState := externalAccountsState{
		LastCreatedAt: oldState.LastCreatedAt,
	}

	var accounts []models.PSPAccount
	hasMore := false
	for page := 0; ; page++ {
		pagedBeneficiarise, err := p.client.GetBeneficiaries(ctx, page, req.PageSize, oldState.LastCreatedAt)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		if len(pagedBeneficiarise.Content) == 0 {
			break
		}

		for _, beneficiary := range pagedBeneficiarise.Content {
			createdTime, err := time.Parse("2006-01-02T15:04:05.999-0700", beneficiary.Created)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, err
			}

			switch createdTime.Compare(oldState.LastCreatedAt) {
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
				Reference: beneficiary.ID,
				CreatedAt: createdTime,
				Name:      &beneficiary.Name,
				Raw:       raw,
			})

			newState.LastCreatedAt = createdTime

			if len(accounts) >= req.PageSize {
				break
			}
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
	}, err
}
