package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/public/mangopay/client"
	"github.com/formancehq/payments/internal/models"
)

type externalAccountsState struct {
	LastPage         int       `json:"last_page"`
	LastCreationDate time.Time `json:"last_creation_date"`
}

func (p Plugin) fetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var oldState externalAccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	} else {
		oldState = externalAccountsState{
			// Mangopay pages start at 1
			LastPage: 1,
		}
	}

	var from client.User
	if req.FromPayload == nil {
		return models.FetchNextExternalAccountsResponse{}, errors.New("missing from payload when fetching external accounts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	newState := externalAccountsState{
		LastPage:         oldState.LastPage,
		LastCreationDate: oldState.LastCreationDate,
	}

	var accounts []models.PSPAccount
	hasMore := false
	for page := oldState.LastPage; ; page++ {
		newState.LastPage = page

		pagedExternalAccounts, err := p.client.GetBankAccounts(ctx, from.ID, page, req.PageSize)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		if len(pagedExternalAccounts) == 0 {
			break
		}

		for _, bankAccount := range pagedExternalAccounts {
			creationDate := time.Unix(bankAccount.CreationDate, 0)
			switch creationDate.Compare(oldState.LastCreationDate) {
			case -1, 0:
				// creationDate <= state.LastCreationDate, nothing to do,
				// we already processed this bank account.
				continue
			default:
			}

			raw, err := json.Marshal(bankAccount)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference: bankAccount.ID,
				CreatedAt: creationDate,
				Name:      &bankAccount.OwnerName,
				Metadata: map[string]string{
					"user_id": from.ID,
				},
				Raw: raw,
			})

			newState.LastCreationDate = creationDate

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
