package storage_test

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	t1ID        models.TransferInitiationID
	t1T         = time.Date(2023, 11, 14, 5, 8, 0, 0, time.UTC)
	t1          *models.TransferInitiation
	adjumentID1 = uuid.New()

	t2ID        models.TransferInitiationID
	t2T         = time.Date(2023, 11, 14, 5, 7, 0, 0, time.UTC)
	t2          *models.TransferInitiation
	adjumentID2 = uuid.New()

	tAddPayments   = time.Date(2023, 11, 14, 5, 9, 10, 0, time.UTC)
	tUpdateStatus1 = time.Date(2023, 11, 14, 5, 9, 15, 0, time.UTC)
	tUpdateStatus2 = time.Date(2023, 11, 14, 5, 9, 16, 0, time.UTC)
)

func TestTransferInitiations(t *testing.T) {
	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreatePayments(t, store)
	testCreateTransferInitiations(t, store)
	testAddTransferInitiationPayments(t, store)
	testUpdateTransferInitiationStatus(t, store)
	testDeleteTransferInitiations(t, store)
	testUninstallConnectors(t, store)
	testTransferInitiationsDeletedAfterConnectorUninstall(t, store)
}

func testCreateTransferInitiations(t *testing.T, store *storage.Storage) {
	t1ID = models.TransferInitiationID{
		Reference:   "test1",
		ConnectorID: connectorID,
	}
	t1 = &models.TransferInitiation{
		ID:          t1ID,
		CreatedAt:   t1T,
		ScheduledAt: t1T,
		Description: "test_description",
		Type:        models.TransferInitiationTypeTransfer,
		ConnectorID: connectorID,
		Provider:    models.ConnectorProviderDummyPay,
		Amount:      big.NewInt(100),
		Asset:       models.Asset("USD/2"),
		RelatedAdjustments: []*models.TransferInitiationAdjustment{
			{
				ID:                   adjumentID1,
				TransferInitiationID: t1ID,
				CreatedAt:            t1T,
				Status:               models.TransferInitiationStatusWaitingForValidation,
			},
		},
	}

	t2ID = models.TransferInitiationID{
		Reference:   "test2",
		ConnectorID: connectorID,
	}
	t2 = &models.TransferInitiation{
		ID:                   t2ID,
		CreatedAt:            t2T,
		ScheduledAt:          t2T,
		Description:          "test_description2",
		Type:                 models.TransferInitiationTypeTransfer,
		ConnectorID:          connectorID,
		Provider:             models.ConnectorProviderDummyPay,
		Amount:               big.NewInt(150),
		Asset:                models.Asset("USD/2"),
		SourceAccountID:      &acc1ID,
		DestinationAccountID: acc2ID,
		RelatedAdjustments: []*models.TransferInitiationAdjustment{
			{
				ID:                   adjumentID2,
				TransferInitiationID: t2ID,
				CreatedAt:            t2T,
				Status:               models.TransferInitiationStatusWaitingForValidation,
			},
		},
	}

	// Missing source account id and destination account id
	err := store.CreateTransferInitiation(context.Background(), t1)
	require.Error(t, err)

	t1.SourceAccountID = &acc1ID
	t1.DestinationAccountID = acc2ID
	err = store.CreateTransferInitiation(context.Background(), t1)
	require.NoError(t, err)

	err = store.CreateTransferInitiation(context.Background(), t2)
	require.NoError(t, err)

	testGetTransferInitiation(t, store, t1ID, false, t1, nil)
	testGetTransferInitiation(t, store, t2ID, false, t2, nil)
}

func testGetTransferInitiation(
	t *testing.T,
	store *storage.Storage,
	id models.TransferInitiationID,
	expand bool,
	expected *models.TransferInitiation,
	expectedErr error,
) {
	tf, err := store.ReadTransferInitiation(context.Background(), id)
	if expectedErr != nil {
		require.EqualError(t, err, expectedErr.Error())
		return
	} else {
		require.NoError(t, err)
	}

	if expand {
		payments, err := store.ReadTransferInitiationPayments(context.Background(), id)
		require.NoError(t, err)
		tf.RelatedPayments = payments
	}

	checkTransferInitiationsEqual(t, expected, tf, true)
}

