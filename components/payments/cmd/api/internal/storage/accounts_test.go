package storage

import (
	"context"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/stretchr/testify/require"
)

func insertAccounts(t *testing.T, store *Storage, connectorID models.ConnectorID) []models.AccountID {
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

func TestListAccounts(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	insertAccounts(t, store, connectorID)

	acc1 := models.Account{
		ID: models.AccountID{
			Reference:   "test_account",
			ConnectorID: connectorID,
		},
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		Reference:   "test_account",
		AccountName: "test",
		Type:        models.AccountTypeInternal,
	}

	acc2 := models.Account{
		ID: models.AccountID{
			Reference:   "test_account2",
			ConnectorID: connectorID,
		},
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
		Reference:   "test_account2",
		AccountName: "test2",
		Type:        models.AccountTypeInternal,
	}

	t.Run("list all accounts with page size 1", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListAccounts(
			context.Background(),
			NewListAccountsQuery(NewPaginatedQueryOptions(AccountQuery{}).WithPageSize(1)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		require.Equal(t, acc2, cursor.Data[0])

		var query ListAccountsQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListAccounts(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		require.Equal(t, acc1, cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListAccounts(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		require.Equal(t, acc2, cursor.Data[0])
	})

	t.Run("list all accounts with page size 2", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListAccounts(
			context.Background(),
			NewListAccountsQuery(NewPaginatedQueryOptions(AccountQuery{}).WithPageSize(2)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		require.Equal(t, acc2, cursor.Data[0])
		require.Equal(t, acc1, cursor.Data[1])
	})

	t.Run("list all accounts with page size > number of accounts", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListAccounts(
			context.Background(),
			NewListAccountsQuery(NewPaginatedQueryOptions(AccountQuery{}).WithPageSize(10)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		require.Equal(t, acc2, cursor.Data[0])
		require.Equal(t, acc1, cursor.Data[1])
	})
}
