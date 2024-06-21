package wallet

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	stdtime "time"

	"github.com/formancehq/formance-sdk-go/v2/pkg/models/sdkerrors"

	"github.com/formancehq/stack/libs/go-libs/query"

	"github.com/formancehq/stack/libs/go-libs/time"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/operations"
	"github.com/formancehq/formance-sdk-go/v2/pkg/models/shared"
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
	Metadata  map[string]string               `json:"metadata,omitempty"`
	Postings  []shared.V2Posting              `json:"postings,omitempty"`
	Reference *string                         `json:"reference,omitempty"`
	Script    *shared.V2PostTransactionScript `json:"script,omitempty"`
	Timestamp *time.Time                      `json:"timestamp,omitempty"`
}

type Account struct {
	Address  string            `json:"address"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

func (a Account) GetMetadata() map[string]string {
	return a.Metadata
}

func (a Account) GetAddress() string {
	return a.Address
}

type AccountWithVolumesAndBalances struct {
	Account
	Balances map[string]*big.Int        `json:"balances,omitempty"`
	Volumes  map[string]shared.V2Volume `json:"volumes,omitempty"`
}

func (a AccountWithVolumesAndBalances) GetBalances() map[string]*big.Int {
	return a.Balances
}

func (a AccountWithVolumesAndBalances) GetVolumes() map[string]shared.V2Volume {
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
	EnsureLedgerExists(ctx context.Context, name string) error
	AddMetadataToAccount(ctx context.Context, ledger, account string, metadata map[string]string) error
	GetAccount(ctx context.Context, ledger, account string) (*AccountWithVolumesAndBalances, error)
	ListAccounts(ctx context.Context, ledger string, query ListAccountsQuery) (*AccountsCursorResponseCursor, error)
	ListTransactions(ctx context.Context, ledger string, query ListTransactionsQuery) (*shared.V2TransactionsCursorResponseCursor, error)
	CreateTransaction(ctx context.Context, ledger string, postTransaction PostTransaction) (*shared.V2Transaction, error)
}

type DefaultLedger struct {
	client *sdk.Formance
}

func (d DefaultLedger) EnsureLedgerExists(ctx context.Context, name string) error {
	_, err := d.client.Ledger.V2GetLedger(ctx, operations.V2GetLedgerRequest{
		Ledger: name,
	})
	if err == nil {
		return nil
	}

	switch err := err.(type) {
	case *sdkerrors.V2ErrorResponse:
		if err.ErrorCode != shared.V2ErrorsEnumNotFound {
			return err
		}
	default:
		return err
	}

	_, err = d.client.Ledger.V2CreateLedger(ctx, operations.V2CreateLedgerRequest{
		V2CreateLedgerRequest: &shared.V2CreateLedgerRequest{
			Bucket: pointer.For(name),
		},
		Ledger: name,
	})
	return err
}

func (d DefaultLedger) ListTransactions(ctx context.Context, ledger string, q ListTransactionsQuery) (*shared.V2TransactionsCursorResponseCursor, error) {
	req := operations.V2ListTransactionsRequest{
		Ledger: ledger,
	}
	if q.Cursor == "" {
		req.PageSize = pointer.For(int64(q.Limit))
		conditions := make([]query.Builder, 0)
		if q.Destination != "" {
			conditions = append(conditions, query.Match("destination", q.Destination))
		}
		if q.Source != "" {
			conditions = append(conditions, query.Match("source", q.Source))
		}
		if q.Account != "" {
			conditions = append(conditions, query.Match("account", q.Source))
		}
		if q.Metadata != nil {
			for k, v := range q.Metadata {
				conditions = append(conditions, query.Match(fmt.Sprintf("metadata[%s]", k), v))
			}
		}
		if len(conditions) > 0 {
			data, err := json.Marshal(query.And(conditions...))
			if err != nil {
				panic(err)
			}
			body := make(map[string]any)
			if err := json.Unmarshal(data, &body); err != nil {
				panic(err)
			}
			req.RequestBody = body
		}
	} else {
		req.Cursor = pointer.For(q.Cursor)
	}

	rsp, err := d.client.Ledger.V2ListTransactions(ctx, req)
	if err != nil {
		return nil, err
	}

	return &rsp.V2TransactionsCursorResponse.Cursor, nil
}

func (d DefaultLedger) CreateTransaction(ctx context.Context, ledger string, transaction PostTransaction) (*shared.V2Transaction, error) {
	ret, err := d.client.Ledger.V2CreateTransaction(ctx, operations.V2CreateTransactionRequest{
		V2PostTransaction: shared.V2PostTransaction{
			Metadata:  transaction.Metadata,
			Postings:  transaction.Postings,
			Reference: transaction.Reference,
			Script:    transaction.Script,
			Timestamp: func() *stdtime.Time {
				if transaction.Timestamp == nil {
					return nil
				}
				return &transaction.Timestamp.Time
			}(),
		},
		Ledger: ledger,
	})
	if err != nil {
		return nil, err
	}

	return &ret.V2CreateTransactionResponse.Data, nil
}

func (d DefaultLedger) AddMetadataToAccount(ctx context.Context, ledger, account string, metadata map[string]string) error {

	_, err := d.client.Ledger.V2AddMetadataToAccount(ctx, operations.V2AddMetadataToAccountRequest{
		RequestBody: metadata,
		Address:     account,
		Ledger:      ledger,
	})
	if err != nil {
		return err
	}
	return nil
}

func (d DefaultLedger) GetAccount(ctx context.Context, ledger, account string) (*AccountWithVolumesAndBalances, error) {
	ret, err := d.client.Ledger.V2GetAccount(ctx, operations.V2GetAccountRequest{
		Address: account,
		Ledger:  ledger,
		Expand:  pointer.For("volumes"),
	})
	if err != nil {
		return nil, err
	}

	balances := make(map[string]*big.Int)
	for asset, volumes := range ret.V2AccountResponse.Data.Volumes {
		balances[asset] = big.NewInt(0).Sub(volumes.Input, volumes.Output)
	}

	return &AccountWithVolumesAndBalances{
		Account: Account{
			Address:  ret.V2AccountResponse.Data.Address,
			Metadata: ret.V2AccountResponse.Data.Metadata,
		},
		Balances: balances,
		Volumes:  ret.V2AccountResponse.Data.Volumes,
	}, nil
}

func (d DefaultLedger) ListAccounts(ctx context.Context, ledger string, q ListAccountsQuery) (*AccountsCursorResponseCursor, error) {
	req := operations.V2ListAccountsRequest{
		Ledger: ledger,
	}
	if q.Cursor == "" {
		req.PageSize = pointer.For(int64(q.Limit))

		conditions := make([]query.Builder, 0)
		if q.Metadata != nil {
			for k, v := range q.Metadata {
				conditions = append(conditions, query.Match(fmt.Sprintf("metadata[%s]", k), v))
			}
		}
		if len(conditions) > 0 {
			data, err := json.Marshal(query.And(conditions...))
			if err != nil {
				panic(err)
			}
			body := make(map[string]any)
			if err := json.Unmarshal(data, &body); err != nil {
				panic(err)
			}
			req.RequestBody = body
		}
	} else {
		req.Cursor = pointer.For(q.Cursor)
	}

	ret, err := d.client.Ledger.V2ListAccounts(ctx, req)
	if err != nil {
		return nil, err
	}

	return &AccountsCursorResponseCursor{
		Data: collectionutils.Map(ret.V2AccountsCursorResponse.Cursor.Data, func(from shared.V2Account) Account {
			return Account{
				Address:  from.Address,
				Metadata: from.Metadata,
			}
		}),
		HasMore:  ret.V2AccountsCursorResponse.Cursor.HasMore,
		Next:     ret.V2AccountsCursorResponse.Cursor.Next,
		PageSize: ret.V2AccountsCursorResponse.Cursor.PageSize,
		Previous: ret.V2AccountsCursorResponse.Cursor.Previous,
	}, nil
}

var _ Ledger = &DefaultLedger{}

func NewDefaultLedger(client *sdk.Formance) *DefaultLedger {
	return &DefaultLedger{
		client: client,
	}
}
