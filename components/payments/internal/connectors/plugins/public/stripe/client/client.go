package client

import (
	"context"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/account"
)

//go:generate mockgen -source client.go -destination client_generated.go -package client . Client
type Client interface {
	GetAccounts(ctx context.Context, lastID *string, pageSize int64) ([]*stripe.Account, bool, error)
}

type client struct {
	accountClient account.Client
}

func New(backend stripe.Backend, apiKey string) Client {
	if backend == nil {
		backend = stripe.GetBackend(stripe.APIBackend)
	}

	return &client{
		accountClient: account.Client{B: backend, Key: apiKey},
	}
}
