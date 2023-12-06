package client

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/stripe/stripe-go/v72"
)

type ClientOption interface {
	Apply(req *http.Request)
}
type ClientOptionFn func(req *http.Request)

func (fn ClientOptionFn) Apply(req *http.Request) {
	fn(req)
}

func QueryParam(key, value string) ClientOptionFn {
	return func(req *http.Request) {
		q := req.URL.Query()
		q.Set(key, value)
		req.URL.RawQuery = q.Encode()
	}
}

type Client interface {
	Accounts(ctx context.Context, options ...ClientOption) ([]*stripe.Account, bool, error)
	ExternalAccounts(ctx context.Context, options ...ClientOption) ([]*stripe.ExternalAccount, bool, error)
	BalanceTransactions(ctx context.Context, options ...ClientOption) ([]*stripe.BalanceTransaction, bool, error)
	Balance(ctx context.Context, options ...ClientOption) (*stripe.Balance, error)
	CreateTransfer(ctx context.Context, CreateTransferRequest *CreateTransferRequest, options ...ClientOption) (*stripe.Transfer, error)
	CreatePayout(ctx context.Context, createPayoutRequest *CreatePayoutRequest, options ...ClientOption) (*stripe.Payout, error)
	GetPayout(ctx context.Context, payoutID string, options ...ClientOption) (*stripe.Payout, error)
	ForAccount(account string) Client
}

type DefaultClient struct {
	httpClient    *http.Client
	apiKey        string
	stripeAccount string
}

func NewDefaultClient(apiKey string) *DefaultClient {
	return &DefaultClient{
		httpClient: newHTTPClient(),
		apiKey:     apiKey,
	}
}

func (d *DefaultClient) ForAccount(account string) Client {
	cp := *d
	cp.stripeAccount = account

	return &cp
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

var _ Client = &DefaultClient{}

func DatePtr(t time.Time) *time.Time {
	return &t
}
