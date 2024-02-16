package storage

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"
	"github.com/stretchr/testify/require"
)

func insertTransferInitiation(t *testing.T, store *Storage, connectorID models.ConnectorID) []models.TransferInitiation {
	tf1 := models.TransferInitiation{
		ID: models.TransferInitiationID{
			Reference:   "tf_1",
			ConnectorID: connectorID,
		},
		CreatedAt:   time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		ScheduledAt: time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		Description: "test_1",
		Type:        models.TransferInitiationTypeTransfer,
		SourceAccountID: &models.AccountID{
			Reference:   "test_account",
			ConnectorID: connectorID,
		},
		DestinationAccountID: models.AccountID{
			Reference:   "test_account2",
			ConnectorID: connectorID,
		},
		Provider:      models.ConnectorProviderDummyPay,
		ConnectorID:   connectorID,
		Amount:        big.NewInt(100),
		InitialAmount: big.NewInt(100),
		Asset:         "EUR/2",
		Metadata: map[string]string{
			"foo": "bar",
		},
	}
	_, err := store.DB().NewInsert().
		Model(&tf1).
		Exec(context.Background())
	require.NoError(t, err)

	tf2 := models.TransferInitiation{
		ID: models.TransferInitiationID{
			Reference:   "tf_2",
			ConnectorID: connectorID,
		},
		CreatedAt:   time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
		ScheduledAt: time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
		Description: "test_2",
		Type:        models.TransferInitiationTypeTransfer,
		SourceAccountID: &models.AccountID{
			Reference:   "test_account",
			ConnectorID: connectorID,
		},
		DestinationAccountID: models.AccountID{
			Reference:   "test_account2",
			ConnectorID: connectorID,
		},
		Provider:      models.ConnectorProviderDummyPay,
		ConnectorID:   connectorID,
		Amount:        big.NewInt(100),
		InitialAmount: big.NewInt(100),
		Asset:         "USD/2",
		Metadata: map[string]string{
			"foo2": "bar2",
		},
	}

	_, err = store.DB().NewInsert().
		Model(&tf2).
		Exec(context.Background())
	require.NoError(t, err)

	return []models.TransferInitiation{tf1, tf2}
}

func TestListTransferInitiation(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	insertAccounts(t, store, connectorID)
	tfs := insertTransferInitiation(t, store, connectorID)

	t.Run("list all transfer initiations with page size 1", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListTransferInitiations(
			context.Background(),
			NewListTransferInitiationsQuery(NewPaginatedQueryOptions(TransferInitiationQuery{}).WithPageSize(1)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].ScheduledAt = cursor.Data[0].ScheduledAt.UTC()
		require.Equal(t, tfs[1], cursor.Data[0])

		var query ListTransferInitiationsQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListTransferInitiations(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].ScheduledAt = cursor.Data[0].ScheduledAt.UTC()
		require.Equal(t, tfs[0], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListTransferInitiations(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].ScheduledAt = cursor.Data[0].ScheduledAt.UTC()
		require.Equal(t, tfs[1], cursor.Data[0])
	})

	t.Run("list all transfer initiations with page size 2", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListTransferInitiations(
			context.Background(),
			NewListTransferInitiationsQuery(NewPaginatedQueryOptions(TransferInitiationQuery{}).WithPageSize(2)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].ScheduledAt = cursor.Data[0].ScheduledAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[1].ScheduledAt = cursor.Data[1].ScheduledAt.UTC()
		require.Equal(t, tfs[1], cursor.Data[0])
		require.Equal(t, tfs[0], cursor.Data[1])
	})

	t.Run("list all transfer initiations with page size > number of transfer initiations", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListTransferInitiations(
			context.Background(),
			NewListTransferInitiationsQuery(NewPaginatedQueryOptions(TransferInitiationQuery{}).WithPageSize(10)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].ScheduledAt = cursor.Data[0].ScheduledAt.UTC()
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[1].ScheduledAt = cursor.Data[1].ScheduledAt.UTC()
		require.Equal(t, tfs[1], cursor.Data[0])
		require.Equal(t, tfs[0], cursor.Data[1])
	})
}
