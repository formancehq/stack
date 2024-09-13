package modulr

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/connectors/plugins/currency"
	"github.com/formancehq/payments/internal/models"
)

func (p Plugin) fetchNextBalances(ctx context.Context, req models.FetchNextBalancesRequest) (models.FetchNextBalancesResponse, error) {
	var from models.PSPAccount
	if req.FromPayload == nil {
		return models.FetchNextBalancesResponse{}, errors.New("missing from payload when fetching payments")
	}
	if err := json.Unmarshal(req.FromPayload, &from); err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	account, err := p.client.GetAccount(ctx, from.Reference)
	if err != nil {
		return models.FetchNextBalancesResponse{}, err
	}

	precision := supportedCurrenciesWithDecimal[account.Currency]

	amount, err := currency.GetAmountWithPrecisionFromString(account.Balance, precision)
	if err != nil {
		return models.FetchNextBalancesResponse{}, fmt.Errorf("failed to parse amount %s: %w", account.Balance, err)
	}

	balance := models.PSPBalance{
		AccountReference: from.Reference,
		CreatedAt:        time.Now().UTC(),
		Amount:           amount,
		Asset:            currency.FormatAsset(supportedCurrenciesWithDecimal, account.Currency),
	}

	return models.FetchNextBalancesResponse{
		Balances: []models.PSPBalance{balance},
		NewState: []byte{},
		HasMore:  false,
	}, nil
}
