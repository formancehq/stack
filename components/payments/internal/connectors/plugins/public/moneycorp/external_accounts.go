package moneycorp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type externalAccountsState struct {
	LastPage int `json:"last_page"`
	// Moneycorp does not allow us to sort by , but we can still
	// sort by ID created (which is incremental when creating accounts).
	LastCreatedAt time.Time `json:"last_created_at"`
}

func (p Plugin) fetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var oldState externalAccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	}

	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextExternalAccountsResponse{}, errors.New("missing from payload when fetching external accounts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	newState := externalAccountsState{
		LastPage:      oldState.LastPage,
		LastCreatedAt: oldState.LastCreatedAt,
	}

	var accounts []models.PSPAccount
	hasMore := false
	for page := oldState.LastPage; ; page++ {
		newState.LastPage = page

		pagedRecipients, err := p.client.GetRecipients(ctx, from.Reference, page, req.PageSize)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		if len(pagedRecipients) == 0 {
			break
		}

		for _, recipient := range pagedRecipients {
			createdAt, err := time.Parse("2006-01-02T15:04:05.999999999", recipient.Attributes.CreatedAt)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, fmt.Errorf("failed to parse transaction date: %v", err)
			}

			switch createdAt.Compare(oldState.LastCreatedAt) {
			case -1, 0:
				continue
			default:
			}

			raw, err := json.Marshal(recipient)
			if err != nil {
				return models.FetchNextExternalAccountsResponse{}, err
			}

			accounts = append(accounts, models.PSPAccount{
				Reference: recipient.ID,
				// Moneycorp does not send the opening date of the account
				CreatedAt:    createdAt,
				Name:         &recipient.Attributes.BankAccountName,
				DefaultAsset: pointer.For(currency.FormatAsset(supportedCurrenciesWithDecimal, recipient.Attributes.BankAccountCurrency)),
				Raw:          raw,
			})

			newState.LastCreatedAt = createdAt

			if len(accounts) >= req.PageSize {
				break
			}
		}

		if len(accounts) >= req.PageSize {
			hasMore = true
			break
		}

		if len(pagedRecipients) < req.PageSize {
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
