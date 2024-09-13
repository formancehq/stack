package wise

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/wise/client"
	"github.com/formancehq/payments/internal/models"
)

type externalAccountsState struct {
	LastSeekPosition uint64 `json:"lastSeekPosition"`
}

func (p Plugin) fetchExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var oldState externalAccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	}

	var from client.Profile
	if req.FromPayload == nil {
		return models.FetchNextExternalAccountsResponse{}, errors.New("missing from payload when fetching external accounts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	newState := externalAccountsState{
		LastSeekPosition: oldState.LastSeekPosition,
	}

	var accounts []models.PSPAccount
	hasMore := false
	for {
		pagedExternalAccounts, err := p.client.GetRecipientAccounts(ctx, from.ID, req.PageSize, newState.LastSeekPosition)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		if len(pagedExternalAccounts.Content) == 0 {
			break
		}

		for _, externalAccount := range pagedExternalAccounts.Content {
			if externalAccount.ID <= oldState.LastSeekPosition {
				continue
			}

			raw, err := json.Marshal(externalAccount)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference:    strconv.FormatUint(externalAccount.ID, 10),
				CreatedAt:    time.Now().UTC(),
				Name:         &externalAccount.Name.FullName,
				DefaultAsset: pointer.For(currency.FormatAsset(supportedCurrenciesWithDecimal, externalAccount.Currency)),
				Raw:          raw,
			})

			if len(accounts) >= req.PageSize {
				break
			}
		}

		if len(accounts) >= req.PageSize {
			hasMore = true
			break
		}

		if pagedExternalAccounts.SeekPositionForNext == 0 {
			// No more data to fetch
			break
		}

		newState.LastSeekPosition = pagedExternalAccounts.SeekPositionForNext
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
