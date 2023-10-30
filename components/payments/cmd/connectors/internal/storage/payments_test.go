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
	p1ID = &models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "test1",
			Type:      models.PaymentTypePayOut,
		},
		Provider: models.ConnectorProviderDummyPay,
	}
	p1T = time.Now().Add(-1 * time.Minute).UTC()
	p1  = &models.Payment{
		ID:        *p1ID,
		CreatedAt: p1T,
		Reference: "ref1",
		Amount:    big.NewInt(100),
		Type:      models.PaymentTypePayOut,
		Status:    models.PaymentStatusSucceeded,
		Scheme:    models.PaymentSchemeCardVisa,
		Asset:     models.Asset("USD/2"),
	}

	p2ID = &models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "test2",
			Type:      models.PaymentTypeTransfer,
		},
		Provider: models.ConnectorProviderDummyPay,
	}
	p2T = time.Now().Add(-2 * time.Minute).UTC()
	p2  = &models.Payment{
		ID:        *p2ID,
		CreatedAt: p2T,
		Reference: "ref2",
		Amount:    big.NewInt(150),
		Type:      models.PaymentTypePayIn,
		Status:    models.PaymentStatusFailed,
		Scheme:    models.PaymentSchemeCardVisa,
		Asset:     models.Asset("EUR/2"),
	}
)

func TestPayments(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreatePayments(t, store)
	testListPayments(t, store)
	testUpdatePayment(t, store)
	testUninstallConnectors(t, store)
	testPaymentsDeletedAfterConnectorUninstall(t, store)
}

func testCreatePayments(t *testing.T, store *storage.Storage) {
	pFail := &models.Payment{
		ID:        *p1ID,
		CreatedAt: p1T,
		Reference: "ref1",
		Amount:    big.NewInt(100),
		Type:      models.PaymentTypePayOut,
		Status:    models.PaymentStatusSucceeded,
		Scheme:    models.PaymentSchemeCardVisa,
		Asset:     models.Asset("USD/2"),
		SourceAccountID: &models.AccountID{
			Reference: "not_existing",
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	err := store.UpsertPayments(context.Background(), models.ConnectorProviderDummyPay, []*models.Payment{pFail})
	require.Error(t, err)

	err = store.UpsertPayments(context.Background(), models.ConnectorProviderDummyPay, []*models.Payment{p1})
	require.NoError(t, err)

	err = store.UpsertPayments(context.Background(), models.ConnectorProviderDummyPay, []*models.Payment{p2})
	require.NoError(t, err)

	testGetPayment(t, store, *p1ID, p1, nil)
	testGetPayment(t, store, *p2ID, p2, nil)
}

func testGetPayment(
	t *testing.T,
	store *storage.Storage,
	paymentID models.PaymentID,
	expected *models.Payment,
	expectedErr error,
) {
	payment, err := store.GetPayment(context.Background(), paymentID.String())
	if expectedErr != nil {
		require.EqualError(t, err, expectedErr.Error())
		return
	} else {
		require.NoError(t, err)
	}

	payment.CreatedAt = payment.CreatedAt.UTC()
	checkPaymentsEqual(t, expected, payment)
}

func checkPaymentsEqual(t *testing.T, p1, p2 *models.Payment) {
	require.Equal(t, p1.ID, p2.ID)
	require.Equal(t, p1.CreatedAt.UTC(), p2.CreatedAt.UTC())
	require.Equal(t, p1.Reference, p2.Reference)
	require.Equal(t, p1.Amount, p2.Amount)
	require.Equal(t, p1.Type, p2.Type)
	require.Equal(t, p1.Status, p2.Status)
	require.Equal(t, p1.Scheme, p2.Scheme)
	require.Equal(t, p1.Asset, p2.Asset)
	require.Equal(t, p1.SourceAccountID, p2.SourceAccountID)
	require.Equal(t, p1.DestinationAccountID, p2.DestinationAccountID)
	require.Equal(t, p1.RawData, p2.RawData)
}

func testUpdatePayment(t *testing.T, store *storage.Storage) {
	p1.CreatedAt = time.Now().UTC()
	p1.Reference = "ref1"
	p1.Amount = big.NewInt(150)
	p1.Type = models.PaymentTypePayIn
	p1.Status = models.PaymentStatusPending
	p1.Scheme = models.PaymentSchemeCardVisa
	p1.Asset = models.Asset("USD/2")

	err := store.UpsertPayments(context.Background(), models.ConnectorProviderDummyPay, []*models.Payment{p1})
	require.NoError(t, err)

	payment, err := store.GetPayment(context.Background(), p1ID.String())
	require.NoError(t, err)

	require.NotEqual(t, p1.CreatedAt, payment.CreatedAt)

	p1.CreatedAt = p1T
	testGetPayment(t, store, *p1ID, p1, nil)
}

func testPaymentsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	testGetPayment(t, store, *p1ID, nil, storage.ErrNotFound)
	testGetPayment(t, store, *p2ID, nil, storage.ErrNotFound)
}

func testListPayments(t *testing.T, store *storage.Storage) {
	query, err := storage.Paginate(1, "", nil, nil)
	require.NoError(t, err)

	payments, paginationDetails, err := store.ListPayments(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.True(t, paginationDetails.HasMore)
	checkPaymentsEqual(t, payments[0], p1)

	query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	payments, paginationDetails, err = store.ListPayments(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.False(t, paginationDetails.HasMore)
	checkPaymentsEqual(t, payments[0], p2)

	query, err = storage.Paginate(1, paginationDetails.PreviousPage, nil, nil)
	require.NoError(t, err)

	payments, paginationDetails, err = store.ListPayments(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, payments, 1)
	require.True(t, paginationDetails.HasMore)
	checkPaymentsEqual(t, payments[0], p1)

	query, err = storage.Paginate(2, "", nil, nil)
	require.NoError(t, err)

	payments, paginationDetails, err = store.ListPayments(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, payments, 2)
	require.False(t, paginationDetails.HasMore)
	checkPaymentsEqual(t, payments[0], p1)
	checkPaymentsEqual(t, payments[1], p2)
}
