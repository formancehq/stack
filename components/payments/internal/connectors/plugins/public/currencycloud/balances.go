package currencycloud

import (
	"context"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) fetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	page := 1
	balances := make([]models.PSPBalance, 0)
	for {
		if page < 0 {
			break
		}

		pagedBalances, nextPage, err := p.client.GetBalances(ctx, page, req.PageSize)
		if err != nil {
			return models.FetchNextBalancesResponse{}, err
		}

		page = nextPage

		for _, balance := range pagedBalances {
			precision, ok := supportedCurrenciesWithDecimal[balance.Currency]
			if !ok {
				return models.FetchNextBalancesResponse{}, nil
			}

			amount, err := currency.GetAmountWithPrecisionFromString(balance.Amount.String(), precision)
			if err != nil {
				return models.FetchNextBalancesResponse{}, err
			}

			balances = append(balances, models.PSPBalance{
				AccountReference: balance.AccountID,
				CreatedAt:        balance.UpdatedAt,
				Amount:           amount,
				Asset:            currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Currency),
			})
		}
	}

	return models.FetchNextBalancesResponse{
		Balances: balances,
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
