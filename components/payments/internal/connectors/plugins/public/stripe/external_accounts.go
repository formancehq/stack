package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/connectors/plugins/public/stripe/client"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

type ExternalAccountsState struct {
	Timeline client.Timeline `json:"timeline"`
}

func (p *Plugin) fetchNextExternalAccounts(ctx context.Context, req models.FetchNextExternalAccountsRequest) (models.FetchNextExternalAccountsResponse, error) {
	var (
		oldState ExternalAccountsState
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

	newState := oldState
	var accounts []models.PSPAccount

	rawAccounts, timeline, hasMore, err := p.client.GetExternalAccounts(
		ctx,
		resolveAccount(from.Reference),
		oldState.Timeline,
		int64(req.PageSize),
	)
	if err != nil {
		return models.FetchNextExternalAccountsResponse{}, err
	}
	newState.Timeline = timeline

	for _, acc := range rawAccounts {
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
