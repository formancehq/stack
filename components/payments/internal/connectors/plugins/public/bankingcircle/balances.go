package bankingcircle

import (
	"context"
	"encoding/json"
	"errors"
	"math/big"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) fetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextBalancesResponse{}, errors.New("missing from payload when fetching balances")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	account, err := p.client.GetAccount(ctx, from.Reference)
	if err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	var balances []models.PSPBalance
	for _, balance := range account.Balances {
		// Note(polo): the last transaction timestamp is wrong in the banking
		// circle response. We will use the current time instead.
		// lastTransactionTimestamp, err := time.Parse("2006-01-02T15:04:05.999999999+00:00", balance.LastTransactionTimestamp)
		// if err != nil {
		// 	return models.FetchNextBalancesResponse{}, fmt.Errorf("failed to parse opening date: %w", err)
		// }
		lastTransactionTimestamp := time.Now().UTC()

		precision := supportedCurrenciesWithDecimal[balance.Currency]

		beginOfDayAmount, err := currency.GetAmountWithPrecisionFromString(balance.BeginOfDayAmount.String(), precision)
		if err != nil {
			return models.FetchNextBalancesResponse{}, err
		}

		intraDayAmount, err := currency.GetAmountWithPrecisionFromString(balance.IntraDayAmount.String(), precision)
		if err != nil {
			return models.FetchNextBalancesResponse{}, err
		}

		amount := big.NewInt(0).Add(beginOfDayAmount, intraDayAmount)

		balances = append(balances, models.PSPBalance{
			AccountReference: from.Reference,
			CreatedAt:        lastTransactionTimestamp,
			Amount:           amount,
			Asset:            currency.FormatAsset(supportedCurrenciesWithDecimal, balance.Currency),
		})
	}

	return models.FetchNextBalancesResponse{
		Balances: balances,
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
