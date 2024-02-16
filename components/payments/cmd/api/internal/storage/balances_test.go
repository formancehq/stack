package storage

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/stretchr/testify/require"
)

func insertBalances(t *testing.T, store *Storage, accountID models.AccountID) []models.Balance {
	b1 := models.Balance{
		AccountID:     accountID,
		Asset:         "EUR/2",
		Balance:       big.NewInt(100),
		CreatedAt:     time.Date(2023, 11, 14, 10, 0, 0, 0, time.UTC),
		LastUpdatedAt: time.Date(2023, 11, 14, 11, 0, 0, 0, time.UTC),
	}

	b2 := models.Balance{
		AccountID:     accountID,
		Asset:         "EUR/2",
		Balance:       big.NewInt(200),
		CreatedAt:     time.Date(2023, 11, 14, 11, 0, 0, 0, time.UTC),
		LastUpdatedAt: time.Date(2023, 11, 14, 11, 30, 0, 0, time.UTC),
	}

	b3 := models.Balance{
		AccountID:     accountID,
		Asset:         "EUR/2",
		Balance:       big.NewInt(150),
		CreatedAt:     time.Date(2023, 11, 14, 11, 30, 0, 0, time.UTC),
		LastUpdatedAt: time.Date(2023, 11, 14, 11, 45, 0, 0, time.UTC),
	}

	b4 := models.Balance{
		AccountID:     accountID,
		Asset:         "USD/2",
		Balance:       big.NewInt(1000),
		CreatedAt:     time.Date(2023, 11, 14, 10, 30, 0, 0, time.UTC),
		LastUpdatedAt: time.Date(2023, 11, 14, 12, 0, 0, 0, time.UTC),
	}

	balances := []models.Balance{b1, b2, b3, b4}
	_, err := store.DB().NewInsert().
		Model(&balances).
		Exec(context.Background())
	require.NoError(t, err)

	return balances
}

func TestListBalances(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	accounts := insertAccounts(t, store, connectorID)
	balancesPerAccountAndAssets := make(map[string]map[string][]models.Balance)
	for _, account := range accounts {
		if balancesPerAccountAndAssets[account.String()] == nil {
			balancesPerAccountAndAssets[account.String()] = make(map[string][]models.Balance)
		}

		balances := insertBalances(t, store, account)
		for _, balance := range balances {
			balancesPerAccountAndAssets[account.String()][balance.Asset.String()] = append(balancesPerAccountAndAssets[account.String()][balance.Asset.String()], balance)
		}
	}

	t.Run("list all balances with page size 1", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBalances(
			context.Background(),
			NewListBalancesQuery(NewPaginatedQueryOptions(NewBalanceQuery().WithAccountID(&accounts[0])).WithPageSize(1)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], cursor.Data[0])

		var query ListBalancesQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListBalances(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListBalances(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListBalances(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListBalances(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], cursor.Data[0])
	})

	t.Run("list all balances with page size 2", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBalances(
			context.Background(),
			NewListBalancesQuery(NewPaginatedQueryOptions(NewBalanceQuery().WithAccountID(&accounts[0])).WithPageSize(2)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		cursor.Data[1].LastUpdatedAt = cursor.Data[1].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], cursor.Data[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], cursor.Data[1])

		var query ListBalancesQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListBalances(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		cursor.Data[1].LastUpdatedAt = cursor.Data[1].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], cursor.Data[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], cursor.Data[1])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListBalances(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		cursor.Data[1].LastUpdatedAt = cursor.Data[1].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], cursor.Data[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], cursor.Data[1])
	})

	t.Run("list balances for asset", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBalances(
			context.Background(),
			NewListBalancesQuery(NewPaginatedQueryOptions(NewBalanceQuery().WithAccountID(&accounts[0]).WithCurrency("USD/2")).WithPageSize(15)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], cursor.Data[0])
	})

	t.Run("list balances for asset and limit", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBalances(
			context.Background(),
			NewListBalancesQuery(
				NewPaginatedQueryOptions(
					NewBalanceQuery().
						WithAccountID(&accounts[0]).
						WithCurrency("EUR/2"),
				).
					WithPageSize(1),
			),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], cursor.Data[0])
	})

	t.Run("list balances for asset and time range", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBalances(
			context.Background(),
			NewListBalancesQuery(
				NewPaginatedQueryOptions(
					NewBalanceQuery().
						WithAccountID(&accounts[0]).
						WithFrom(time.Date(2023, 11, 14, 10, 15, 0, 0, time.UTC)).
						WithTo(time.Date(2023, 11, 14, 11, 15, 0, 0, time.UTC)),
				).
					WithPageSize(15),
			),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 3)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[1].LastUpdatedAt = cursor.Data[1].LastUpdatedAt.UTC()
		cursor.Data[2].CreatedAt = cursor.Data[2].CreatedAt.UTC()
		cursor.Data[2].LastUpdatedAt = cursor.Data[2].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], cursor.Data[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], cursor.Data[1])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], cursor.Data[2])
	})

	t.Run("get balances at a precise time", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBalances(
			context.Background(),
			NewListBalancesQuery(
				NewPaginatedQueryOptions(
					NewBalanceQuery().
						WithAccountID(&accounts[0]).
						WithCurrency("EUR/2").
						WithTo(time.Date(2023, 11, 14, 11, 15, 0, 0, time.UTC)),
				).WithPageSize(1),
			),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].LastUpdatedAt = cursor.Data[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], cursor.Data[0])

		cursor, err = store.ListBalances(
			context.Background(),
			NewListBalancesQuery(
				NewPaginatedQueryOptions(
					NewBalanceQuery().
						WithAccountID(&accounts[0]).
						WithCurrency("EUR/2").
						WithTo(time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC)),
				).WithPageSize(1),
			),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 0)
		require.False(t, cursor.HasMore)
	})
}

