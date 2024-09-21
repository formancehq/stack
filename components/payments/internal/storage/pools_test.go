package storage

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/query"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	poolID1      = uuid.New()
	poolID2      = uuid.New()
	poolID3      = uuid.New()
	defaultPools = []models.Pool{
		{
			ID:        poolID1,
			Name:      "test1",
			CreatedAt: now.Add(-60 * time.Minute).UTC().Time,
			PoolAccounts: []models.PoolAccounts{
				{
					PoolID:    poolID1,
					AccountID: defaultAccounts[0].ID,
				},
				{
					PoolID:    poolID1,
					AccountID: defaultAccounts[1].ID,
				},
			},
		},
		{
			ID:        poolID2,
			Name:      "test2",
			CreatedAt: now.Add(-30 * time.Minute).UTC().Time,
			PoolAccounts: []models.PoolAccounts{
				{
					PoolID:    poolID2,
					AccountID: defaultAccounts[2].ID,
				},
			},
		},
		{
			ID:        poolID3,
			Name:      "test3",
			CreatedAt: now.Add(-55 * time.Minute).UTC().Time,
			PoolAccounts: []models.PoolAccounts{
				{
					PoolID:    poolID3,
					AccountID: defaultAccounts[2].ID,
				},
			},
		},
	}
)

func upsertPool(t *testing.T, ctx context.Context, storage Storage, pool models.Pool) {
	require.NoError(t, storage.PoolsUpsert(ctx, pool))
}

func TestPoolsUpsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertPool(t, ctx, store, defaultPools[0])
	upsertPool(t, ctx, store, defaultPools[1])

	t.Run("upsert with same name", func(t *testing.T) {
		poolID3 := uuid.New()
		p := models.Pool{
			ID:        poolID3,
			Name:      "test1",
			CreatedAt: now.Add(-30 * time.Minute).UTC().Time,
			PoolAccounts: []models.PoolAccounts{
				{
					PoolID:    poolID3,
					AccountID: defaultAccounts[2].ID,
				},
			},
		}

		err := store.PoolsUpsert(ctx, p)
		require.Error(t, err)
	})

	t.Run("upsert with same id", func(t *testing.T) {
		upsertPool(t, ctx, store, defaultPools[1])

		actual, err := store.PoolsGet(ctx, defaultPools[1].ID)
		require.NoError(t, err)
		require.Equal(t, defaultPools[1], *actual)
	})

	t.Run("upsert with same id but more related accounts", func(t *testing.T) {
		p := defaultPools[0]
		p.PoolAccounts = append(p.PoolAccounts, models.PoolAccounts{
			PoolID:    p.ID,
			AccountID: defaultAccounts[2].ID,
		})

		upsertPool(t, ctx, store, p)

		actual, err := store.PoolsGet(ctx, defaultPools[0].ID)
		require.NoError(t, err)
		require.Equal(t, p, *actual)
	})

	t.Run("upsert with same id, but wrong related account pool id", func(t *testing.T) {
		p := defaultPools[0]
		p.PoolAccounts = append(p.PoolAccounts, models.PoolAccounts{
			PoolID:    uuid.New(),
			AccountID: defaultAccounts[2].ID,
		})

		err := store.PoolsUpsert(ctx, p)
		require.Error(t, err)
	})

	t.Run("upsert with same id, but wrong related account account id", func(t *testing.T) {
		p := defaultPools[0]
		p.PoolAccounts = append(p.PoolAccounts, models.PoolAccounts{
			PoolID: p.ID,
			AccountID: models.AccountID{
				Reference:   "unknown",
				ConnectorID: defaultConnector.ID,
			},
		})

		err := store.PoolsUpsert(ctx, p)
		require.Error(t, err)
	})
}

func TestPoolsGet(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertPool(t, ctx, store, defaultPools[0])
	upsertPool(t, ctx, store, defaultPools[1])
	upsertPool(t, ctx, store, defaultPools[2])

	t.Run("get existing pool", func(t *testing.T) {
		for _, p := range defaultPools {
			actual, err := store.PoolsGet(ctx, p.ID)
			require.NoError(t, err)
			require.Equal(t, p, *actual)
		}
	})

	t.Run("get non-existing pool", func(t *testing.T) {
		p, err := store.PoolsGet(ctx, uuid.New())
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)
		require.Nil(t, p)
	})
}

