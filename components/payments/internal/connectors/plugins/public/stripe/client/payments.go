package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
)

var (
	expandSource                    = "data.source"
	expandSourceCharge              = "data.source.charge"
	expandSourceDispute             = "data.source.dispute"
	expandSourcePayout              = "data.source.payout"
	expandSourceRefund              = "data.source.refund"
	expandSourceTransfer            = "data.source.transfer"
	expandSourcePaymentIntent       = "data.source.payment_intent"
	expandSourceRefundPaymentIntent = "data.source.refund.payment_intent"
)

func (c *client) GetPayments(ctx context.Context, accountID *string, lastID *string, pageSize int64) ([]*stripe.BalanceTransaction, bool, error) {
	// TODO
	//	f := connectors.ClientMetrics(ctx, "stripe", "list_accounts")
	//	now := time.Now()
	//	defer f(ctx, now)
	filters := stripe.ListParams{
		Limit: &pageSize,
	}
	if lastID == nil {
		filters.StartingAfter = lastID
	}

	expand := []*string{
		&expandSource,
		&expandSourceCharge,
		&expandSourceDispute,
		&expandSourcePayout,
		&expandSourceRefund,
		&expandSourceTransfer,
		&expandSourcePaymentIntent,
		&expandSourceRefundPaymentIntent,
	}

	itr := c.balanceTransactionClient.List(&stripe.BalanceTransactionListParams{
		ListParams: filters,
		Expand:     expand,
	})
	if err := itr.Err(); err != nil {
		return nil, false, handleError(itr.Iter)
	}
	return itr.BalanceTransactionList().Data, itr.BalanceTransactionList().ListMeta.HasMore, nil
}
