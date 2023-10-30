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
	b1T = time.Now().UTC().Add(-50 * time.Second)

	b2T = time.Now().UTC().Add(-40 * time.Second)
)

func TestBalances(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreateBalances(t, store)
	testListBalances(t, store)
	testUninstallConnectors(t, store)
	testBalancesDeletedAfterConnectorUninstall(t, store)
}

func testCreateBalances(t *testing.T, store *storage.Storage) {
	b1 := &models.Balance{
		AccountID: models.AccountID{
			Reference: "not_existing",
			Provider:  models.ConnectorProviderDummyPay,
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
	query, err := storage.Paginate(100, "", nil, nil)
	require.NoError(t, err)

	balances, paginationDetails, err := store.ListBalances(context.Background(), storage.NewBalanceQuery(query).WithAccountID(&accountID))
	require.NoError(t, err)
	require.Len(t, balances, len(expectedBalances))
	require.False(t, paginationDetails.HasMore)
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

func testListBalances(t *testing.T, store *storage.Storage) {
	query, err := storage.Paginate(1, "", nil, nil)
	require.NoError(t, err)

	balances, paginationDetails, err := store.ListBalances(context.Background(), storage.NewBalanceQuery(query))
	require.NoError(t, err)
	require.Len(t, balances, 1)
	require.True(t, paginationDetails.HasMore)
	require.Equal(t, b2T, balances[0].CreatedAt.UTC())

	query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	balances, paginationDetails, err = store.ListBalances(context.Background(), storage.NewBalanceQuery(query))
	require.NoError(t, err)
	require.Len(t, balances, 1)
	require.False(t, paginationDetails.HasMore)
	require.Equal(t, b1T, balances[0].CreatedAt.UTC())

	query, err = storage.Paginate(1, paginationDetails.PreviousPage, nil, nil)
	require.NoError(t, err)

	balances, paginationDetails, err = store.ListBalances(context.Background(), storage.NewBalanceQuery(query))
	require.NoError(t, err)
	require.Len(t, balances, 1)
	require.True(t, paginationDetails.HasMore)
	require.Equal(t, b2T, balances[0].CreatedAt.UTC())

	query, err = storage.Paginate(2, "", nil, nil)
	require.NoError(t, err)

	balances, paginationDetails, err = store.ListBalances(context.Background(), storage.NewBalanceQuery(query))
	require.NoError(t, err)
	require.Len(t, balances, 2)
	require.False(t, paginationDetails.HasMore)
	require.Equal(t, b2T, balances[0].CreatedAt.UTC())
	require.Equal(t, b1T, balances[1].CreatedAt.UTC())
}

func testBalancesDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	query, err := storage.Paginate(1, "", nil, nil)
	require.NoError(t, err)

	balances, paginationDetails, err := store.ListBalances(context.Background(), storage.NewBalanceQuery(query))
	require.NoError(t, err)
	require.Len(t, balances, 0)
	require.False(t, paginationDetails.HasMore)

}
