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
	bankAccount1ID uuid.UUID
	bankAccount2ID uuid.UUID

	bankAccount1T = time.Date(2023, 11, 14, 5, 2, 0, 0, time.UTC)
	bankAccount2T = time.Date(2023, 11, 14, 5, 1, 0, 0, time.UTC)
)

func TestBankAccounts(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreateBankAccounts(t, store)
	testListBankAccounts(t, store)
	testUninstallConnectors(t, store)
	testBankAccountsDeletedAfterConnectorUninstall(t, store)
}

func testCreateBankAccounts(t *testing.T, store *storage.Storage) {
	bankAccount1 := &models.BankAccount{
		CreatedAt:    bankAccount1T,
		Provider:     models.ConnectorProviderDummyPay,
		Name:         "test1",
		IBAN:         "FR7630006000011234567890189",
		SwiftBicCode: "BNPAFRPPXXX",
		Country:      "FR",
	}

	err := store.CreateBankAccount(context.Background(), bankAccount1)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, bankAccount1.ID)
	bankAccount1ID = bankAccount1.ID

	bankAccount2 := &models.BankAccount{
		CreatedAt:     bankAccount2T,
		Provider:      models.ConnectorProviderDummyPay,
		Name:          "test2",
		AccountNumber: "123456789",
		Country:       "FR",
		AccountID:     &acc1ID,
	}

	err = store.CreateBankAccount(context.Background(), bankAccount2)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, bankAccount2.ID)
	bankAccount2ID = bankAccount2.ID

	bankAccountFail := &models.BankAccount{
		CreatedAt:     bankAccount2T,
		Provider:      models.ConnectorProviderDummyPay,
		Name:          "test2",
		AccountNumber: "123456789",
		Country:       "FR",
		AccountID: &models.AccountID{
			Reference: "not_existing",
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	err = store.CreateBankAccount(context.Background(), bankAccountFail)
	require.Error(t, err)

	testGetBankAccount(t, store, bankAccount1ID, true, bankAccount1, nil)
	testGetBankAccount(t, store, bankAccount2ID, true, bankAccount2, nil)
}

func testGetBankAccount(
	t *testing.T,
	store *storage.Storage,
	bankAccountID uuid.UUID,
	expand bool,
	expectedBankAccount *models.BankAccount,
	expectedError error,
) {
	bankAccount, err := store.GetBankAccount(context.Background(), bankAccountID, expand)
	if expectedError != nil {
		require.EqualError(t, err, expectedError.Error())
		return
	} else {
		require.NoError(t, err)
	}

	require.Equal(t, bankAccount.AccountID, expectedBankAccount.AccountID)
	require.Equal(t, bankAccount.Country, expectedBankAccount.Country)
	require.Equal(t, bankAccount.CreatedAt.UTC(), expectedBankAccount.CreatedAt.UTC())
	require.Equal(t, bankAccount.Name, expectedBankAccount.Name)
	require.Equal(t, bankAccount.Provider, expectedBankAccount.Provider)

	if expand {
		require.Equal(t, bankAccount.SwiftBicCode, expectedBankAccount.SwiftBicCode)
		require.Equal(t, bankAccount.IBAN, expectedBankAccount.IBAN)
		require.Equal(t, bankAccount.AccountNumber, expectedBankAccount.AccountNumber)
	}
}

func testListBankAccounts(t *testing.T, store *storage.Storage) {
	query, err := storage.Paginate(1, "", nil, nil)
	require.NoError(t, err)

	bankAccounts, paginationDetails, err := store.ListBankAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, bankAccounts, 1)
	require.True(t, paginationDetails.HasMore)
	require.Equal(t, bankAccount1ID, bankAccounts[0].ID)

	query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	bankAccounts, paginationDetails, err = store.ListBankAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, bankAccounts, 1)
	require.False(t, paginationDetails.HasMore)
	require.Equal(t, bankAccount2ID, bankAccounts[0].ID)

	query, err = storage.Paginate(1, paginationDetails.PreviousPage, nil, nil)
	require.NoError(t, err)

	bankAccounts, paginationDetails, err = store.ListBankAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, bankAccounts, 1)
	require.True(t, paginationDetails.HasMore)
	require.Equal(t, bankAccount1ID, bankAccounts[0].ID)

	query, err = storage.Paginate(2, "", nil, nil)
	require.NoError(t, err)

	bankAccounts, paginationDetails, err = store.ListBankAccounts(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, bankAccounts, 2)
	require.False(t, paginationDetails.HasMore)
	require.Equal(t, bankAccount1ID, bankAccounts[0].ID)
	require.Equal(t, bankAccount2ID, bankAccounts[1].ID)

}

func testBankAccountsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	testGetBankAccount(t, store, bankAccount1ID, false, nil, storage.ErrNotFound)
	testGetBankAccount(t, store, bankAccount2ID, false, nil, storage.ErrNotFound)
}
