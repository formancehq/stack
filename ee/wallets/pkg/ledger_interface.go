package wallet

import (
	"context"
	"fmt"
	"math/big"
	"time"

	sdk "github.com/formancehq/formance-sdk-go"
	"github.com/formancehq/formance-sdk-go/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/pkg/models/shared"
	"github.com/formancehq/stack/libs/go-libs/collectionutils"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/pointer"
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

type PostTransaction struct {
	Metadata  map[string]string             `json:"metadata,omitempty"`
	Postings  []shared.Posting              `json:"postings,omitempty"`
	Reference *string                       `json:"reference,omitempty"`
	Script    *shared.PostTransactionScript `json:"script,omitempty"`
	Timestamp *time.Time                    `json:"timestamp,omitempty"`
}

type Account struct {
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata,omitempty"`
	Type     *string           `json:"type,omitempty"`
}

func (a Account) GetMetadata() map[string]string {
	return a.Metadata
}

func (a Account) GetAddress() string {
	return a.Address
}

type AccountWithVolumesAndBalances struct {
	Account
	Balances map[string]*big.Int      `json:"balances,omitempty"`
	Volumes  map[string]shared.Volume `json:"volumes,omitempty"`
}

func (a AccountWithVolumesAndBalances) GetBalances() map[string]*big.Int {
	return a.Balances
}

func (a AccountWithVolumesAndBalances) GetVolumes() map[string]shared.Volume {
	return a.Volumes
}

type AccountsCursorResponseCursor struct {
	Data     []Account `json:"data"`
	HasMore  bool      `json:"hasMore"`
	Next     *string   `json:"next,omitempty"`
	PageSize int64     `json:"pageSize"`
	Previous *string   `json:"previous,omitempty"`
}

func (c AccountsCursorResponseCursor) GetNext() *string {
	return c.Next
}

func (c AccountsCursorResponseCursor) GetPrevious() *string {
	return c.Previous
}

func (c AccountsCursorResponseCursor) GetData() []Account {
	return c.Data
}

func (c AccountsCursorResponseCursor) GetHasMore() bool {
	return c.HasMore
}

type Ledger interface {
	AddMetadataToAccount(ctx context.Context, ledger, account string, metadata map[string]string) error
	GetAccount(ctx context.Context, ledger, account string) (*AccountWithVolumesAndBalances, error)
	ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*AccountsCursorResponseCursor, error)
	ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*shared.TransactionsCursorResponseCursor, error)
	CreateTransaction(ctx context.Context, ledger string, postTransaction PostTransaction) (*shared.Transaction, error)
}

type DefaultLedger struct {
	client *sdk.Formance
}

func (d DefaultLedger) ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*shared.TransactionsCursorResponseCursor, error) {
	req := operations.ListTransactionsRequest{
		Ledger: ledger,
	}
	if query.Cursor == "" {
		req.PageSize = pointer.For(int64(query.Limit))
		req.Destination = pointer.For(query.Destination)
		req.Source = pointer.For(query.Source)
		req.Account = pointer.For(query.Account)
		req.Metadata = make(map[string]any)

		for key, value := range query.Metadata {
			req.Metadata[fmt.Sprintf("metadata[%s]", key)] = value
		}
	} else {
		req.Cursor = pointer.For(query.Cursor)
	}

	rsp, err := d.client.Ledger.ListTransactions(ctx, req)
	if err != nil {
		return nil, err
	}

	return &rsp.TransactionsCursorResponse.Cursor, nil
}

func (d DefaultLedger) CreateTransaction(ctx context.Context, ledger string, transaction PostTransaction) (*shared.Transaction, error) {
	txMetadata := make(map[string]any, 0)
	for k, v := range transaction.Metadata {
		txMetadata[k] = v
	}
	ret, err := d.client.Ledger.CreateTransaction(ctx, operations.CreateTransactionRequest{
		PostTransaction: shared.PostTransaction{
			Metadata:  txMetadata,
			Postings:  transaction.Postings,
			Reference: transaction.Reference,
			Script:    transaction.Script,
			Timestamp: transaction.Timestamp,
		},
		Ledger: ledger,
	})
	if err != nil {
		return nil, err
	}

	return &ret.TransactionsResponse.Data[0], nil
}

func (d DefaultLedger) AddMetadataToAccount(ctx context.Context, ledger, account string, metadata map[string]string) error {

	m := make(map[string]any)
	for k, v := range metadata {
		m[k] = v
	}

	_, err := d.client.Ledger.AddMetadataToAccount(ctx, operations.AddMetadataToAccountRequest{
		RequestBody: m,
		Address:     account,
		Ledger:      ledger,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultLedger) GetAccount(ctx context.Context, ledger, account string) (*AccountWithVolumesAndBalances, error) {
	ret, err := d.client.Ledger.GetAccount(ctx, operations.GetAccountRequest{
		Address: account,
		Ledger:  ledger,
	})
	if err != nil {
		return nil, err
	}

	return &AccountWithVolumesAndBalances{
		Account: Account{
			Address:  ret.AccountResponse.Data.Address,
			Metadata: convertAccountMetadata(ret.AccountResponse.Data.Metadata),
			Type:     ret.AccountResponse.Data.Type,
		},
		Balances: ret.AccountResponse.Data.Balances,
		Volumes:  ret.AccountResponse.Data.Volumes,
	}, nil
}

func (d DefaultLedger) ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*AccountsCursorResponseCursor, error) {
	req := operations.ListAccountsRequest{
		Ledger: ledger,
	}
	if query.Cursor == "" {
		req.PageSize = pointer.For(int64(query.Limit))
		req.Metadata = make(map[string]any)
		for key, value := range query.Metadata {
			req.Metadata[key] = value
		}
	} else {
		req.Cursor = pointer.For(query.Cursor)
	}

	ret, err := d.client.Ledger.ListAccounts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AccountsCursorResponseCursor{
		Data: collectionutils.Map(ret.AccountsCursorResponse.Cursor.Data, func(from shared.Account) Account {
			return Account{
				Address:  from.Address,
				Metadata: convertAccountMetadata(from.Metadata),
				Type:     from.Type,
			}
		}),
		HasMore:  ret.AccountsCursorResponse.Cursor.HasMore,
		Next:     ret.AccountsCursorResponse.Cursor.Next,
		PageSize: ret.AccountsCursorResponse.Cursor.PageSize,
		Previous: ret.AccountsCursorResponse.Cursor.Previous,
	}, nil
}

var _ Ledger = &DefaultLedger{}

func NewDefaultLedger(client *sdk.Formance) *DefaultLedger {
	return &DefaultLedger{
		client: client,
	}
}

func convertAccountMetadata(m map[string]any) map[string]string {
	ret := make(map[string]string)
	for k, v := range m {
		switch v := v.(type) {
		case string:
			ret[k] = v
		case map[string]any:
			ret[k] = metadata.MarshalValue(v)
		default:
			ret[k] = fmt.Sprint(v)
		}
	}
	return ret
}
