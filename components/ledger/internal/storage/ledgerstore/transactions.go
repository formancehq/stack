package ledgerstore

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"strings"

	ledger "github.com/formancehq/ledger/internal"
	"github.com/formancehq/ledger/internal/storage/paginate"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/uptrace/bun"
)

const (
	MovesTableName = "moves"
)

var (
	metadataRegex = regexp.MustCompile("metadata\\[(.+)\\]")
)

type Transaction struct {
	bun.BaseModel `bun:"transactions,alias:transactions"`

	ID         *paginate.BigInt  `bun:"id,type:numeric"`
	Timestamp  ledger.Time       `bun:"timestamp,type:timestamp without time zone"`
	Reference  string            `bun:"reference,type:varchar,unique,nullzero"`
	Postings   []ledger.Posting  `bun:"postings,type:jsonb"`
	Metadata   metadata.Metadata `bun:"metadata,type:jsonb,default:'{}'"`
	RevertedAt *ledger.Time      `bun:"reverted_at"`
	LastUpdate *ledger.Time      `bun:"last_update"`
}

func (t *Transaction) toCore() *ledger.Transaction {
	return &ledger.Transaction{
		TransactionData: ledger.TransactionData{
			Reference: t.Reference,
			Metadata:  t.Metadata,
			Timestamp: t.Timestamp,
			Postings:  t.Postings,
		},
		ID:       (*big.Int)(t.ID),
		Reverted: t.RevertedAt != nil && !t.RevertedAt.IsZero(),
	}
}

type ExpandedTransaction struct {
	Transaction
	bun.BaseModel `bun:"transactions,alias:transactions"`

	ID                         *paginate.BigInt             `bun:"id,type:numeric"`
	Timestamp                  ledger.Time                  `bun:"timestamp,type:timestamp without time zone"`
	Reference                  string                       `bun:"reference,type:varchar,unique,nullzero"`
	Postings                   []ledger.Posting             `bun:"postings,type:jsonb"`
	Metadata                   metadata.Metadata            `bun:"metadata,type:jsonb,default:'{}'"`
	PostCommitEffectiveVolumes ledger.AccountsAssetsVolumes `bun:"post_commit_effective_volumes,type:jsonb"`
	PostCommitVolumes          ledger.AccountsAssetsVolumes `bun:"post_commit_volumes,type:jsonb"`
	RevertedAt                 *ledger.Time                 `bun:"reverted_at"`
	LastUpdate                 *ledger.Time                 `bun:"last_update"`
}

func (t *ExpandedTransaction) toCore() *ledger.ExpandedTransaction {
	var (
		preCommitEffectiveVolumes ledger.AccountsAssetsVolumes
		preCommitVolumes          ledger.AccountsAssetsVolumes
	)
	if t.PostCommitEffectiveVolumes != nil {
		preCommitEffectiveVolumes = t.PostCommitEffectiveVolumes.Copy()
		for _, posting := range t.Postings {
			preCommitEffectiveVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
			preCommitEffectiveVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		}
	}
	if t.PostCommitVolumes != nil {
		preCommitVolumes = t.PostCommitVolumes.Copy()
		for _, posting := range t.Postings {
			preCommitVolumes.AddOutput(posting.Source, posting.Asset, big.NewInt(0).Neg(posting.Amount))
			preCommitVolumes.AddInput(posting.Destination, posting.Asset, big.NewInt(0).Neg(posting.Amount))
		}
	}
	return &ledger.ExpandedTransaction{
		Transaction: ledger.Transaction{
			TransactionData: ledger.TransactionData{
				Reference: t.Reference,
				Metadata:  t.Metadata,
				Timestamp: t.Timestamp,
				Postings:  t.Postings,
			},
			ID:       (*big.Int)(t.ID),
			Reverted: t.RevertedAt != nil && !t.RevertedAt.IsZero(),
		},
		PreCommitEffectiveVolumes:  preCommitEffectiveVolumes,
		PostCommitEffectiveVolumes: t.PostCommitEffectiveVolumes,
		PreCommitVolumes:           preCommitVolumes,
		PostCommitVolumes:          t.PostCommitVolumes,
	}
}

type account string

var _ driver.Valuer = account("")

