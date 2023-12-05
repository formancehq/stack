package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/stretchr/testify/require"
)

func insertAccounts(t *testing.T, store *storage.Storage, connectorID models.ConnectorID) []models.AccountID {
	id1 := models.AccountID{
		Reference:   "test_account",
		ConnectorID: connectorID,
	}
	acc1 := models.Account{
		ID:          id1,
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		Reference:   "test_account",
		AccountName: "test",
		Type:        models.AccountTypeInternal,
	}

	_, err := store.DB().NewInsert().
		Model(&acc1).
		Exec(context.Background())
	require.NoError(t, err)

	id2 := models.AccountID{
		Reference:   "test_account2",
		ConnectorID: connectorID,
	}
	acc2 := models.Account{
		ID:          id2,
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
		Reference:   "test_account2",
		AccountName: "test2",
		Type:        models.AccountTypeInternal,
	}

	_, err = store.DB().NewInsert().
		Model(&acc2).
		Exec(context.Background())
	require.NoError(t, err)

	return []models.AccountID{id1, id2}
}
