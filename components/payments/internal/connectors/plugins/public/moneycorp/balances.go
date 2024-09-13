package moneycorp

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
)

func (p Plugin) fetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextBalancesResponse{}, errors.New("missing from payload when fetching balances")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	balances, err := p.client.GetAccountBalances(ctx, from.Reference)
	if err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	var accountBalances []models.PSPBalance
	for _, balance := range balances {
		precision, err := currency.GetPrecision(supportedCurrenciesWithDecimal, balance.Attributes.CurrencyCode)
		if err != nil {
			return models.FetchNextBalancesResponse{}, err
		}

		amount, err := currency.GetAmountWithPrecisionFromString(balance.Attributes.AvailableBalance.String(), precision)
		if err != nil {
			return models.FetchNextBalancesResponse{}, err
		}

		accountBalances = append(accountBalances, models.PSPBalance{
			AccountReference: from.Reference,
			CreatedAt:        time.Now(),
			Amount:           amount,
			Asset:            currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Attributes.CurrencyCode),
		})
	}

	return models.FetchNextBalancesResponse{
		Balances: accountBalances,
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
