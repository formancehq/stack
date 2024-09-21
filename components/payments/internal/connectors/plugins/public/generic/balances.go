package generic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) fetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextBalancesResponse{}, errors.New("from payload is required")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	balances, err := p.client.GetBalances(ctx, from.Reference)
	if err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	var res []models.PSPBalance
	for _, balance := range balances.Balances {
		var amount big.Int
		_, ok := amount.SetString(balance.Amount, 10)
		if !ok {
			return models.FetchNextBalancesResponse{}, fmt.Errorf("failed to parse amount: %s", balance.Amount)
		}

		res = append(res, models.PSPBalance{
			AccountReference: balances.AccountID,
			CreatedAt:        balances.At,
			Amount:           &amount,
			Asset:            currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Currency),
		})
	}

	return models.FetchNextBalancesResponse{
		Balances: res,
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
