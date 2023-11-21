package client

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// HTTP client howto: https://www.youtube.com/watch?v=evorkFq3Y5k

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
	Accounts(ctx context.Context, options ...ClientOption) ([]*Account, string, error)
	Transactions(ctx context.Context, options ...ClientOption) ([]*Transaction, string, error)
	ForAccount(account string) Client
}

type DefaultClient struct {
	httpClient   *http.Client
	baseUrl      string
	accessKey    string
	secret       string
	atlarAccount string
}

func NewDefaultClient(baseUrl, accessKey, secret string) *DefaultClient {
	return &DefaultClient{
		httpClient: newHTTPClient(),
		baseUrl:    baseUrl,
		accessKey:  accessKey,
		secret:     secret,
	}
}

func (d *DefaultClient) ForAccount(account string) Client {
	cp := *d
	cp.atlarAccount = account

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

type ErrorResponse struct {
	Status  uint   `json:"status"`
	Message string `json:"message"`
}

type ListResponse struct {
	NextToken string        `json:"nextToken"`
	Token     string        `json:"token"`
	Limit     int           `json:"limit"`
	Items     []interface{} `json:"items"`
}
