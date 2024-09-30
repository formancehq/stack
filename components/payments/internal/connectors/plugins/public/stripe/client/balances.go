package client

import (
	"context"
	"fmt"

	"github.com/stripe/stripe-go/v79"
)

func (c *client) GetAccountBalances(ctx context.Context, accountID string) (*stripe.Balance, error) {
	// TODO
	//	f := connectors.ClientMetrics(ctx, "stripe", "get_balances")
	//	now := time.Now()
	//	defer f(ctx, now)
	var filters stripe.Params
	if accountID != "" {
		filters.StripeAccount = &accountID
	}

	balance, err := c.balanceClient.Get(&stripe.BalanceParams{Params: filters})
	if err != nil {
		return nil, fmt.Errorf("failed to get stripe balance: %w", err)
	}
	return balance, nil
}
