package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/formancehq/stack/libs/go-libs/pointer"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

type pool struct {
	bun.BaseModel `bun:"table:pools"`

	// Mandatory fields
	ID        uuid.UUID `bun:"id,pk,type:uuid,notnull"`
	Name      string    `bun:"name,type:text,notnull"`
	CreatedAt time.Time `bun:"created_at,type:timestamp without time zone,notnull"`

	PoolAccounts []*poolAccounts `bun:"rel:has-many,join:id=pool_id"`
}

type poolAccounts struct {
	bun.BaseModel `bun:"table:pool_accounts"`

	PoolID    uuid.UUID        `bun:"pool_id,pk,type:uuid,notnull"`
	AccountID models.AccountID `bun:"account_id,pk,type:character varying,notnull"`
}

func (s *store) PoolsUpsert(ctx context.Context, pool models.Pool) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("begin transaction: %w", err)
	}
	defer tx.Rollback()

	poolToInsert, accountsToInsert := fromPoolModel(pool)

	_, err = tx.NewInsert().
		Model(&poolToInsert).
		On("CONFLICT (id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("insert pool: %w", err)
	}

	_, err = tx.NewInsert().
		Model(&accountsToInsert).
		On("CONFLICT (pool_id, account_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("insert pool accounts: %w", err)
	}

	return e("commit transaction: %w", tx.Commit())
}

func (s *store) PoolsGet(ctx context.Context, id uuid.UUID) (*models.Pool, error) {
	var pool pool
	err := s.db.NewSelect().
		Model(&pool).
		Relation("PoolAccounts").
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return nil, e("get pool: %w", err)
	}

	return pointer.For(toPoolModel(pool)), nil
}

func (s *store) PoolsDelete(ctx context.Context, id uuid.UUID) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return e("begin transaction: %w", err)
	}
	defer tx.Rollback()

	_, err = tx.NewDelete().
		Model((*pool)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("delete pool: %w", err)
	}

	_, err = tx.NewDelete().
		Model((*poolAccounts)(nil)).
		Where("pool_id = ?", id).
		Exec(ctx)
	if err != nil {
		return e("delete pool accounts: %w", err)
	}

	return e("commit transaction: %w", tx.Commit())
}

func (s *store) PoolsAddAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	_, err := s.db.NewInsert().
		Model(&poolAccounts{
			PoolID:    id,
			AccountID: accountID,
		}).
		On("CONFLICT (pool_id, account_id) DO NOTHING").
		Exec(ctx)
	if err != nil {
		return e("insert pool account: %w", err)
	}
	return nil
}

func (s *store) PoolsRemoveAccount(ctx context.Context, id uuid.UUID, accountID models.AccountID) error {
	_, err := s.db.NewDelete().
		Model((*poolAccounts)(nil)).
		Where("pool_id = ? AND account_id = ?", id, accountID).
		Exec(ctx)
	if err != nil {
		return e("delete pool account: %w", err)
	}
	return nil
}

type PoolQuery struct{}

type ListPoolsQuery bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[PoolQuery]]

func NewListPoolsQuery(opts bunpaginate.PaginatedQueryOptions[PoolQuery]) ListPoolsQuery {
	return ListPoolsQuery{
		Order:    bunpaginate.OrderAsc,
		PageSize: opts.PageSize,
		Options:  opts,
	}
}

func (s *store) poolsQueryContext(qb query.Builder) (string, []any, error) {
	return qb.Build(query.ContextFn(func(key, operator string, value any) (string, []any, error) {
		switch {
		case key == "name":
			if operator != "$match" {
				return "", nil, errors.Wrap(ErrValidation, "'metadata' column can only be used with $match")
			}

			return fmt.Sprintf("%s = ?", key), []any{value}, nil
			// TODO(polo): add filters for accounts ID
		default:
			return "", nil, errors.Wrap(ErrValidation, fmt.Sprintf("unknown key '%s' when building query", key))
		}
	}))
}

func (s *store) PoolsList(ctx context.Context, q ListPoolsQuery) (*bunpaginate.Cursor[models.Pool], error) {
	var (
		where string
		args  []any
		err   error
	)
	if q.Options.QueryBuilder != nil {
		where, args, err = s.poolsQueryContext(q.Options.QueryBuilder)
		if err != nil {
			return nil, err
		}
	}

	cursor, err := paginateWithOffset[bunpaginate.PaginatedQueryOptions[PoolQuery], pool](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[bunpaginate.PaginatedQueryOptions[PoolQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Relation("PoolAccounts")

			if where != "" {
				query = query.Where(where, args...)
			}

			query = query.Order("created_at DESC")

			return query
		},
	)
	if err != nil {
		return nil, e("failed to fetch pools", err)
	}

	pools := make([]models.Pool, 0, len(cursor.Data))
	for _, p := range cursor.Data {
		pools = append(pools, toPoolModel(p))
	}

	return &bunpaginate.Cursor[models.Pool]{
		PageSize: cursor.PageSize,
		HasMore:  cursor.HasMore,
		Previous: cursor.Previous,
		Next:     cursor.Next,
		Data:     pools,
	}, nil
}

func fromPoolModel(from models.Pool) (pool, []poolAccounts) {
	p := pool{
		ID:        from.ID,
		Name:      from.Name,
		CreatedAt: from.CreatedAt,
	}

	var accounts []poolAccounts
	for _, pa := range from.PoolAccounts {
		accounts = append(accounts, poolAccounts{
			PoolID:    pa.PoolID,
			AccountID: pa.AccountID,
		})
	}

	return p, accounts
}

func toPoolModel(from pool) models.Pool {
	var accounts []models.PoolAccounts
	for _, pa := range from.PoolAccounts {
		accounts = append(accounts, models.PoolAccounts{
			PoolID:    pa.PoolID,
			AccountID: pa.AccountID,
		})
	}

	return models.Pool{
		ID:           from.ID,
		Name:         from.Name,
		CreatedAt:    from.CreatedAt,
		PoolAccounts: accounts,
	}
}
