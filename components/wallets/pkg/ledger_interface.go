package wallet

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go"
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

type Ledger interface {
	AddMetadataToAccount(ctx context.Context, ledger, account string, metadata metadata.Metadata) error
	GetAccount(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error)
	ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*sdk.AccountsCursorResponseCursor, error)
	ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*sdk.TransactionsCursorResponseCursor, error)
	CreateTransaction(ctx context.Context, ledger string, postTransaction sdk.PostTransaction) (*sdk.CreateTransactionResponse, error)
}

type DefaultLedger struct {
	client  *http.Client
	baseUrl string
}

func (d DefaultLedger) ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*sdk.TransactionsCursorResponseCursor, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/transactions", d.baseUrl, ledger), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(req.Context())
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
		errorResponse := &sdk.ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	ret := &sdk.TransactionsCursorResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(ret); err != nil {
		return nil, err
	}

	return &ret.Cursor, nil
}

func (d DefaultLedger) CreateTransaction(ctx context.Context, ledger string, transaction sdk.PostTransaction) (*sdk.CreateTransactionResponse, error) {

	data, err := json.Marshal(transaction)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/%s/transactions", d.baseUrl, ledger), bytes.NewBuffer(data))
	if err != nil {
		panic(err)
	}
	req = req.WithContext(req.Context())

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusCreated {
		errorResponse := &sdk.ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	ret := &sdk.CreateTransactionResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(ret); err != nil {
		return nil, err
	}

	return ret, err
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
	req = req.WithContext(req.Context())

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK && httpResponse.StatusCode != http.StatusNoContent {
		errorResponse := &sdk.ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	return err
}

func (d DefaultLedger) GetAccount(ctx context.Context, ledger, account string) (*sdk.AccountWithVolumesAndBalances, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/accounts/%s", d.baseUrl, ledger, account), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(req.Context())

	httpResponse, err := d.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	if httpResponse.StatusCode != http.StatusOK {
		errorResponse := &sdk.ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	ret := &sdk.AccountResponse{}
	if err := json.NewDecoder(httpResponse.Body).Decode(ret); err != nil {
		return nil, err
	}

	return &ret.Data, nil
}

func (d DefaultLedger) ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*sdk.AccountsCursorResponseCursor, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/%s/accounts", d.baseUrl, ledger), nil)
	if err != nil {
		panic(err)
	}
	req = req.WithContext(req.Context())
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
		errorResponse := &sdk.ErrorResponse{}
		if err := json.NewDecoder(httpResponse.Body).Decode(errorResponse); err != nil {
			panic(err)
		}
		return nil, fmt.Errorf("%s", errorResponse.ErrorMessage)
	}

	ret := &sdk.AccountsCursorResponse{}
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
