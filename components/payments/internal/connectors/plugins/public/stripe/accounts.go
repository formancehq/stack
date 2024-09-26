package stripe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

var (
	rootAccountReference = "root"
)

// root account reference is internal so we don't pass it to Stripe API clients
func resolveAccount(ref *string) *string {
	if *ref == rootAccountReference {
		return nil
	}
	return ref
}

type AccountsState struct {
	InitFinished bool   `json:"init_finished"`
	LastID       string `json:"lastID,omitempty"`
}

func (p *Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState AccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	}

	accounts := make([]models.PSPAccount, 0, req.PageSize)
	if !oldState.InitFinished {
		// create a root account if this is the first time this is being run
		accounts = append(accounts, models.PSPAccount{
			Name:      &rootAccountReference,
			Reference: rootAccountReference,
			CreatedAt: time.Now().UTC(),
			Raw:       json.RawMessage("{}"),
			Metadata:  map[string]string{},
		})
		oldState.InitFinished = true
	}

	needed := req.PageSize - len(accounts)

	newState := AccountsState{
		InitFinished: oldState.InitFinished,
	}
	rawAccounts, hasMore, err := p.client.GetAccounts(ctx, &oldState.LastID, int64(needed))
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