func TestPoolsDelete(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertPool(t, ctx, store, defaultPools[0])
	upsertPool(t, ctx, store, defaultPools[1])
	upsertPool(t, ctx, store, defaultPools[2])

	t.Run("delete unknown pool", func(t *testing.T) {
		require.NoError(t, store.PoolsDelete(ctx, uuid.New()))
		for _, p := range defaultPools {
			actual, err := store.PoolsGet(ctx, p.ID)
			require.NoError(t, err)
			require.Equal(t, p, *actual)
		}
	})

	t.Run("delete existing pool", func(t *testing.T) {
		require.NoError(t, store.PoolsDelete(ctx, defaultPools[0].ID))

		_, err := store.PoolsGet(ctx, defaultPools[0].ID)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotFound)

		actual, err := store.PoolsGet(ctx, defaultPools[1].ID)
		require.NoError(t, err)
		require.Equal(t, defaultPools[1], *actual)
	})
}

func TestPoolsAddAccount(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertPool(t, ctx, store, defaultPools[0])
	upsertPool(t, ctx, store, defaultPools[1])

	t.Run("add unknown account to pool", func(t *testing.T) {
		err := store.PoolsAddAccount(ctx, defaultPools[0].ID, models.AccountID{
			Reference:   "unknown",
			ConnectorID: defaultConnector.ID,
		})
		require.Error(t, err)
	})

	t.Run("add account to unknown pool", func(t *testing.T) {
		err := store.PoolsAddAccount(ctx, uuid.New(), defaultAccounts[0].ID)
		require.Error(t, err)
	})

	t.Run("add account to pool", func(t *testing.T) {
		require.NoError(t, store.PoolsAddAccount(ctx, defaultPools[0].ID, defaultAccounts[2].ID))

		p := defaultPools[0]
		p.PoolAccounts = append(p.PoolAccounts, models.PoolAccounts{
			PoolID:    p.ID,
			AccountID: defaultAccounts[2].ID,
		})

		actual, err := store.PoolsGet(ctx, defaultPools[0].ID)
		require.NoError(t, err)
		require.Equal(t, p, *actual)
	})
}

func TestPoolsRemoveAccount(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertPool(t, ctx, store, defaultPools[0])
	upsertPool(t, ctx, store, defaultPools[1])

	t.Run("remove unknown account from pool", func(t *testing.T) {
		require.NoError(t, store.PoolsRemoveAccount(ctx, defaultPools[0].ID, models.AccountID{
			Reference:   "unknown",
			ConnectorID: defaultConnector.ID,
		}))
	})

	t.Run("remove account from unknown pool", func(t *testing.T) {
		require.NoError(t, store.PoolsRemoveAccount(ctx, uuid.New(), defaultAccounts[0].ID))
	})

	t.Run("remove account from pool", func(t *testing.T) {
		require.NoError(t, store.PoolsRemoveAccount(ctx, defaultPools[0].ID, defaultAccounts[1].ID))

		p := defaultPools[0]
		p.PoolAccounts = p.PoolAccounts[:1]

		actual, err := store.PoolsGet(ctx, defaultPools[0].ID)
		require.NoError(t, err)
		require.Equal(t, p, *actual)
	})
}

func TestPoolsList(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertPool(t, ctx, store, defaultPools[0])
	upsertPool(t, ctx, store, defaultPools[1])
	upsertPool(t, ctx, store, defaultPools[2])

	t.Run("list pools by name", func(t *testing.T) {
		q := NewListPoolsQuery(
			bunpaginate.NewPaginatedQueryOptions(PoolQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("name", "test1")),
		)

		cursor, err := store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		require.Equal(t, []models.Pool{defaultPools[0]}, cursor.Data)
	})

	t.Run("list pools by unknown name", func(t *testing.T) {
		q := NewListPoolsQuery(
			bunpaginate.NewPaginatedQueryOptions(PoolQuery{}).
				WithPageSize(15).
				WithQueryBuilder(query.Match("name", "unknown")),
		)

		cursor, err := store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Empty(t, cursor.Data)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
	})

	t.Run("list pools test cursor", func(t *testing.T) {
		q := NewListPoolsQuery(
			bunpaginate.NewPaginatedQueryOptions(PoolQuery{}).
				WithPageSize(1),
		)

		cursor, err := store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, []models.Pool{defaultPools[1]}, cursor.Data)

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, []models.Pool{defaultPools[2]}, cursor.Data)

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.Empty(t, cursor.Next)
		require.Equal(t, []models.Pool{defaultPools[0]}, cursor.Data)

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.NotEmpty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, []models.Pool{defaultPools[2]}, cursor.Data)

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.PoolsList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, []models.Pool{defaultPools[1]}, cursor.Data)
	})
}
