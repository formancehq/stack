package stripe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

const (
	rootAccountReference = "root"
)

type AccountsState struct {
	LastID string `json:"lastID,omitempty"`
}

func (p *Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState AccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	}

	var accounts []models.PSPAccount
	newState := AccountsState{}
	rawAccounts, hasMore, err := p.client.GetAccounts(ctx, &oldState.LastID, PageLimit)
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}
	for _, acc := range rawAccounts {
		newState.LastID = acc.ID

		raw, err := json.Marshal(acc)
		if err != nil {
			return models.FetchNextAccountsResponse{}, err
		}

		metadata := make(map[string]string)
		for k, v := range acc.Metadata {
			metadata[k] = v
		}

		defaultAsset := currency.FormatAsset(supportedCurrenciesWithDecimal, string(acc.DefaultCurrency))
		accounts = append(accounts, models.PSPAccount{
			Reference:    acc.ID,
			CreatedAt:    time.Unix(acc.Created, 0).UTC(),
			DefaultAsset: &defaultAsset,
			Raw:          raw,
			Metadata:     metadata,
		})
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
