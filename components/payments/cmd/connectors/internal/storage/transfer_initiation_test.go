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
	t1ID = models.TransferInitiationID{
		Reference: "test1",
		Provider:  models.ConnectorProviderDummyPay,
	}
	t1T = time.Date(2023, 11, 14, 5, 8, 0, 0, time.UTC)
	t1  = &models.TransferInitiation{
		ID:          t1ID,
		CreatedAt:   t1T,
		ScheduledAt: t1T,
		UpdatedAt:   t1T,
		Description: "test_description",
		Type:        models.TransferInitiationTypeTransfer,
		Provider:    models.ConnectorProviderDummyPay,
		Amount:      big.NewInt(100),
		Asset:       models.Asset("USD/2"),
		Status:      models.TransferInitiationStatusWaitingForValidation,
	}

	t2ID = models.TransferInitiationID{
		Reference: "test2",
		Provider:  models.ConnectorProviderDummyPay,
	}
	t2T = time.Date(2023, 11, 14, 5, 7, 0, 0, time.UTC)
	t2  = &models.TransferInitiation{
		ID:                   t2ID,
		CreatedAt:            t2T,
		ScheduledAt:          t2T,
		UpdatedAt:            t2T,
		Description:          "test_description2",
		Type:                 models.TransferInitiationTypeTransfer,
		Provider:             models.ConnectorProviderDummyPay,
		Amount:               big.NewInt(150),
		Asset:                models.Asset("USD/2"),
		SourceAccountID:      acc1ID,
		DestinationAccountID: acc2ID,
		Status:               models.TransferInitiationStatusWaitingForValidation,
	}

	tAddPayments  = time.Date(2023, 11, 14, 5, 1, 10, 0, time.UTC)
	tUpdateStatus = time.Date(2023, 11, 14, 5, 1, 15, 0, time.UTC)
)

func TestTransferInitiations(t *testing.T) {
	t.Parallel()

	store := newStore(t)

	testInstallConnectors(t, store)
	testCreateAccounts(t, store)
	testCreatePayments(t, store)
	testCreateTransferInitiations(t, store)
	testListTransferInitiations(t, store)
	testAddTransferInitiationPayments(t, store)
	testUpdateTransferInitiationStatus(t, store)
	testDeleteTransferInitiations(t, store)
	testUninstallConnectors(t, store)
	testTransferInitiationsDeletedAfterConnectorUninstall(t, store)
}

func testCreateTransferInitiations(t *testing.T, store *storage.Storage) {
	// Missing source account id and destination account id
	err := store.CreateTransferInitiation(context.Background(), t1)
	require.Error(t, err)

	t1.SourceAccountID = acc1ID
	t1.DestinationAccountID = acc2ID
	err = store.CreateTransferInitiation(context.Background(), t1)
	require.NoError(t, err)

	err = store.CreateTransferInitiation(context.Background(), t2)
	require.NoError(t, err)

	testGetTransferInitiation(t, store, t1ID, false, t1, nil)
	testGetTransferInitiation(t, store, t2ID, false, t2, nil)
}

func testListTransferInitiations(t *testing.T, store *storage.Storage) {
	query, err := storage.Paginate(1, "", nil, nil)
	require.NoError(t, err)

	tfs, paginationDetails, err := store.ListTransferInitiations(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, tfs, 1)
	require.True(t, paginationDetails.HasMore)
	checkTransferInitiationsEqual(t, t1, tfs[0])

	query, err = storage.Paginate(1, paginationDetails.NextPage, nil, nil)
	require.NoError(t, err)

	tfs, paginationDetails, err = store.ListTransferInitiations(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, tfs, 1)
	require.False(t, paginationDetails.HasMore)
	checkTransferInitiationsEqual(t, t2, tfs[0])

	query, err = storage.Paginate(1, paginationDetails.PreviousPage, nil, nil)
	require.NoError(t, err)

	tfs, paginationDetails, err = store.ListTransferInitiations(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, tfs, 1)
	require.True(t, paginationDetails.HasMore)
	checkTransferInitiationsEqual(t, t1, tfs[0])

	query, err = storage.Paginate(2, "", nil, nil)
	require.NoError(t, err)

	tfs, paginationDetails, err = store.ListTransferInitiations(context.Background(), query)
	require.NoError(t, err)
	require.Len(t, tfs, 2)
	require.False(t, paginationDetails.HasMore)
	checkTransferInitiationsEqual(t, t1, tfs[0])
	checkTransferInitiationsEqual(t, t2, tfs[1])
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

	checkTransferInitiationsEqual(t, expected, tf)
}

