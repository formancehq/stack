package stripe

import (
	"context"
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

func (p *Plugin) fetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextBalancesResponse{}, errors.New("missing from payload when fetching balances")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	balance, err := p.client.GetAccountBalances(ctx, resolveAccount(from.Reference))
	if err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	var accountBalances []models.PSPBalance
	for _, available := range balance.Available {
		timestamp := time.Now()
		accountBalances = append(accountBalances, models.PSPBalance{
			AccountReference: from.Reference,
			Asset:            currency.FormatAsset(supportedCurrenciesWithDecimal, string(available.Currency)),
			Amount:           big.NewInt(available.Amount),
			CreatedAt:        timestamp,
		})
	}

	return models.FetchNextBalancesResponse{
		Balances: accountBalances,
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
