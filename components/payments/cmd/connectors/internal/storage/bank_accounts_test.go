package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

var (
	bankAccount1ID uuid.UUID
	bankAccount2ID uuid.UUID

	bankAccount1T = time.Date(2023, 11, 14, 5, 2, 0, 0, time.UTC)
	bankAccount2T = time.Date(2023, 11, 14, 5, 1, 0, 0, time.UTC)
)

func TestBankAccounts(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreateBankAccounts(t, store)
	testUpdateBankAccountMetadata(t, store)
	testUninstallConnectors(t, store)
	testBankAccountsDeletedAfterConnectorUninstall(t, store)
}

func testCreateBankAccounts(t *testing.T, store *storage.Storage) {
	bankAccount1 := &models.BankAccount{
		CreatedAt:    bankAccount1T,
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
		Name:          "test2",
		AccountNumber: "123456789",
		Country:       "FR",
	}

	err = store.CreateBankAccount(context.Background(), bankAccount2)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil, bankAccount2.ID)
	bankAccount2ID = bankAccount2.ID

	relatedAccount := &models.BankAccountRelatedAccount{
		ID:            uuid.New(),
		CreatedAt:     bankAccount2T,
		BankAccountID: bankAccount2ID,
		ConnectorID:   connectorID,
		AccountID:     acc1ID,
	}
	err = store.AddBankAccountRelatedAccount(context.Background(), relatedAccount)
	require.NoError(t, err)
	bankAccount2.RelatedAccounts = append(bankAccount2.RelatedAccounts, relatedAccount)

	err = store.AddBankAccountRelatedAccount(context.Background(), &models.BankAccountRelatedAccount{
		ID:            uuid.New(),
		CreatedAt:     bankAccount2T,
		BankAccountID: bankAccount2ID,
		ConnectorID:   connectorID,
		AccountID: models.AccountID{
			Reference:   "not_existing",
			ConnectorID: connectorID,
		},
	})
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

	require.Equal(t, bankAccount.Country, expectedBankAccount.Country)
	require.Equal(t, bankAccount.CreatedAt.UTC(), expectedBankAccount.CreatedAt.UTC())
	require.Equal(t, bankAccount.Name, expectedBankAccount.Name)

	if expand {
		require.Equal(t, bankAccount.SwiftBicCode, expectedBankAccount.SwiftBicCode)
		require.Equal(t, bankAccount.IBAN, expectedBankAccount.IBAN)
		require.Equal(t, bankAccount.AccountNumber, expectedBankAccount.AccountNumber)
	}

	require.Len(t, bankAccount.RelatedAccounts, len(expectedBankAccount.RelatedAccounts))
	for i, adj := range bankAccount.RelatedAccounts {
		require.Equal(t, adj.BankAccountID, expectedBankAccount.RelatedAccounts[i].BankAccountID)
		require.Equal(t, adj.CreatedAt.UTC(), expectedBankAccount.RelatedAccounts[i].CreatedAt.UTC())
		require.Equal(t, adj.ConnectorID, expectedBankAccount.RelatedAccounts[i].ConnectorID)
		require.Equal(t, adj.AccountID, expectedBankAccount.RelatedAccounts[i].AccountID)
	}
}

func testUpdateBankAccountMetadata(t *testing.T, store *storage.Storage) {
	metadata := map[string]string{
		"key": "value",
	}

	err := store.UpdateBankAccountMetadata(context.Background(), bankAccount1ID, metadata)
	require.NoError(t, err)

	bankAccount, err := store.GetBankAccount(context.Background(), bankAccount1ID, false)
	require.NoError(t, err)
	require.Equal(t, metadata, bankAccount.Metadata)

	// Bank account not existing
	err = store.UpdateBankAccountMetadata(context.Background(), uuid.New(), metadata)
	require.True(t, errors.Is(err, storage.ErrNotFound))
}

func testBankAccountsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	// Connector has been uninstalled, related adjustments are deleted, but not the bank
	// accounts themselves.
	bankAccount1 := &models.BankAccount{
		CreatedAt:    bankAccount1T,
		Name:         "test1",
		IBAN:         "FR7630006000011234567890189",
		SwiftBicCode: "BNPAFRPPXXX",
		Country:      "FR",
	}

	bankAccount2 := &models.BankAccount{
		CreatedAt:     bankAccount2T,
		Name:          "test2",
		AccountNumber: "123456789",
		Country:       "FR",
	}

	testGetBankAccount(t, store, bankAccount1ID, true, bankAccount1, nil)
	testGetBankAccount(t, store, bankAccount2ID, true, bankAccount2, nil)
}
