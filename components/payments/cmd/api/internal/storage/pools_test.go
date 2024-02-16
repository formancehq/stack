package storage

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func insertPools(t *testing.T, store *Storage, accountIDs []models.AccountID) []uuid.UUID {
	pool1 := models.Pool{
		Name:      "test",
		CreatedAt: time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
	}
	var uuid1 uuid.UUID
	err := store.DB().NewInsert().
		Model(&pool1).
		Returning("id").
		Scan(context.Background(), &uuid1)
	require.NoError(t, err)

	poolAccounts1 := models.PoolAccounts{
		PoolID:    uuid1,
		AccountID: accountIDs[0],
	}
	_, err = store.DB().NewInsert().
		Model(&poolAccounts1).
		Exec(context.Background())
	require.NoError(t, err)

	var uuid2 uuid.UUID
	pool2 := models.Pool{
		Name:      "test2",
		CreatedAt: time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
	}
	err = store.DB().NewInsert().
		Model(&pool2).
		Returning("id").
		Scan(context.Background(), &uuid2)
	require.NoError(t, err)

	poolAccounts2 := []*models.PoolAccounts{
		{
			PoolID:    uuid2,
			AccountID: accountIDs[0],
		},
		{
			PoolID:    uuid2,
			AccountID: accountIDs[1],
		},
	}
	_, err = store.DB().NewInsert().
		Model(&poolAccounts2).
		Exec(context.Background())
	require.NoError(t, err)

	return []uuid.UUID{uuid1, uuid2}
}

func TestCreatePools(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	accounts := insertAccounts(t, store, connectorID)

	pool := &models.Pool{
		Name:         "test",
		CreatedAt:    time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		PoolAccounts: []*models.PoolAccounts{},
	}
	for _, account := range accounts {
		pool.PoolAccounts = append(pool.PoolAccounts, &models.PoolAccounts{
			AccountID: account,
		})
	}

	err := store.CreatePool(context.Background(), pool)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, pool.ID)
}

func TestAddAccountsToPool(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	accounts := insertAccounts(t, store, connectorID)
	poolIDs := insertPools(t, store, accounts)

	poolAccounts := []*models.PoolAccounts{
		{
			PoolID:    poolIDs[0],
			AccountID: accounts[1],
		},
	}

	err := store.AddAccountsToPool(context.Background(), poolAccounts)
	require.NoError(t, err)

	pool, err := store.GetPool(context.Background(), poolIDs[0])
	require.NoError(t, err)
	require.Equal(t, 2, len(pool.PoolAccounts))
	require.Equal(t, accounts[0], pool.PoolAccounts[0].AccountID)
	require.Equal(t, accounts[1], pool.PoolAccounts[1].AccountID)
}

func TestRemoveAccoutsToPool(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	accounts := insertAccounts(t, store, connectorID)
	poolIDs := insertPools(t, store, accounts)

	poolAccounts := []*models.PoolAccounts{
		{
			PoolID:    poolIDs[0],
			AccountID: accounts[0],
		},
	}

	err := store.RemoveAccountFromPool(context.Background(), poolAccounts[0])
	require.NoError(t, err)

	pool, err := store.GetPool(context.Background(), poolIDs[0])
	require.NoError(t, err)
	require.Equal(t, 0, len(pool.PoolAccounts))
}

func TestListPools(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	accounts := insertAccounts(t, store, connectorID)
	insertedPools := insertPools(t, store, accounts)

	t.Run("list all pools", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListPools(
			context.Background(),
			NewListPoolsQuery(NewPaginatedQueryOptions(PoolQuery{}).WithPageSize(15)),
		)
		require.NoError(t, err)
		require.Equal(t, 2, len(cursor.Data))
		require.Equal(t, 15, cursor.PageSize)
		require.Equal(t, false, cursor.HasMore)
		require.Equal(t, "", cursor.Previous)
		require.Equal(t, "", cursor.Next)
		require.Equal(t, insertedPools[1], cursor.Data[0].ID)
		require.Equal(t, 2, len(cursor.Data[0].PoolAccounts))
		require.Equal(t, insertedPools[0], cursor.Data[1].ID)
		require.Equal(t, 1, len(cursor.Data[1].PoolAccounts))
	})

	t.Run("list all pools with page size 1", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListPools(
			context.Background(),
			NewListPoolsQuery(NewPaginatedQueryOptions(PoolQuery{}).WithPageSize(1)),
		)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.Equal(t, 1, cursor.PageSize)
		require.Equal(t, true, cursor.HasMore)
		require.Equal(t, "", cursor.Previous)
		require.Equal(t, insertedPools[1], cursor.Data[0].ID)
		require.Equal(t, 2, len(cursor.Data[0].PoolAccounts))

		var query ListPoolsQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListPools(context.Background(), query)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.Equal(t, 1, cursor.PageSize)
		require.Equal(t, false, cursor.HasMore)
		require.Equal(t, insertedPools[0], cursor.Data[0].ID)
		require.Equal(t, 1, len(cursor.Data[0].PoolAccounts))

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListPools(context.Background(), query)
		require.NoError(t, err)
		require.Equal(t, 1, len(cursor.Data))
		require.Equal(t, 1, cursor.PageSize)
		require.Equal(t, true, cursor.HasMore)
		require.Equal(t, insertedPools[1], cursor.Data[0].ID)
		require.Equal(t, 2, len(cursor.Data[0].PoolAccounts))
	})
}