func (m1 account) Value() (driver.Value, error) {
	ret, err := json.Marshal(strings.Split(string(m1), ":"))
	if err != nil {
		return nil, err
	}
	return string(ret), nil
}

// Scan - Implement the database/sql scanner interface
func (m1 *account) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	v, err := driver.String.ConvertValue(value)
	if err != nil {
		return err
	}

	array := make([]string, 0)
	switch vv := v.(type) {
	case []uint8:
		err = json.Unmarshal(vv, &array)
	case string:
		err = json.Unmarshal([]byte(vv), &array)
	default:
		panic("not handled type")
	}
	if err != nil {
		return err
	}
	*m1 = account(strings.Join(array, ":"))
	return nil
}

func (store *Store) buildTransactionQuery(p PITFilterWithVolumes, query *bun.SelectQuery) *bun.SelectQuery {

	selectMetadata := query.NewSelect().
		Table("transactions_metadata").
		Where("transactions.id = transactions_metadata.transaction_id").
		Order("revision desc").
		Limit(1)

	if p.PIT != nil && !p.PIT.IsZero() {
		selectMetadata = selectMetadata.Where("date <= ?", p.PIT)
	}

	query = query.
		Table("transactions").
		ColumnExpr("distinct on(transactions.id) transactions.*, transactions_metadata.metadata").
		Join(fmt.Sprintf(`left join lateral (%s) as transactions_metadata on true`, selectMetadata.String()))

	if p.PIT != nil && !p.PIT.IsZero() {
		query = query.
			Where("timestamp <= ?", p.PIT).
			ColumnExpr(fmt.Sprintf("case when reverted_at is not null and reverted_at > '%s' then null else reverted_at end", p.PIT.Format(ledger.DateFormat)))
	}

	if p.ExpandEffectiveVolumes {
		query = query.ColumnExpr("get_aggregated_effective_volumes_for_transaction(transactions) as post_commit_effective_volumes")
	}
	if p.ExpandVolumes {
		query = query.ColumnExpr("get_aggregated_volumes_for_transaction(transactions) as post_commit_volumes")
	}
	return query
}

func (store *Store) transactionQueryContext(qb query.Builder) (string, []any, error) {

	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "reference" || key == "timestamp":
			return fmt.Sprintf("%s %s ?", key, query.DefaultComparisonOperatorsMapping[operator]), []any{value}, nil
		case key == "account":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, newErrInvalidQuery("'account' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddressOnTransactions(address, true, true), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'account'", address)
			}
		case key == "source":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, errors.New("'source' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddressOnTransactions(address, true, false), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'source'", address)
			}
		case key == "destination":
			// TODO: Should allow comparison operator only if segments not used
			if operator != "$match" {
				return "", nil, errors.New("'destination' column can only be used with $match")
			}
			switch address := value.(type) {
			case string:
				return filterAccountAddressOnTransactions(address, false, true), nil, nil
			default:
				return "", nil, newErrInvalidQuery("unexpected type %T for column 'destination'", address)
			}
		case metadataRegex.Match([]byte(key)):
			if operator != "$match" {
				return "", nil, newErrInvalidQuery("'account' column can only be used with $match")
			}
			match := metadataRegex.FindAllStringSubmatch(key, 3)

			return "metadata @> ?", []any{map[string]any{
				match[0][1]: value,
			}}, nil
		default:
			return "", nil, newErrInvalidQuery("unknown key '%s' when building query", key)
		}
	}))
}

func (store *Store) buildTransactionListQuery(selectQuery *bun.SelectQuery, q PaginatedQueryOptions[PITFilterWithVolumes], where string, args []any) *bun.SelectQuery {

	selectQuery = store.buildTransactionQuery(q.Options, selectQuery)
	if where != "" {
		return selectQuery.Where(where, args...)
	}

	return selectQuery
}

