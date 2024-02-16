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
)

func TestBalances(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreateBalances(t, store)
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

	testGetBalance(t, store, acc1ID, []*models.Balance{b2, b1}, nil)
}

func testGetBalance(
	t *testing.T,
	store *storage.Storage,
	accountID models.AccountID,
	expectedBalances []*models.Balance,
	expectedError error,
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
