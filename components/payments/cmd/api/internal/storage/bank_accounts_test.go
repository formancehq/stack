package storage

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func insertBankAccounts(t *testing.T, store *Storage, connectorID models.ConnectorID) []models.BankAccount {
	acc1 := models.BankAccount{
		ID:        uuid.New(),
		CreatedAt: time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		Name:      "test_1",
		IBAN:      "FR7630006000011234567890189",
		Country:   "FR",
		Metadata: map[string]string{
			"foo": "bar",
		},
	}
	_, err := store.DB().NewInsert().
		Model(&acc1).
		Exec(context.Background())
	require.NoError(t, err)

	acc2 := models.BankAccount{
		ID:        uuid.New(),
		CreatedAt: time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
		Name:      "test_2",
		IBAN:      "FR7630006000011234567891234",
		Country:   "GB",
		Metadata: map[string]string{
			"foo2": "bar2",
		},
	}
	_, err = store.DB().NewInsert().
		Model(&acc2).
		Exec(context.Background())
	require.NoError(t, err)

	return []models.BankAccount{acc1, acc2}
}

func TestListBankAccounts(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	bankAccounts := insertBankAccounts(t, store, connectorID)

	for i := range bankAccounts {
		bankAccounts[i].CreatedAt = bankAccounts[i].CreatedAt.UTC()
		// The listing of bank accounts does not sent the IBAN info
		bankAccounts[i].IBAN = ""
	}

	t.Run("list all bank accounts with page size 1", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBankAccounts(
			context.Background(),
			NewListBankAccountQuery(NewPaginatedQueryOptions(BankAccountQuery{}).WithPageSize(1)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		require.Equal(t, bankAccounts[1], cursor.Data[0])

		var query ListBankAccountQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListBankAccounts(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		require.Equal(t, bankAccounts[0], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListBankAccounts(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		require.Equal(t, bankAccounts[1], cursor.Data[0])
	})

	t.Run("list all bank accounts with page size 2", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBankAccounts(
			context.Background(),
			NewListBankAccountQuery(NewPaginatedQueryOptions(BankAccountQuery{}).WithPageSize(2)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		require.Equal(t, bankAccounts[1], cursor.Data[0])
		require.Equal(t, bankAccounts[0], cursor.Data[1])
	})

	t.Run("list all bank accounts with page size > number of bank accounts", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListBankAccounts(
			context.Background(),
			NewListBankAccountQuery(NewPaginatedQueryOptions(BankAccountQuery{}).WithPageSize(10)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		require.Equal(t, bankAccounts[1], cursor.Data[0])
		require.Equal(t, bankAccounts[0], cursor.Data[1])
	})
}