func (store *Store) GetTransactions(ctx context.Context, q GetTransactionsQuery) (*api.Cursor[ledger.ExpandedTransaction], error) {

	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = store.transactionQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	transactions, err := paginateWithColumn[PaginatedQueryOptions[PITFilterWithVolumes], ExpandedTransaction](store, ctx,
		(*paginate.ColumnPaginatedQuery[PaginatedQueryOptions[PITFilterWithVolumes]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return store.buildTransactionListQuery(query, q.Options, where, args)
		},
	)
	if err != nil {
		return nil, err
	}

	return api.MapCursor(transactions, func(from ExpandedTransaction) ledger.ExpandedTransaction {
		return *from.toCore()
	}), nil
}

func (store *Store) CountTransactions(ctx context.Context, q GetTransactionsQuery) (int, error) {

	var (
		where string
		args  []any
		err   error
	)

	if q.Options.QueryBuilder != nil {
		where, args, err = store.transactionQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return 0, err
		}
	}

	return count(store, ctx, func(query *bun.SelectQuery) *bun.SelectQuery {
		return store.buildTransactionListQuery(query, q.Options, where, args)
	})
}

func (store *Store) GetTransactionWithVolumes(ctx context.Context, filter GetTransactionQuery) (*ledger.ExpandedTransaction, error) {
	ret, err := fetch[*ExpandedTransaction](store, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return store.buildTransactionQuery(filter.PITFilterWithVolumes, query).
				Where("transactions.id = ?", filter.ID).
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return ret.toCore(), nil
}

func (store *Store) GetTransaction(ctx context.Context, txId *big.Int) (*ledger.Transaction, error) {
	tx, err := fetch[*Transaction](store, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Table("transactions").
				ColumnExpr(`transactions.id, transactions.reference, transactions.postings, transactions.timestamp, transactions.reverted_at, tm.metadata`).
				Join("left join transactions_metadata tm on tm.transaction_id = transactions.id").
				Where("transactions.id = ?", (*paginate.BigInt)(txId)).
				Order("tm.revision desc").
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return tx.toCore(), nil
}

func (store *Store) GetTransactionByReference(ctx context.Context, ref string) (*ledger.ExpandedTransaction, error) {
	ret, err := fetch[*ExpandedTransaction](store, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Table("transactions").
				ColumnExpr(`transactions.id, transactions.reference, transactions.postings, transactions.timestamp, transactions.reverted_at, tm.metadata`).
				Join("left join transactions_metadata tm on tm.transaction_id = transactions.id").
				Where("transactions.reference = ?", ref).
				Order("tm.revision desc").
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return ret.toCore(), nil
}

func (store *Store) GetLastTransaction(ctx context.Context) (*ledger.ExpandedTransaction, error) {
	ret, err := fetch[*ExpandedTransaction](store, ctx,
		func(query *bun.SelectQuery) *bun.SelectQuery {
			return query.
				Table("transactions").
				ColumnExpr(`transactions.id, transactions.reference, transactions.postings, transactions.timestamp, transactions.reverted_at, tm.metadata`).
				Join("left join transactions_metadata tm on tm.transaction_id = transactions.id").
				Order("transactions.id desc", "tm.revision desc").
				Limit(1)
		})
	if err != nil {
		return nil, err
	}

	return ret.toCore(), nil
}

type GetTransactionsQuery paginate.ColumnPaginatedQuery[PaginatedQueryOptions[PITFilterWithVolumes]]

func (q GetTransactionsQuery) WithExpandVolumes() GetTransactionsQuery {
	q.Options.Options.ExpandVolumes = true

	return q
}

func (q GetTransactionsQuery) WithExpandEffectiveVolumes() GetTransactionsQuery {
	q.Options.Options.ExpandEffectiveVolumes = true

	return q
}

func NewGetTransactionsQuery(options PaginatedQueryOptions[PITFilterWithVolumes]) GetTransactionsQuery {
	return GetTransactionsQuery{
		PageSize: options.PageSize,
		Column:   "id",
		Order:    paginate.OrderDesc,
		Options:  options,
	}
}

type GetTransactionQuery struct {
	PITFilterWithVolumes
	ID *big.Int
}

func (q GetTransactionQuery) WithExpandVolumes() GetTransactionQuery {
	q.ExpandVolumes = true

	return q
}

func (q GetTransactionQuery) WithExpandEffectiveVolumes() GetTransactionQuery {
	q.ExpandEffectiveVolumes = true

	return q
}

func NewGetTransactionQuery(id *big.Int) GetTransactionQuery {
	return GetTransactionQuery{
		PITFilterWithVolumes: PITFilterWithVolumes{},
		ID:                   id,
	}
}
