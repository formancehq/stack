package storage

import (
	"context"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

func (s *Storage) CreatePool(ctx context.Context, pool *models.Pool) error {
	var id uuid.UUID
	err := s.db.NewInsert().
		Model(pool).
		Returning("id").
		Scan(ctx, &id)
	if err != nil {
		return e("failed to create pool", err)
	}
	pool.ID = id

	return nil
}

func (s *Storage) AddAccountsToPool(ctx context.Context, poolAccounts []*models.PoolAccounts) error {
	_, err := s.db.NewInsert().
		Model(&poolAccounts).
		Exec(ctx)
	if err != nil {
		return e("failed to add accounts to pool", err)
	}

	return nil
}

func (s *Storage) AddAccountToPool(ctx context.Context, poolAccount *models.PoolAccounts) error {
	_, err := s.db.NewInsert().
		Model(poolAccount).
		Exec(ctx)
	if err != nil {
		return e("failed to add account to pool", err)
	}

	return nil
}

func (s *Storage) RemoveAccountFromPool(ctx context.Context, poolAccount *models.PoolAccounts) error {
	_, err := s.db.NewDelete().
		Model(poolAccount).
		Where("pool_id = ?", poolAccount.PoolID).
		Where("account_id = ?", poolAccount.AccountID).
		Exec(ctx)
	if err != nil {
		return e("failed to remove account from pool", err)
	}

	return nil
}

type PoolQuery struct{}

type ListPoolsQuery bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PoolQuery]]

func NewListPoolsQuery(opts PaginatedQueryOptions[PoolQuery]) ListPoolsQuery {
	return ListPoolsQuery{
		PageSize: opts.PageSize,
		Order:    bunpaginate.OrderAsc,
		Options:  opts,
	}
}

func (s *Storage) ListPools(ctx context.Context, q ListPoolsQuery) (*api.Cursor[models.Pool], error) {
	cursor, err := PaginateWithOffset[PaginatedQueryOptions[PoolQuery], models.Pool](s, ctx,
		(*bunpaginate.OffsetPaginatedQuery[PaginatedQueryOptions[PoolQuery]])(&q),
		func(query *bun.SelectQuery) *bun.SelectQuery {
			query = query.
				Relation("PoolAccounts").
				Order("created_at DESC")

			if q.Options.Sorter != nil {
				query = q.Options.Sorter.Apply(query)
			}

			return query
		},
	)
	return cursor, err
}

func (s *Storage) GetPool(ctx context.Context, poolID uuid.UUID) (*models.Pool, error) {
	var pool models.Pool

	err := s.db.NewSelect().
		Model(&pool).
		Where("id = ?", poolID).
		Relation("PoolAccounts").
		Scan(ctx)
	if err != nil {
		return nil, e("failed to get pool", err)
	}

	return &pool, nil
}

func (s *Storage) DeletePool(ctx context.Context, poolID uuid.UUID) error {
	_, err := s.db.NewDelete().
		Model(&models.Pool{}).
		Where("id = ?", poolID).
		Exec(ctx)
	if err != nil {
		return e("failed to delete pool", err)
	}

	return nil
}