func TestGetBalanceAt(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	accounts := insertAccounts(t, store, connectorID)
	balancesPerAccountAndAssets := make(map[string]map[string][]models.Balance)
	for _, account := range accounts {
		if balancesPerAccountAndAssets[account.String()] == nil {
			balancesPerAccountAndAssets[account.String()] = make(map[string][]models.Balance)
		}

		balances := insertBalances(t, store, account)
		for _, balance := range balances {
			balancesPerAccountAndAssets[account.String()][balance.Asset.String()] = append(balancesPerAccountAndAssets[account.String()][balance.Asset.String()], balance)
		}
	}

	// Should have only one EUR/2 balance of 100
	t.Run("get balance at 10:15", func(t *testing.T) {
		balances, err := store.GetBalancesAt(context.Background(), accounts[0], time.Date(2023, 11, 14, 10, 15, 0, 0, time.UTC))
		require.NoError(t, err)
		require.Len(t, balances, 1)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, &balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], balances[0])
		require.Equal(t, big.NewInt(100), balances[0].Balance)
	})

	t.Run("get balance at 11:15", func(t *testing.T) {
		balances, err := store.GetBalancesAt(context.Background(), accounts[0], time.Date(2023, 11, 14, 11, 15, 0, 0, time.UTC))
		require.NoError(t, err)
		require.Len(t, balances, 2)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, &balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], balances[0])
		require.Equal(t, big.NewInt(200), balances[0].Balance)
		balances[1].CreatedAt = balances[1].CreatedAt.UTC()
		balances[1].LastUpdatedAt = balances[1].LastUpdatedAt.UTC()
		require.Equal(t, &balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[1])
		require.Equal(t, big.NewInt(1000), balances[1].Balance)
	})

	t.Run("get balance at 11:45", func(t *testing.T) {
		balances, err := store.GetBalancesAt(context.Background(), accounts[0], time.Date(2023, 11, 14, 11, 45, 0, 0, time.UTC))
		require.NoError(t, err)
		require.Len(t, balances, 2)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, &balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], balances[0])
		require.Equal(t, big.NewInt(150), balances[0].Balance)
		balances[1].CreatedAt = balances[1].CreatedAt.UTC()
		balances[1].LastUpdatedAt = balances[1].LastUpdatedAt.UTC()
		require.Equal(t, &balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[1])
		require.Equal(t, big.NewInt(1000), balances[1].Balance)
	})

	t.Run("get balance at 12:00", func(t *testing.T) {
		balances, err := store.GetBalancesAt(context.Background(), accounts[0], time.Date(2023, 11, 14, 12, 0, 0, 0, time.UTC))
		require.NoError(t, err)
		require.Len(t, balances, 1)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, &balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[0])
		require.Equal(t, big.NewInt(1000), balances[0].Balance)
	})
}
