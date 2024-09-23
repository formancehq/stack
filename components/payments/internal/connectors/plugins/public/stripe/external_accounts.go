package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

func (p *Plugin) fetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var (
		oldState AccountsState
		from     models.PSPAccount
	)
	if req.State != nil {
		if err := json.Unmarshal(req.State, &oldState); err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}
	}
	if req.FromPayload == nil {
		return models.FetchNextExternalAccountsResponse{}, errors.New("missing from payload when fetching external accounts")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}

	newState := AccountsState{}
	var accounts []models.PSPAccount

	rawAccounts, hasMore, err := p.client.GetExternalAccounts(ctx, &from.Reference, &oldState.LastID, PageLimit)
	if err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}
	for _, acc := range rawAccounts {
		newState.LastID = acc.ID

		raw, err := json.Marshal(acc)
		if err != nil {
			return models.FetchNextExternalAccountsResponse{}, err
		}

		metadata := make(map[string]string)
		for k, v := range acc.Metadata {
			metadata[k] = v
		}

		defaultAsset := currency.FormatAsset(supportedCurrenciesWithDecimal, string(acc.Currency))

		if acc.Account == nil {
			return models.FetchNextExternalAccountsResponse{}, fmt.Errorf("internal account %q is missing from response for %q", from.Reference, acc.ID)
		}

		accounts = append(accounts, models.PSPAccount{
			Reference:    acc.ID,
			CreatedAt:    time.Unix(acc.Account.Created, 0).UTC(),
			DefaultAsset: &defaultAsset,
			Raw:          raw,
			Metadata:     metadata,
		})
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
