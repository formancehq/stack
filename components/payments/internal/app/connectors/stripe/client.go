package stripe

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/pkg/errors"
	"github.com/stripe/stripe-go/v72"
)

const (
	balanceTransactionsEndpoint = "https://api.stripe.com/v1/balance_transactions"
	accountsEndpoint            = "https://api.stripe.com/v1/accounts"
	balanceEndpoint             = "https://api.stripe.com/v1/balance"
)

type ClientOption interface {
	apply(req *http.Request)
}
type ClientOptionFn func(req *http.Request)

func (fn ClientOptionFn) apply(req *http.Request) {
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
	BalanceTransactions(ctx context.Context, options ...ClientOption) ([]*stripe.BalanceTransaction, bool, error)
	Balance(ctx context.Context, options ...ClientOption) (*stripe.Balance, error)
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

func (d *DefaultClient) BalanceTransactions(ctx context.Context,
	options ...ClientOption,
) ([]*stripe.BalanceTransaction, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, balanceTransactionsEndpoint, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.apply(req)
	}

	if d.stripeAccount != "" {
		req.Header.Set("Stripe-Account", d.stripeAccount)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.apiKey, "") // gfyrag: really weird authentication right?

	var httpResponse *http.Response

	httpResponse, err = d.httpClient.Do(req)
	if err != nil {
		return nil, false, errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type listResponse struct {
		TransactionsListResponse
		Data []json.RawMessage `json:"data"`
	}

	rsp := &listResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, false, errors.Wrap(err, "decoding response")
	}

	asBalanceTransactions := make([]*stripe.BalanceTransaction, 0)

	if len(rsp.Data) > 0 {
		for _, data := range rsp.Data {
			asBalanceTransaction := &stripe.BalanceTransaction{}

			err = json.Unmarshal(data, &asBalanceTransaction)
			if err != nil {
				return nil, false, err
			}

			asBalanceTransactions = append(asBalanceTransactions, asBalanceTransaction)
		}
	}

	return asBalanceTransactions, rsp.HasMore, nil
}

type BalanceResponse struct {
	*stripe.Balance
}

func (d *DefaultClient) Accounts(ctx context.Context,
	options ...ClientOption,
) ([]*stripe.Account, bool, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, accountsEndpoint, nil)
	if err != nil {
		return nil, false, errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.apply(req)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.apiKey, "") // gfyrag: really weird authentication right?

	var httpResponse *http.Response

	httpResponse, err = d.httpClient.Do(req)
	if err != nil {
		return nil, false, errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, false, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type listResponse struct {
		AccountsListResponse
		Data []json.RawMessage `json:"data"`
	}

	rsp := &listResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, false, errors.Wrap(err, "decoding response")
	}

	accounts := make([]*stripe.Account, 0)

	if len(rsp.Data) > 0 {
		for _, data := range rsp.Data {
			account := &stripe.Account{}

			err = json.Unmarshal(data, &account)
			if err != nil {
				return nil, false, err
			}

			accounts = append(accounts, account)
		}
	}

	return accounts, rsp.HasMore, nil
}

func (d *DefaultClient) Balance(ctx context.Context, options ...ClientOption) (*stripe.Balance, error) {
	if d.stripeAccount == "" {
		return nil, errors.New("stripe account is required")
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, balanceEndpoint, nil)
	if err != nil {
		return nil, errors.Wrap(err, "creating http request")
	}

	for _, opt := range options {
		opt.apply(req)
	}

	req.Header.Set("Stripe-Account", d.stripeAccount)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(d.apiKey, "") // gfyrag: really weird authentication right?

	var httpResponse *http.Response

	httpResponse, err = d.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "doing request")
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", httpResponse.StatusCode)
	}

	type balanceResponse struct {
		BalanceResponse
	}

	rsp := &balanceResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(rsp)
	if err != nil {
		return nil, errors.Wrap(err, "decoding response")
	}

	return rsp.Balance, nil
}

func newHTTPClient() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

var _ Client = &DefaultClient{}
