package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/metadata"
)

type ListAccountsQuery struct {
	Cursor   string
	Limit    int
	Metadata metadata.Metadata
}

type ListTransactionsQuery struct {
	Cursor      string
	Limit       int
	Metadata    metadata.Metadata
	Destination string
	Source      string
	Account     string
}

type TransactionsCursorResponseCursor struct {
	PageSize int64                 `json:"pageSize"`
	HasMore  bool                  `json:"hasMore"`
	Previous string                `json:"previous,omitempty"`
	Next     string                `json:"next,omitempty"`
	Data     []ExpandedTransaction `json:"data"`
}

func (c TransactionsCursorResponseCursor) GetData() []ExpandedTransaction {
	return c.Data
}
func (c TransactionsCursorResponseCursor) GetNext() string {
	return c.Next
}
func (c TransactionsCursorResponseCursor) GetPrevious() string {
	return c.Previous
}
func (c TransactionsCursorResponseCursor) GetHasMore() bool {
	return c.HasMore
}

type TransactionsCursorResponse struct {
	Cursor TransactionsCursorResponseCursor `json:"cursor"`
}

type CreateTransactionResponse struct {
	Data Transaction `json:"data"`
}

type PostTransactionScript struct {
	Plain string         `json:"plain"`
	Vars  map[string]any `json:"vars,omitempty"`
}

type PostTransaction struct {
	Timestamp *time.Time             `json:"timestamp,omitempty"`
	Postings  []Posting              `json:"postings,omitempty"`
	Script    *PostTransactionScript `json:"script,omitempty"`
	Reference *string                `json:"reference,omitempty"`
	Metadata  map[string]string      `json:"metadata"`
}

type AccountsCursorResponseCursor struct {
	PageSize int64     `json:"pageSize"`
	HasMore  bool      `json:"hasMore"`
	Previous string    `json:"previous,omitempty"`
	Next     string    `json:"next,omitempty"`
	Data     []Account `json:"data"`
}

func (c AccountsCursorResponseCursor) GetData() []Account {
	return c.Data
}

func (c AccountsCursorResponseCursor) GetNext() string {
	return c.Next
}

func (c AccountsCursorResponseCursor) GetPrevious() string {
	return c.Previous
}

func (c AccountsCursorResponseCursor) GetHasMore() bool {
	return c.HasMore
}

type AccountsCursorResponse struct {
	Cursor AccountsCursorResponseCursor `json:"cursor"`
}

type Ledger interface {
	AddMetadataToAccount(ctx context.Context, ledger, account string, metadata metadata.Metadata) error
	GetAccount(ctx context.Context, ledger, account string) (*AccountWithVolumesAndBalances, error)
	ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*AccountsCursorResponseCursor, error)
	ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*TransactionsCursorResponseCursor, error)
	CreateTransaction(ctx context.Context, ledger string, postTransaction PostTransaction) (*CreateTransactionResponse, error)
}

type DefaultLedger struct {
	client  *http.Client
	baseUrl string
}

func (d DefaultLedger) ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*TransactionsCursorResponseCursor, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/transactions", d.baseUrl, ledger), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)
	urlValues := req.URL.Query()

	if query.Cursor == "" {
		urlValues.Set("pageSize", fmt.Sprint(query.Limit))
		urlValues.Set("destination", query.Destination)
		urlValues.Set("account", query.Account)
		urlValues.Set("source", query.Source)
		for key, value := range query.Metadata {
			urlValues.Set(fmt.Sprintf("metadata[%s]", key), value)
		}
	} else {
		urlValues.Set("cursor", query.Cursor)
	}
	req.URL.RawQuery = urlValues.Encode()
	if err != nil {
		return nil, err
	}

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		errorResponse := &shared.V2ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	ret := &TransactionsCursorResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(ret); err != nil {
		return nil, err
	}

	return &ret.Cursor, nil
}

func (d DefaultLedger) CreateTransaction(ctx context.Context, ledger string, transaction PostTransaction) (*CreateTransactionResponse, error) {

	data, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/transactions", d.baseUrl, ledger), bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		errorResponse := &shared.V2ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	// TODO(gfyrag): Remove when ledger v2 will be released
	type v1Response struct {
		Data []Transaction `json:"data"`
	}

	data, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	v1 := &v1Response{}
	if err := json.Unmarshal(data, v1); err != nil {
		v2 := &CreateTransactionResponse{}
		if err := json.Unmarshal(data, v2); err != nil {
			return nil, err
		}
		return v2, nil
	}

	return &CreateTransactionResponse{
		Data: v1.Data[0],
	}, nil
}

func (d DefaultLedger) AddMetadataToAccount(ctx context.Context, ledger, account string, metadata metadata.Metadata) error {

	data, err := json.Marshal(metadata)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/accounts/%s/metadata", d.baseUrl, ledger, account), bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
		errorResponse := &shared.V2ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	return err
}

func (d DefaultLedger) GetAccount(ctx context.Context, ledger, account string) (*AccountWithVolumesAndBalances, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/accounts/%s", d.baseUrl, ledger, account), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		errorResponse := &shared.V2ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	type accountResponse struct {
		Data AccountWithVolumesAndBalances `json:"data"`
	}

	ret := &accountResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(ret); err != nil {
		return nil, err
	}

	return &ret.Data, nil
}

func (d DefaultLedger) ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*AccountsCursorResponseCursor, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/accounts", d.baseUrl, ledger), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(ctx)
	urlValues := req.URL.Query()

	if query.Cursor == "" {
		urlValues.Set("pageSize", fmt.Sprint(query.Limit))
		for key, value := range query.Metadata {
			urlValues.Set(fmt.Sprintf("metadata[%s]", key), value)
		}
	} else {
		urlValues.Set("cursor", query.Cursor)
	}
	req.URL.RawQuery = urlValues.Encode()
	if err != nil {
		return nil, err
	}

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		errorResponse := &shared.V2ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	ret := &AccountsCursorResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(ret); err != nil {
		return nil, err
	}

	return &ret.Cursor, nil
}

var _ Ledger = &DefaultLedger{}

func NewDefaultLedger(client *http.Client, baseURL string) *DefaultLedger {
	return &DefaultLedger{
		client:  client,
		baseUrl: baseURL,
	}
}
