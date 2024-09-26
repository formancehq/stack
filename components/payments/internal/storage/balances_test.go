package storage

import (
	"context"
	"math/big"
	"testing"

	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/pointer"
	"github.com/formancehq/go-libs/time"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	defaultBalances = []models.Balance{
		{
			AccountID:     defaultAccounts[0].ID,
			CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-60 * time.Minute).UTC().Time,
			Asset:         "USD/2",
			Balance:       big.NewInt(100),
		},
		{
			AccountID:     defaultAccounts[1].ID,
			CreatedAt:     now.Add(-30 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-30 * time.Minute).UTC().Time,
			Asset:         "EUR/2",
			Balance:       big.NewInt(1000),
		},
		{
			AccountID:     defaultAccounts[0].ID,
			CreatedAt:     now.Add(-55 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-55 * time.Minute).UTC().Time,
			Asset:         "EUR/2",
			Balance:       big.NewInt(150),
		},
	}

	defaultBalances2 = []models.Balance{
		{
			AccountID:     defaultAccounts[2].ID,
			CreatedAt:     now.Add(-59 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-59 * time.Minute).UTC().Time,
			Asset:         "USD/2",
			Balance:       big.NewInt(100),
		},
		{
			AccountID:     defaultAccounts[2].ID,
			CreatedAt:     now.Add(-31 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-31 * time.Minute).UTC().Time,
			Asset:         "DKK/2",
			Balance:       big.NewInt(1000),
		},
	}
)

func upsertBalances(t *testing.T, ctx context.Context, storage Storage, balances []models.Balance) {
	require.NoError(t, storage.BalancesUpsert(ctx, balances))
}

func TestBalancesUpsert(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBalances(t, ctx, store, defaultBalances)
	upsertBalances(t, ctx, store, defaultBalances2)

	t.Run("insert balances with same asset and same balance", func(t *testing.T) {
		b := models.Balance{
			AccountID:     defaultAccounts[2].ID,
			CreatedAt:     now.Add(-20 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-20 * time.Minute).UTC().Time,
			Asset:         "USD/2",
			Balance:       big.NewInt(100),
		}

		upsertBalances(t, ctx, store, []models.Balance{b})

		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				AccountID: pointer.For(defaultAccounts[2].ID),
				Asset:     "USD/2",
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-59 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-20 * time.Minute).UTC().Time, // Last updated at should be updated to the new balance value
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		balances, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, balances.Data, 1)
		require.Equal(t, expectedBalances, balances.Data)
	})

	t.Run("insert balances same asset different balance", func(t *testing.T) {
		b := models.Balance{
			AccountID:     defaultAccounts[0].ID,
			CreatedAt:     now.Add(-20 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-20 * time.Minute).UTC().Time,
			Asset:         "USD/2",
			Balance:       big.NewInt(200),
		}

		upsertBalances(t, ctx, store, []models.Balance{b})

		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				AccountID: pointer.For(defaultAccounts[0].ID),
				Asset:     "USD/2",
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			// We should have one more balance with the new balance value
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-20 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-20 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(200),
			},
			// and the old balance should have its updated at to the new balance created at
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-20 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		balances, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, balances.Data, 2)
		require.Equal(t, expectedBalances, balances.Data)
	})

	t.Run("insert balances with new asset", func(t *testing.T) {
		b := models.Balance{
			AccountID:     defaultAccounts[2].ID,
			CreatedAt:     now.Add(-10 * time.Minute).UTC().Time,
			LastUpdatedAt: now.Add(-10 * time.Minute).UTC().Time,
			Asset:         "EUR/2",
			Balance:       big.NewInt(200),
		}

		upsertBalances(t, ctx, store, []models.Balance{b})

		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				AccountID: pointer.For(defaultAccounts[2].ID),
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-10 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-10 * time.Minute).UTC().Time,
				Asset:         "EUR/2",
				Balance:       big.NewInt(200),
			},
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-31 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-31 * time.Minute).UTC().Time,
				Asset:         "DKK/2",
				Balance:       big.NewInt(1000),
			},
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-59 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-20 * time.Minute).UTC().Time, // Because on the first function it was modified
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 3)
		require.Equal(t, expectedBalances, cursor.Data)
	})
}

