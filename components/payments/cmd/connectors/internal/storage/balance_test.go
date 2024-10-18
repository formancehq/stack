package storage_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	b1T = time.Date(2023, 11, 14, 5, 1, 10, 0, time.UTC)
	b2T = time.Date(2023, 11, 14, 5, 1, 20, 0, time.UTC)
	b3T = time.Date(2023, 11, 14, 5, 1, 40, 0, time.UTC)
)

func TestBalances(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreateBalances(t, store)
	testUpdateBalances(t, store)
	testUninstallConnectors(t, store)
	testBalancesDeletedAfterConnectorUninstall(t, store)
}

func testCreateBalances(t *testing.T, store *storage.Storage) {
	b1 := &models.Balance{
		AccountID: models.AccountID{
			Reference:   "not_existing",
			ConnectorID: connectorID,
		},
		Asset:         "USD",
		Balance:       big.NewInt(int64(100)),
		CreatedAt:     b1T,
		LastUpdatedAt: b1T,
	}

	// Cannot insert balance for non-existing account
	err := store.InsertBalances(context.Background(), []*models.Balance{b1}, false)
	require.Error(t, err)

	// When inserting with ignore, no error is returned
	err = store.InsertBalances(context.Background(), []*models.Balance{b1}, true)
	require.NoError(t, err)

	b1.AccountID = acc1ID
	err = store.InsertBalances(context.Background(), []*models.Balance{b1}, true)
	require.NoError(t, err)

	b2 := &models.Balance{
		AccountID:     acc1ID,
		Asset:         "USD",
		Balance:       big.NewInt(int64(200)),
		CreatedAt:     b2T,
		LastUpdatedAt: b2T,
	}
	err = store.InsertBalances(context.Background(), []*models.Balance{b2}, true)
	require.NoError(t, err)

	testGetBalance(t, store, acc1ID, []*models.Balance{b2, b1})
}

func testUpdateBalances(t *testing.T, store *storage.Storage) {
	b1 := &models.Balance{
		AccountID:     acc2ID,
		Asset:         "USD",
		Balance:       big.NewInt(int64(338737362)),
		CreatedAt:     time.Date(2024, 10, 8, 22, 28, 18, 893000, time.UTC),
		LastUpdatedAt: time.Date(2024, 10, 8, 22, 28, 18, 893000, time.UTC),
	}

	err := store.InsertBalances(context.Background(), []*models.Balance{b1}, false)
	require.NoError(t, err)

	testGetBalance(t, store, acc2ID, []*models.Balance{b1})

	b2 := &models.Balance{
		AccountID:     acc2ID,
		Asset:         "USD",
		Balance:       big.NewInt(int64(317070162)),
		CreatedAt:     time.Date(2024, 10, 15, 15, 00, 01, 960000, time.UTC),
		LastUpdatedAt: time.Date(2024, 10, 15, 15, 00, 01, 960000, time.UTC),
	}

	b1.LastUpdatedAt = b2.CreatedAt
	err = store.InsertBalances(context.Background(), []*models.Balance{b1, b2}, false)
	require.NoError(t, err)

	testGetBalance(t, store, acc2ID, []*models.Balance{b2, b1})

	b3 := &models.Balance{
		AccountID:     acc2ID,
		Asset:         "USD",
		Balance:       big.NewInt(int64(327762162)),
		CreatedAt:     time.Date(2024, 10, 16, 19, 36, 29, 850000, time.UTC),
		LastUpdatedAt: time.Date(2024, 10, 16, 19, 36, 29, 850000, time.UTC),
	}

	err = store.InsertBalances(context.Background(), []*models.Balance{b1, b2, b3}, false)
	require.NoError(t, err)

	b2.LastUpdatedAt = b3.CreatedAt
	testGetBalance(t, store, acc2ID, []*models.Balance{b3, b2, b1})

	b4 := &models.Balance{
		AccountID:     acc2ID,
		Asset:         "USD",
		Balance:       big.NewInt(int64(327762162)),
		CreatedAt:     time.Date(2024, 10, 16, 19, 38, 06, 766000, time.UTC),
		LastUpdatedAt: time.Date(2024, 10, 16, 19, 38, 06, 766000, time.UTC),
	}

	err = store.InsertBalances(context.Background(), []*models.Balance{b1, b2, b3, b4}, false)
	require.NoError(t, err)

	b3.LastUpdatedAt = b3T
	testGetBalance(t, store, acc2ID, []*models.Balance{b3, b2, b1})
}

func testGetBalance(
	t *testing.T,
	store *storage.Storage,
	accountID models.AccountID,
	expectedBalances []*models.Balance,
) {
	balances, err := store.GetBalancesForAccountID(context.Background(), accountID)
	require.NoError(t, err)
	require.Len(t, balances, len(expectedBalances))
	for i := range balances {
		if i < len(balances)-1 {
			require.Equal(t, balances[i+1].LastUpdatedAt.UTC(), balances[i].CreatedAt.UTC())
		}
		require.Equal(t, expectedBalances[i].CreatedAt.UTC(), balances[i].CreatedAt.UTC())
		require.Equal(t, expectedBalances[i].AccountID, balances[i].AccountID)
		require.Equal(t, expectedBalances[i].Asset, balances[i].Asset)
		require.Equal(t, expectedBalances[i].Balance, balances[i].Balance)
	}
}

func testBalancesDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	balances, err := store.GetBalancesForAccountID(context.Background(), acc1ID)
	require.NoError(t, err)
	require.Len(t, balances, 0)
}
