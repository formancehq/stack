package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/stretchr/testify/require"
)

var (
	acc1ID = models.AccountID{
		Reference: "test1",
		Provider:  models.ConnectorProviderDummyPay,
	}
	acc1T = time.Now().Add(-1 * time.Minute).UTC().Round(time.Microsecond)

	acc2ID = models.AccountID{
		Reference: "test2",
		Provider:  models.ConnectorProviderDummyPay,
	}
	acc2T = time.Now().Add(-2 * time.Minute).UTC().Round(time.Microsecond)

	acc3ID = models.AccountID{
		Reference: "test3",
		Provider:  models.ConnectorProviderDummyPay,
	}
	acc3T = time.Now().Add(-3 * time.Minute).UTC().Round(time.Microsecond)
)

func TestAccounts(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testListAccounts(t, store)
	testUpdateAccounts(t, store)
	testUninstallConnectors(t, store)
	testAccountsDeletedAfterConnectorUninstall(t, store)
}

func testCreateAccounts(t *testing.T, store *storage.Storage) {
	acc1 := &models.Account{
		ID:           acc1ID,
		CreatedAt:    acc1T,
		Reference:    "test1",
		Provider:     models.ConnectorProviderDummyPay,
		DefaultAsset: "USD",
		AccountName:  "test1",
		Type:         models.AccountTypeInternal,
		Metadata: map[string]string{
			"foo": "bar",
		},
	}

	acc2 := &models.Account{
		ID:        acc2ID,
		CreatedAt: acc2T,
		Reference: "test2",
		Provider:  models.ConnectorProviderDummyPay,
		Type:      models.AccountTypeExternal,
	}

	acc3 := &models.Account{
		ID:        acc3ID,
		CreatedAt: acc3T,
		Reference: "test3",
		Provider:  models.ConnectorProviderDummyPay,
		Type:      models.AccountTypeInternal,
	}

	// Try to insert accounts from a not installed connector
	err := store.UpsertAccounts(
		context.Background(),
		models.ConnectorProviderMangopay,
		[]*models.Account{acc1, acc2, acc3},
	)
	require.Error(t, err)

	err = store.UpsertAccounts(
		context.Background(),
		models.ConnectorProviderDummyPay,
		[]*models.Account{acc1, acc2, acc3},
	)
	require.NoError(t, err)

	testGetAccount(t, store, acc1.ID, acc1, false)
	testGetAccount(t, store, acc2.ID, acc2, false)
	testGetAccount(t, store, acc3.ID, acc3, false)
	testGetAccount(t, store, models.AccountID{Reference: "test4", Provider: models.ConnectorProviderDummyPay}, nil, true)
}

func testGetAccount(
	t *testing.T,
	store *storage.Storage,
	id models.AccountID,
	expectedAccount *models.Account,
	expectedError bool,
) {
	account, err := store.GetAccount(context.Background(), id.String())
	if expectedError {
		require.Error(t, err)
		return
	} else {
		require.NoError(t, err)
	}

	account.CreatedAt = account.CreatedAt.UTC()
	require.Equal(t, expectedAccount, account)
}

func testUpdateAccounts(t *testing.T, store *storage.Storage) {
	acc1Updated := &models.Account{
		ID: models.AccountID{
			Reference: "test1",
			Provider:  models.ConnectorProviderDummyPay,
		},
		CreatedAt:    time.Now().UTC().Round(time.Microsecond), // New timestamps, but should not be updated in the database
		Reference:    "test1",
		Provider:     models.ConnectorProviderDummyPay,
		DefaultAsset: "EUR",
		AccountName:  "test1-update",
		Type:         models.AccountTypeInternal,
		Metadata: map[string]string{
			"foo2": "bar2",
		},
	}

	err := store.UpsertAccounts(
		context.Background(),
		models.ConnectorProviderDummyPay,
		[]*models.Account{acc1Updated},
	)
	require.NoError(t, err)

	// CreatedAt should not be updated
	acc1Updated.CreatedAt = acc1T
	testGetAccount(t, store, acc1Updated.ID, acc1Updated, false)
}

func testAccountsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	// Accounts should be deleted after uninstalling the connector
	testGetAccount(t, store, acc1ID, nil, true)
	testGetAccount(t, store, acc2ID, nil, true)
	testGetAccount(t, store, acc3ID, nil, true)
}

func testListAccounts(t *testing.T, store *storage.Storage) {
	query, err := storage.Paginate(1, "", nil, nil)
	require.NoError(t, err)

	accounts, paginationDetails, err := store.ListAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, acc1ID, accounts[0].ID)

	query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	accounts, paginationDetails, err = store.ListAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, acc2ID, accounts[0].ID)

	query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	accounts, paginationDetails, err = store.ListAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, acc3ID, accounts[0].ID)

	query, err = storage.Paginate(1, paginationDetails.PreviousPage, nil, nil)
	require.NoError(t, err)

	accounts, paginationDetails, err = store.ListAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, acc2ID, accounts[0].ID)

	query, err = storage.Paginate(2, "", nil, nil)
	require.NoError(t, err)

	accounts, paginationDetails, err = store.ListAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, accounts, 2)
	require.Equal(t, acc1ID, accounts[0].ID)
	require.Equal(t, acc2ID, accounts[1].ID)

	query, err = storage.Paginate(2, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	accounts, paginationDetails, err = store.ListAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, accounts, 1)
	require.Equal(t, acc3ID, accounts[0].ID)
}
