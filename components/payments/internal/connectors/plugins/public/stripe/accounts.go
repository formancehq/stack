package stripe

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/stripe/client"
	"github.com/formancehq/payments/internal/models"
)

var (
	rootAccountReference = "root"
)

// root account reference is internal so we don't pass it to Stripe API clients
func resolveAccount(ref string) string {
	if ref == rootAccountReference {
		return ""
	}
	return ref
}

type AccountsState struct {
	RootCreated bool            `json:"root_created"`
	Timeline    client.Timeline `json:"timeline"`
}

func (p *Plugin) fetchNextAccounts(ctx context.Context, req models.FetchNextAccountsRequest) (models.FetchNextAccountsResponse, error) {
	var oldState AccountsState
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextAccountsResponse{}, err
		}
	}

	accounts := make([]models.PSPAccount, 0, req.PageSize)
	if !oldState.RootCreated {
		// create a root account if this is the first time this is being run
		accounts = append(accounts, models.PSPAccount{
			Name:      &rootAccountReference,
			Reference: rootAccountReference,
			CreatedAt: time.Now().UTC(),
			Raw:       json.RawMessage("{}"),
			Metadata:  map[string]string{},
		})
		oldState.RootCreated = true
	}

	needed := req.PageSize - len(accounts)

	newState := oldState
	rawAccounts, timeline, hasMore, err := p.client.GetAccounts(ctx, oldState.Timeline, int64(needed))
	if err != nil {
		return models.FetchNextAccountsResponse{}, err
	}
	newState.Timeline = timeline

	for _, acc := range rawAccounts {
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
