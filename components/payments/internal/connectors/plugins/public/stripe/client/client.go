package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/account"
	"github.com/stripe/stripe-go/v79/balance"
	"github.com/stripe/stripe-go/v79/balancetransaction"
	"github.com/stripe/stripe-go/v79/bankaccount"
)

//go:generate mockgen -source client.go -destination client_generated.go -package client . Client
type Client interface {
	GetAccounts(ctx context.Context, timeline Timeline, pageSize int64) ([]*stripe.Account, Timeline, bool, error)
	GetAccountBalances(ctx context.Context, accountID string) (*stripe.Balance, error)
	GetExternalAccounts(ctx context.Context, accountID string, timeline Timeline, pageSize int64) ([]*stripe.BankAccount, Timeline, bool, error)
	GetPayments(ctx context.Context, accountID string, timeline Timeline, pageSize int64) ([]*stripe.BalanceTransaction, Timeline, bool, error)
}

type client struct {
	accountClient            account.Client
	balanceClient            balance.Client
	bankAccountClient        bankaccount.Client
	balanceTransactionClient balancetransaction.Client
}

func New(backend stripe.Backend, apiKey string) Client {
	if backend == nil {
		backend = stripe.GetBackend(stripe.APIBackend)
	}

	return &client{
		accountClient:            account.Client{B: backend, Key: apiKey},
		balanceClient:            balance.Client{B: backend, Key: apiKey},
		bankAccountClient:        bankaccount.Client{B: backend, Key: apiKey},
		balanceTransactionClient: balancetransaction.Client{B: backend, Key: apiKey},
	}
}

func limit(wanted int64, have int) *int64 {
	needed := wanted - int64(have)
	return &needed
}
