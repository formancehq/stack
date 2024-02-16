package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	acc1ID models.AccountID
	acc1T  = time.Date(2023, 11, 14, 4, 59, 0, 0, time.UTC)

	acc2ID models.AccountID
	acc2T  = time.Date(2023, 11, 14, 4, 58, 0, 0, time.UTC)

	acc3ID models.AccountID
	acc3T  = time.Date(2023, 11, 14, 4, 57, 0, 0, time.UTC)
)

func TestAccounts(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testUpdateAccounts(t, store)
	testUninstallConnectors(t, store)
	testAccountsDeletedAfterConnectorUninstall(t, store)
}

func testCreateAccounts(t *testing.T, store *storage.Storage) {
	acc1ID = models.AccountID{
		Reference:   "test1",
		ConnectorID: connectorID,
	}
	acc2ID = models.AccountID{
		Reference:   "test2",
		ConnectorID: connectorID,
	}
	acc3ID = models.AccountID{
		Reference:   "test3",
		ConnectorID: connectorID,
	}

	acc1 := &models.Account{
		ID:           acc1ID,
		CreatedAt:    acc1T,
		Reference:    "test1",
		ConnectorID:  connectorID,
		DefaultAsset: "USD",
		AccountName:  "test1",
		Type:         models.AccountTypeInternal,
		Metadata: map[string]string{
			"foo": "bar",
		},
	}

	acc2 := &models.Account{
		ID:          acc2ID,
		CreatedAt:   acc2T,
		Reference:   "test2",
		ConnectorID: connectorID,
		Type:        models.AccountTypeExternal,
	}

	acc3 := &models.Account{
		ID:          acc3ID,
		CreatedAt:   acc3T,
		Reference:   "test3",
		ConnectorID: connectorID,
		Type:        models.AccountTypeInternal,
	}

	connectorIDFail := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
	accFail := &models.Account{
		ID:          models.AccountID{Reference: "test4", ConnectorID: connectorIDFail},
		CreatedAt:   acc3T,
		Reference:   "test4",
		ConnectorID: connectorIDFail,
		Type:        models.AccountTypeInternal,
	}

	// Try to insert accounts from a not installed connector
	err := store.UpsertAccounts(
		context.Background(),
		[]*models.Account{accFail},
	)
	require.Error(t, err)

	err = store.UpsertAccounts(
		context.Background(),
		[]*models.Account{acc1, acc2, acc3},
	)
	require.NoError(t, err)

	testGetAccount(t, store, acc1.ID, acc1, false)
	testGetAccount(t, store, acc2.ID, acc2, false)
	testGetAccount(t, store, acc3.ID, acc3, false)
	testGetAccount(t, store, models.AccountID{Reference: "test4", ConnectorID: connectorID}, nil, true)
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
		ID:           acc1ID,
		CreatedAt:    time.Date(2023, 11, 14, 5, 59, 0, 0, time.UTC), // New timestamps, but should not be updated in the database
		Reference:    "test1",
		ConnectorID:  connectorID,
		DefaultAsset: "EUR",
		AccountName:  "test1-update",
		Type:         models.AccountTypeInternal,
		Metadata: map[string]string{
			"foo2": "bar2",
		},
	}

	err := store.UpsertAccounts(
		context.Background(),
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
