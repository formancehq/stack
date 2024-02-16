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
	p1ID *models.PaymentID
	p1T  = time.Date(2023, 11, 14, 4, 55, 0, 0, time.UTC)
	p1   *models.Payment

	p2ID *models.PaymentID
	p2T  = time.Date(2023, 11, 14, 4, 54, 0, 0, time.UTC)
	p2   *models.Payment
)

func TestPayments(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreatePayments(t, store)
	testUpdatePayment(t, store)
	testUninstallConnectors(t, store)
	testPaymentsDeletedAfterConnectorUninstall(t, store)
}

func testCreatePayments(t *testing.T, store *storage.Storage) {
	p1ID = &models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "test1",
			Type:      models.PaymentTypePayOut,
		},
		ConnectorID: connectorID,
	}
	p1 = &models.Payment{
		ID:          *p1ID,
		CreatedAt:   p1T,
		Reference:   "ref1",
		Amount:      big.NewInt(100),
		ConnectorID: connectorID,
		Type:        models.PaymentTypePayOut,
		Status:      models.PaymentStatusSucceeded,
		Scheme:      models.PaymentSchemeCardVisa,
		Asset:       models.Asset("USD/2"),
	}

	p2ID = &models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "test2",
			Type:      models.PaymentTypeTransfer,
		},
		ConnectorID: connectorID,
	}
	p2 = &models.Payment{
		ID:          *p2ID,
		CreatedAt:   p2T,
		Reference:   "ref2",
		Amount:      big.NewInt(150),
		ConnectorID: connectorID,
		Type:        models.PaymentTypePayIn,
		Status:      models.PaymentStatusFailed,
		Scheme:      models.PaymentSchemeCardVisa,
		Asset:       models.Asset("EUR/2"),
	}

	pFail := &models.Payment{
		ID:          *p1ID,
		CreatedAt:   p1T,
		Reference:   "ref1",
		ConnectorID: connectorID,
		Amount:      big.NewInt(100),
		Type:        models.PaymentTypePayOut,
		Status:      models.PaymentStatusSucceeded,
		Scheme:      models.PaymentSchemeCardVisa,
		Asset:       models.Asset("USD/2"),
		SourceAccountID: &models.AccountID{
			Reference:   "not_existing",
			ConnectorID: connectorID,
		},
	}

	ids, err := store.UpsertPayments(context.Background(), []*models.Payment{pFail})
	require.Error(t, err)
	require.Len(t, ids, 0)

	ids, err = store.UpsertPayments(context.Background(), []*models.Payment{p1})
	require.NoError(t, err)
	require.Len(t, ids, 1)

	ids, err = store.UpsertPayments(context.Background(), []*models.Payment{p2})
	require.NoError(t, err)
	require.Len(t, ids, 1)

	p1.Status = models.PaymentStatusPending
	p2.Status = models.PaymentStatusSucceeded
	ids, err = store.UpsertPayments(context.Background(), []*models.Payment{p1, p2})
	require.NoError(t, err)
	require.Len(t, ids, 2)

	ids, err = store.UpsertPayments(context.Background(), []*models.Payment{p1, p2})
	require.NoError(t, err)
	require.Len(t, ids, 0)

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
	p1.CreatedAt = time.Date(2023, 11, 14, 5, 55, 0, 0, time.UTC)
	p1.Reference = "ref1_updated"
	p1.Amount = big.NewInt(150)
	p1.Type = models.PaymentTypePayIn
	p1.Status = models.PaymentStatusPending
	p1.Scheme = models.PaymentSchemeCardVisa
	p1.Asset = models.Asset("USD/2")

	ids, err := store.UpsertPayments(context.Background(), []*models.Payment{p1})
	require.NoError(t, err)
	require.Len(t, ids, 1)

	payment, err := store.GetPayment(context.Background(), p1ID.String())
	require.NoError(t, err)

	require.NotEqual(t, p1.Reference, payment.Reference)

	p1.Reference = payment.Reference
	testGetPayment(t, store, *p1ID, p1, nil)
}

func testPaymentsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	testGetPayment(t, store, *p1ID, nil, storage.ErrNotFound)
	testGetPayment(t, store, *p2ID, nil, storage.ErrNotFound)
}
