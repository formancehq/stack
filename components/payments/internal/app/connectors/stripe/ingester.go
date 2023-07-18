package stripe

import (
	"context"

	"github.com/stripe/stripe-go/v72"
)

type ingestTransaction func(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error
type ingestAccounts func(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error

type Ingester interface {
	IngestTransactions(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error
	IngestAccounts(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error
}

type ingester struct {
	it ingestTransaction
	ia ingestAccounts
}

func NewIngester(it ingestTransaction, ia ingestAccounts) Ingester {
	return &ingester{
		it: it,
		ia: ia,
	}
}

func (i *ingester) IngestTransactions(ctx context.Context, batch []*stripe.BalanceTransaction, commitState TimelineState, tail bool) error {
	return i.it(ctx, batch, commitState, tail)
}

func (i *ingester) IngestAccounts(ctx context.Context, batch []*stripe.Account, commitState TimelineState, tail bool) error {
	return i.ia(ctx, batch, commitState, tail)
}
