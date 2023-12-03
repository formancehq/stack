package storage

import (
	"context"
	"sort"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
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
		Exec(ctx)
	if err != nil {
		return e("failed to remove account from pool", err)
	}

	return nil
}

func (s *Storage) ListPools(ctx context.Context, pagination PaginatorQuery) ([]*models.Pool, PaginationDetails, error) {
	var pools []*models.Pool

	query := s.db.NewSelect().
		Model(&pools).
		Relation("PoolAccounts")

	query = pagination.apply(query, "pools.created_at")

	err := query.Scan(ctx)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to list pools", err)
	}

	var (
		hasMore                       = len(pools) > pagination.pageSize
		hasPrevious                   bool
		firstReference, lastReference string
	)

	if hasMore {
		if pagination.cursor.Next || pagination.cursor.Reference == "" {
			pools = pools[:pagination.pageSize]
		} else {
			pools = pools[1:]
		}
	}

	sort.Slice(pools, func(i, j int) bool {
		return pools[i].CreatedAt.After(pools[j].CreatedAt)
	})

	if len(pools) > 0 {
		firstReference = pools[0].CreatedAt.Format(time.RFC3339Nano)
		lastReference = pools[len(pools)-1].CreatedAt.Format(time.RFC3339Nano)

		query = s.db.NewSelect().Model(&pools)

		hasPrevious, err = pagination.hasPrevious(ctx, query, "pools.created_at", firstReference)
		if err != nil {
			return nil, PaginationDetails{}, e("failed to check if there is a previous page", err)
		}
	}

	paginationDetails, err := pagination.paginationDetails(hasMore, hasPrevious, firstReference, lastReference)
	if err != nil {
		return nil, PaginationDetails{}, e("failed to get pagination details", err)
	}

	return pools, paginationDetails, nil
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
