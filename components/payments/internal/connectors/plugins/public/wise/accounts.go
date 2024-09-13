package wise

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/wise/client"
	"github.com/formancehq/payments/internal/models"
)

type accountsState struct {
	// Accounts are ordered by their ID
	LastAccountID uint64 `json:"lastAccountID"`
}

func (p Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState accountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	}

	var from client.Profile
	if req.FromPayload == nil {
		return models.FetchNextAccountsResponse{}, errors.New("missing from payload when fetching accounts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	newState := accountsState{
		LastAccountID: oldState.LastAccountID,
	}

	var accounts []models.PSPAccount
	hasMore := false
	// Wise balances are considered as accounts on our side.
	balances, err := p.client.GetBalances(ctx, from.ID)
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}

	for _, balance := range balances {
		if balance.ID <= oldState.LastAccountID {
			continue
		}

		raw, err := json.Marshal(balance)
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		accounts = append(accounts, models.PSPAccount{
			Reference:    strconv.FormatUint(balance.ID, 10),
			CreatedAt:    balance.CreationTime,
			Name:         &balance.Name,
			DefaultAsset: pointer.For(currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Amount.Currency)),
			Metadata: map[string]string{
				metadataProfileIDKey: strconv.FormatUint(from.ID, 10),
			},
			Raw: raw,
		})

		newState.LastAccountID = balance.ID

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