func TestBalancesDeleteFromConnectorID(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBalances(t, ctx, store, defaultBalances)
	upsertBalances(t, ctx, store, defaultBalances2)

	t.Run("delete balances from unknown connector id", func(t *testing.T) {
		err := store.BalancesDeleteFromConnectorID(ctx, models.ConnectorID{
			Reference: uuid.New(),
			Provider:  "unknown",
		})
		require.NoError(t, err)
	})

	t.Run("delete balances from known connector id", func(t *testing.T) {
		err := store.BalancesDeleteFromConnectorID(ctx, defaultConnector.ID)
		require.NoError(t, err)

		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				AccountID: pointer.For(defaultAccounts[0].ID),
			}).WithPageSize(15),
		)

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
	})
}

func TestBalancesList(t *testing.T) {
	t.Parallel()

	ctx := logging.TestingContext()
	store := newStore(t)

	upsertConnector(t, ctx, store, defaultConnector)
	upsertAccounts(t, ctx, store, defaultAccounts)
	upsertBalances(t, ctx, store, defaultBalances)
	upsertBalances(t, ctx, store, defaultBalances2)

	t.Run("list balances with account id", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				AccountID: pointer.For(defaultAccounts[0].ID),
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-55 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-55 * time.Minute).UTC().Time,
				Asset:         "EUR/2",
				Balance:       big.NewInt(150),
			},
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-60 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		require.Equal(t, expectedBalances, cursor.Data)
	})

	t.Run("list balances with asset 1", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				Asset: "USD/2",
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-59 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-59 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-60 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		require.Equal(t, expectedBalances, cursor.Data)
	})

	t.Run("list balances with asset 2", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				Asset: "DKK/2",
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-31 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-31 * time.Minute).UTC().Time,
				Asset:         "DKK/2",
				Balance:       big.NewInt(1000),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.Equal(t, expectedBalances, cursor.Data)
	})

	t.Run("list balances with from", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				From: now.Add(-40 * time.Minute).UTC().Time,
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[1].ID,
				CreatedAt:     now.Add(-30 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-30 * time.Minute).UTC().Time,
				Asset:         "EUR/2",
				Balance:       big.NewInt(1000),
			},
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-31 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-31 * time.Minute).UTC().Time,
				Asset:         "DKK/2",
				Balance:       big.NewInt(1000),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		require.Equal(t, expectedBalances, cursor.Data)
	})

	t.Run("list balances with from 2", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				From: now.Add(-20 * time.Minute).UTC().Time,
			}).WithPageSize(15),
		)

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
		require.False(t, cursor.HasMore)
	})

	t.Run("list balances with to", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				To: now.Add(-40 * time.Minute).UTC().Time,
			}).WithPageSize(15),
		)

		expectedBalances := []models.Balance{
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-55 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-55 * time.Minute).UTC().Time,
				Asset:         "EUR/2",
				Balance:       big.NewInt(150),
			},
			{
				AccountID:     defaultAccounts[2].ID,
				CreatedAt:     now.Add(-59 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-59 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-60 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 3)
		require.False(t, cursor.HasMore)
		require.Equal(t, expectedBalances, cursor.Data)
	})

	t.Run("list balances with to 2", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				To: now.Add(-70 * time.Minute).UTC().Time,
			}).WithPageSize(15),
		)

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
		require.False(t, cursor.HasMore)
	})

	t.Run("list balances test cursor", func(t *testing.T) {
		q := NewListBalancesQuery(
			bunpaginate.NewPaginatedQueryOptions(BalanceQuery{
				AccountID: pointer.For(defaultAccounts[0].ID),
			}).WithPageSize(1),
		)

		expectedBalances1 := []models.Balance{
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-55 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-55 * time.Minute).UTC().Time,
				Asset:         "EUR/2",
				Balance:       big.NewInt(150),
			},
		}

		cursor, err := store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, expectedBalances1, cursor.Data)

		expectedBalances2 := []models.Balance{
			{
				AccountID:     defaultAccounts[0].ID,
				CreatedAt:     now.Add(-60 * time.Minute).UTC().Time,
				LastUpdatedAt: now.Add(-60 * time.Minute).UTC().Time,
				Asset:         "USD/2",
				Balance:       big.NewInt(100),
			},
		}

		err = bunpaginate.UnmarshalCursor(cursor.Next, &q)
		require.NoError(t, err)
		cursor, err = store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		require.Empty(t, cursor.Next)
		require.NotEmpty(t, cursor.Previous)
		require.Equal(t, expectedBalances2, cursor.Data)

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &q)
		require.NoError(t, err)
		cursor, err = store.BalancesList(ctx, q)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		require.Empty(t, cursor.Previous)
		require.NotEmpty(t, cursor.Next)
		require.Equal(t, expectedBalances1, cursor.Data)
	})
}