func checkTransferInitiationsEqual(t *testing.T, t1, t2 *models.TransferInitiation) {
	require.Equal(t, t1.ID, t2.ID)
	require.Equal(t, t1.CreatedAt.UTC(), t2.CreatedAt.UTC())
	require.Equal(t, t1.ScheduledAt.UTC(), t2.ScheduledAt.UTC())
	require.Equal(t, t1.UpdatedAt.UTC(), t2.UpdatedAt.UTC())
	require.Equal(t, t1.Description, t2.Description)
	require.Equal(t, t1.Type, t2.Type)
	require.Equal(t, t1.Provider, t2.Provider)
	require.Equal(t, t1.Amount, t2.Amount)
	require.Equal(t, t1.Asset, t2.Asset)
	require.Equal(t, t1.SourceAccountID, t2.SourceAccountID)
	require.Equal(t, t1.DestinationAccountID, t2.DestinationAccountID)
	require.Equal(t, t1.Status, t2.Status)
	for i := range t1.RelatedPayments {
		require.Equal(t, t1.RelatedPayments[i].TransferInitiationID, t2.RelatedPayments[i].TransferInitiationID)
		require.Equal(t, t1.RelatedPayments[i].PaymentID, t2.RelatedPayments[i].PaymentID)
		require.Equal(t, t1.RelatedPayments[i].CreatedAt.UTC(), t2.RelatedPayments[i].CreatedAt.UTC())
		require.Equal(t, t1.RelatedPayments[i].Status, t2.RelatedPayments[i].Status)
		require.Equal(t, t1.RelatedPayments[i].Error, t2.RelatedPayments[i].Error)
	}
}

func testAddTransferInitiationPayments(t *testing.T, store *storage.Storage) {
	err := store.AddTransferInitiationPaymentID(
		context.Background(),
		t1ID,
		p1ID,
		tAddPayments,
	)
	require.NoError(t, err)

	t1.RelatedPayments = []*models.TransferInitiationPayments{
		{
			TransferInitiationID: t1ID,
			PaymentID:            *p1ID,
			CreatedAt:            tAddPayments,
			Status:               models.TransferInitiationStatusProcessing,
			Error:                "",
		},
	}
	testGetTransferInitiation(t, store, t1ID, true, t1, nil)

	err = store.AddTransferInitiationPaymentID(
		context.Background(),
		t1ID,
		nil,
		tAddPayments,
	)
	require.Error(t, err)

	err = store.AddTransferInitiationPaymentID(
		context.Background(),
		models.TransferInitiationID{
			Reference: "not_existing",
			Provider:  models.ConnectorProviderDummyPay,
		},
		p1ID,
		tAddPayments,
	)
	require.Error(t, err)
}

func testUpdateTransferInitiationStatus(t *testing.T, store *storage.Storage) {
	err := store.UpdateTransferInitiationPaymentsStatus(
		context.Background(),
		t1ID,
		nil,
		models.TransferInitiationStatusRejected,
		"test_error",
		2,
		tUpdateStatus,
	)
	require.NoError(t, err)

	t1.Status = models.TransferInitiationStatusRejected
	t1.Error = "test_error"
	t1.UpdatedAt = tUpdateStatus
	t1.Attempts = 2
	testGetTransferInitiation(t, store, t1ID, true, t1, nil)

	err = store.UpdateTransferInitiationPaymentsStatus(
		context.Background(),
		t1ID,
		p1ID,
		models.TransferInitiationStatusFailed,
		"test_error2",
		3,
		tUpdateStatus,
	)
	require.NoError(t, err)

	t1.RelatedPayments[0].Status = models.TransferInitiationStatusFailed
	t1.RelatedPayments[0].Error = "test_error2"
	t1.Status = models.TransferInitiationStatusFailed
	t1.Error = "test_error2"
	t1.UpdatedAt = tUpdateStatus
	t1.Attempts = 3
	testGetTransferInitiation(t, store, t1ID, true, t1, nil)
}

func testDeleteTransferInitiations(t *testing.T, store *storage.Storage) {
	err := store.DeleteTransferInitiation(context.Background(), t1ID)
	require.NoError(t, err)

	testGetTransferInitiation(t, store, t1ID, false, nil, storage.ErrNotFound)

	// Delete does not generate an error when not existing
	err = store.DeleteTransferInitiation(context.Background(), models.TransferInitiationID{
		Reference: "not_existing",
		Provider:  models.ConnectorProviderDummyPay,
	})
	require.NoError(t, err)
}

func testTransferInitiationsDeletedAfterConnectorUninstall(t *testing.T, store *storage.Storage) {
	testGetTransferInitiation(t, store, t1ID, false, nil, storage.ErrNotFound)
	testGetTransferInitiation(t, store, t2ID, false, nil, storage.ErrNotFound)
}
