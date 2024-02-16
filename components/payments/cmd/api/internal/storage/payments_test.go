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

func insertPayments(t *testing.T, store *Storage, connectorID models.ConnectorID) []models.Payment {
	p1 := models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: "test_1",
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		CreatedAt:     time.Date(2023, 11, 14, 8, 0, 0, 0, time.UTC),
		Reference:     "test_1",
		Amount:        big.NewInt(100),
		InitialAmount: big.NewInt(100),
		Type:          models.PaymentTypePayIn,
		Status:        models.PaymentStatusPending,
		Scheme:        models.PaymentSchemeA2A,
		Asset:         models.Asset("USD/2"),
		SourceAccountID: &models.AccountID{
			Reference:   "test_account",
			ConnectorID: connectorID,
		},
		DestinationAccountID: &models.AccountID{
			Reference:   "test_account2",
			ConnectorID: connectorID,
		},
	}
	_, err := store.DB().NewInsert().
		Model(&p1).
		Exec(context.Background())
	require.NoError(t, err)

	p2 := models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: "test_2",
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		CreatedAt:     time.Date(2023, 11, 14, 9, 0, 0, 0, time.UTC),
		Reference:     "test_2",
		Amount:        big.NewInt(200),
		InitialAmount: big.NewInt(100),
		Type:          models.PaymentTypePayOut,
		Status:        models.PaymentStatusPending,
		Scheme:        models.PaymentSchemeA2A,
		Asset:         models.Asset("EUR/2"),
		SourceAccountID: &models.AccountID{
			Reference:   "test_account",
			ConnectorID: connectorID,
		},
		DestinationAccountID: &models.AccountID{
			Reference:   "test_account2",
			ConnectorID: connectorID,
		},
	}
	_, err = store.DB().NewInsert().
		Model(&p2).
		Exec(context.Background())
	require.NoError(t, err)

	return []models.Payment{p1, p2}
}

func TestListPayments(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	connectorID := installConnector(t, store)
	insertAccounts(t, store, connectorID)
	payments := insertPayments(t, store, connectorID)

	t.Run("list all payments with page size 1", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListPayments(
			context.Background(),
			NewListPaymentsQuery(NewPaginatedQueryOptions(PaymentQuery{}).WithPageSize(1)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].Connector = nil
		require.Equal(t, payments[1], cursor.Data[0])

		var query ListPaymentsQuery
		err = bunpaginate.UnmarshalCursor(cursor.Next, &query)
		require.NoError(t, err)
		cursor, err = store.ListPayments(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].Connector = nil
		require.Equal(t, payments[0], cursor.Data[0])

		err = bunpaginate.UnmarshalCursor(cursor.Previous, &query)
		require.NoError(t, err)
		cursor, err = store.ListPayments(
			context.Background(),
			query,
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 1)
		require.True(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].Connector = nil
		require.Equal(t, payments[1], cursor.Data[0])
	})

	t.Run("list all payments with page size 2", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListPayments(
			context.Background(),
			NewListPaymentsQuery(NewPaginatedQueryOptions(PaymentQuery{}).WithPageSize(2)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].Connector = nil
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[1].Connector = nil
		require.Equal(t, payments[1], cursor.Data[0])
		require.Equal(t, payments[0], cursor.Data[1])
	})

	t.Run("list all payments with page size > number of payments", func(t *testing.T) {
		t.Parallel()

		cursor, err := store.ListPayments(
			context.Background(),
			NewListPaymentsQuery(NewPaginatedQueryOptions(PaymentQuery{}).WithPageSize(10)),
		)
		require.NoError(t, err)
		require.Len(t, cursor.Data, 2)
		require.False(t, cursor.HasMore)
		cursor.Data[0].CreatedAt = cursor.Data[0].CreatedAt.UTC()
		cursor.Data[0].Connector = nil
		cursor.Data[1].CreatedAt = cursor.Data[1].CreatedAt.UTC()
		cursor.Data[1].Connector = nil
		require.Equal(t, payments[1], cursor.Data[0])
		require.Equal(t, payments[0], cursor.Data[1])
	})
}
