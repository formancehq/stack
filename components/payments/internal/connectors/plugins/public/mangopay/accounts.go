package mangopay

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/mangopay/client"
	"github.com/formancehq/payments/internal/models"
)

type accountsState struct {
	LastPage         int       `json:"lastPage"`
	LastCreationDate time.Time `json:"lastCreationDate"`
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

	var from client.User
	if req.FromPayload == nil {
		return models.FetchNextAccountsResponse{}, errors.New("missing from payload when fetching accounts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	newState := accountsState{
		LastPage:         oldState.LastPage,
		LastCreationDate: oldState.LastCreationDate,
	}

	var accounts []models.PSPAccount
	hasMore := false
	page := oldState.LastPage
	for {
		pagedAccounts, err := p.client.GetWallets(ctx, from.ID, page, req.PageSize)
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		if len(pagedAccounts) == 0 {
			break
		}

		for _, account := range pagedAccounts {
			accountCreationDate := time.Unix(account.CreationDate, 0)
			switch accountCreationDate.Compare(oldState.LastCreationDate) {
			case -1, 0:
				// creationDate <= state.LastCreationDate, nothing to do,
				// we already processed this account.
				continue
			default:
			}

			raw, err := json.Marshal(account)
			if err != nil {
				return models.FetchNextAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference:    account.ID,
				CreatedAt:    accountCreationDate,
				Name:         &account.Description,
				DefaultAsset: pointer.For(currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency)),
				Metadata: map[string]string{
					"user_id": from.ID,
				},
				Raw: raw,
			})

			newState.LastCreationDate = accountCreationDate

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

		page++
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