func checkTransferInitiationsEqual(t *testing.T, t1, t2 *models.TransferInitiation, checkRelatedAdjusment bool) {
	require.Equal(t, t1.ID, t2.ID)
	require.Equal(t, t1.CreatedAt.UTC(), t2.CreatedAt.UTC())
	require.Equal(t, t1.ScheduledAt.UTC(), t2.ScheduledAt.UTC())
	require.Equal(t, t1.Description, t2.Description)
	require.Equal(t, t1.Type, t2.Type)
	require.Equal(t, t1.Provider, t2.Provider)
	require.Equal(t, t1.Amount, t2.Amount)
	require.Equal(t, t1.Asset, t2.Asset)
	require.Equal(t, t1.SourceAccountID, t2.SourceAccountID)
	require.Equal(t, t1.DestinationAccountID, t2.DestinationAccountID)
	for i := range t1.RelatedPayments {
		require.Equal(t, t1.RelatedPayments[i].TransferInitiationID, t2.RelatedPayments[i].TransferInitiationID)
		require.Equal(t, t1.RelatedPayments[i].PaymentID, t2.RelatedPayments[i].PaymentID)
		require.Equal(t, t1.RelatedPayments[i].CreatedAt.UTC(), t2.RelatedPayments[i].CreatedAt.UTC())
		require.Equal(t, t1.RelatedPayments[i].Status, t2.RelatedPayments[i].Status)
		require.Equal(t, t1.RelatedPayments[i].Error, t2.RelatedPayments[i].Error)
	}
	if checkRelatedAdjusment {
		for i := range t1.RelatedAdjustments {
			require.Equal(t, t1.RelatedAdjustments[i].TransferInitiationID, t2.RelatedAdjustments[i].TransferInitiationID)
			require.Equal(t, t1.RelatedAdjustments[i].CreatedAt.UTC(), t2.RelatedAdjustments[i].CreatedAt.UTC())
			require.Equal(t, t1.RelatedAdjustments[i].Status, t2.RelatedAdjustments[i].Status)
			require.Equal(t, t1.RelatedAdjustments[i].Error, t2.RelatedAdjustments[i].Error)
			require.Equal(t, t1.RelatedAdjustments[i].Metadata, t2.RelatedAdjustments[i].Metadata)
		}
	}
}

func testAddTransferInitiationPayments(t *testing.T, store *storage.Storage) {
	err := store.AddTransferInitiationPaymentID(
		context.Background(),
		t1ID,
		p1ID,
		tAddPayments,
		map[string]string{
			"test": "test",
		},
	)
	require.NoError(t, err)

	t1.RelatedPayments = []*models.TransferInitiationPayment{
		{
			TransferInitiationID: t1ID,
			PaymentID:            *p1ID,
			CreatedAt:            tAddPayments,
			Status:               models.TransferInitiationStatusProcessing,
			Error:                "",
		},
	}
	t1.Metadata = map[string]string{
		"test": "test",
	}
	testGetTransferInitiation(t, store, t1ID, true, t1, nil)

	err = store.AddTransferInitiationPaymentID(
		context.Background(),
		t1ID,
		nil,
		tAddPayments,
		nil,
	)
	require.Error(t, err)

	err = store.AddTransferInitiationPaymentID(
		context.Background(),
		models.TransferInitiationID{
			Reference:   "not_existing",
			ConnectorID: connectorID,
		},
		p1ID,
		tAddPayments,
		nil,
	)
	require.Error(t, err)
}

func testUpdateTransferInitiationStatus(t *testing.T, store *storage.Storage) {
	adjustment1 := &models.TransferInitiationAdjustment{
		ID:                   uuid.New(),
		TransferInitiationID: t1ID,
		CreatedAt:            tUpdateStatus1,
		Status:               models.TransferInitiationStatusRejected,
		Error:                "test_error",
	}
	err := store.UpdateTransferInitiationPaymentsStatus(
		context.Background(),
		t1ID,
		nil,
		adjustment1,
	)
	require.NoError(t, err)

	t1.RelatedAdjustments = append([]*models.TransferInitiationAdjustment{
		adjustment1,
	}, t1.RelatedAdjustments...)
	testGetTransferInitiation(t, store, t1ID, true, t1, nil)

	adjustment2 := &models.TransferInitiationAdjustment{
		ID:                   uuid.New(),
		TransferInitiationID: t1ID,
		CreatedAt:            tUpdateStatus2,
		Status:               models.TransferInitiationStatusFailed,
		Error:                "test_error2",
	}
	err = store.UpdateTransferInitiationPaymentsStatus(
		context.Background(),
		t1ID,
		p1ID,
		adjustment2,
	)
	require.NoError(t, err)

	t1.RelatedPayments[0].Status = models.TransferInitiationStatusFailed
	t1.RelatedPayments[0].Error = "test_error2"
	t1.RelatedAdjustments = append([]*models.TransferInitiationAdjustment{
		adjustment2,
	}, t1.RelatedAdjustments...)
	testGetTransferInitiation(t, store, t1ID, true, t1, nil)
}

func testDeleteTransferInitiations(t *testing.T, store *storage.Storage) {
	err := store.DeleteTransferInitiation(context.Background(), t1ID)
	require.NoError(t, err)

	testGetTransferInitiation(t, store, t1ID, false, nil, storage.ErrNotFound)

	// Delete does not generate an error when not existing
	err = store.DeleteTransferInitiation(context.Background(), models.TransferInitiationID{
		Reference:   "not_existing",
		ConnectorID: connectorID,
	})
	require.NoError(t, err)
}

func testTransferInitiationsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	testGetTransferInitiation(t, store, t1ID, false, nil, storage.ErrNotFound)
	testGetTransferInitiation(t, store, t2ID, false, nil, storage.ErrNotFound)
}
