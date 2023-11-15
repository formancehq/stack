package storage_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/stretchr/testify/require"
)

func insertBalances(t *testing.T, store *storage.Storage, accountID models.AccountID) []*models.Balance {
	b1 := &models.Balance{
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

	balances := []*models.Balance{b1, &b2, &b3, &b4}
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
	balancesPerAccountAndAssets := make(map[string]map[string][]*models.Balance)
	for _, account := range accounts {
		if balancesPerAccountAndAssets[account.String()] == nil {
			balancesPerAccountAndAssets[account.String()] = make(map[string][]*models.Balance)
		}

		balances := insertBalances(t, store, account)
		for _, balance := range balances {
			balancesPerAccountAndAssets[account.String()][balance.Asset.String()] = append(balancesPerAccountAndAssets[account.String()][balance.Asset.String()], balance)
		}
	}

	t.Run("list all balances with page size 1", func(t *testing.T) {
		query, err := storage.Paginate(1, "", nil, nil)
		require.NoError(t, err)

		balances, paginationDetails, err := store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.True(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], balances[0])

		query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
		require.NoError(t, err)
		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.True(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], balances[0])

		query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
		require.NoError(t, err)
		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.True(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[0])

		query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
		require.NoError(t, err)
		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.False(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], balances[0])

		query, err = storage.Paginate(1, paginationDetails.PreviousPage, nil, nil)
		require.NoError(t, err)
		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.True(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[0])
	})

	t.Run("list all balances with page size 2", func(t *testing.T) {
		query, err := storage.Paginate(2, "", nil, nil)
		require.NoError(t, err)

		balances, paginationDetails, err := store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 2)
		require.True(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[1].CreatedAt = balances[1].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		balances[1].LastUpdatedAt = balances[1].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], balances[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], balances[1])

		query, err = storage.Paginate(2, paginationDetails.NextPage, nil, nil)
		require.NoError(t, err)
		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 2)
		require.False(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[1].CreatedAt = balances[1].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		balances[1].LastUpdatedAt = balances[1].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], balances[1])

		query, err = storage.Paginate(2, paginationDetails.PreviousPage, nil, nil)
		require.NoError(t, err)
		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]),
		)
		require.NoError(t, err)
		require.Len(t, balances, 2)
		require.True(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[1].CreatedAt = balances[1].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		balances[1].LastUpdatedAt = balances[1].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], balances[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], balances[1])
	})

	t.Run("list balances for asset", func(t *testing.T) {
		query, err := storage.Paginate(15, "", nil, nil)
		require.NoError(t, err)

		balances, paginationDetails, err := store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]).WithCurrency("USD/2"),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.False(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[0])
	})

	t.Run("list balances for asset and limit", func(t *testing.T) {
		query, err := storage.Paginate(15, "", nil, nil)
		require.NoError(t, err)

		balances, paginationDetails, err := store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).WithAccountID(&accounts[0]).WithCurrency("EUR/2").WithLimit(1),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.False(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][2], balances[0])
	})

	t.Run("list balances for asset and time range", func(t *testing.T) {
		query, err := storage.Paginate(15, "", nil, nil)
		require.NoError(t, err)

		balances, paginationDetails, err := store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).
				WithAccountID(&accounts[0]).
				WithFrom(time.Date(2023, 11, 14, 10, 15, 0, 0, time.UTC)).
				WithTo(time.Date(2023, 11, 14, 11, 15, 0, 0, time.UTC)),
		)
		require.NoError(t, err)
		require.Len(t, balances, 3)
		require.False(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		balances[1].CreatedAt = balances[1].CreatedAt.UTC()
		balances[1].LastUpdatedAt = balances[1].LastUpdatedAt.UTC()
		balances[2].CreatedAt = balances[2].CreatedAt.UTC()
		balances[2].LastUpdatedAt = balances[2].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], balances[0])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["USD/2"][0], balances[1])
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][0], balances[2])
	})

	t.Run("get balances at a precise time", func(t *testing.T) {
		query, err := storage.Paginate(15, "", nil, nil)
		require.NoError(t, err)

		balances, paginationDetails, err := store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).
				WithAccountID(&accounts[0]).
				WithCurrency("EUR/2").
				WithTo(time.Date(2023, 11, 14, 11, 15, 0, 0, time.UTC)).
				WithLimit(1),
		)
		require.NoError(t, err)
		require.Len(t, balances, 1)
		require.False(t, paginationDetails.HasMore)
		balances[0].CreatedAt = balances[0].CreatedAt.UTC()
		balances[0].LastUpdatedAt = balances[0].LastUpdatedAt.UTC()
		require.Equal(t, balancesPerAccountAndAssets[accounts[0].String()]["EUR/2"][1], balances[0])

		balances, paginationDetails, err = store.ListBalances(
			context.Background(),
			storage.NewBalanceQuery(query).
				WithAccountID(&accounts[0]).
				WithCurrency("EUR/2").
				WithTo(time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC)).
				WithLimit(1),
		)
		require.NoError(t, err)
		require.Len(t, balances, 0)
		require.False(t, paginationDetails.HasMore)
	})
}
