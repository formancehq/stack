package service

import (
	"context"
	"math/big"
	"testing"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateTransferInitiation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                string
		req                 *CreateTransferInitiationRequest
		expectedTF          *models.TransferInitiation
		listConnectorLength int
		errorPublish        bool
		errorPaymentHandler error
		noPaymentsHandler   bool
		expectedError       error
	}

	sourceAccountID := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorDummyPay.ID,
	}

	destinationAccountID := models.AccountID{
		Reference:   "acc2",
		ConnectorID: connectorDummyPay.ID,
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedTF: &models.TransferInitiation{
				ID: models.TransferInitiationID{
					Reference:   "ref1",
					ConnectorID: connectorDummyPay.ID,
				},
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				DestinationAccountID: destinationAccountID,
				SourceAccountID:      &sourceAccountID,
				ConnectorID:          connectorDummyPay.ID,
				Provider:             models.ConnectorProviderDummyPay,
				Type:                 models.TransferInitiationTypeTransfer,
				Amount:               big.NewInt(100),
				InitialAmount:        big.NewInt(100),
				Asset:                models.Asset("EUR/2"),
				RelatedAdjustments: []*models.TransferInitiationAdjustment{
					{
						TransferInitiationID: models.TransferInitiationID{
							Reference:   "ref1",
							ConnectorID: connectorDummyPay.ID,
						},
						Status: models.TransferInitiationStatusWaitingForValidation,
					},
				},
			},
		},
		{
			name: "nominal without description",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedTF: &models.TransferInitiation{
				ID: models.TransferInitiationID{
					Reference:   "ref1",
					ConnectorID: connectorDummyPay.ID,
				},
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				DestinationAccountID: destinationAccountID,
				SourceAccountID:      &sourceAccountID,
				ConnectorID:          connectorDummyPay.ID,
				Provider:             models.ConnectorProviderDummyPay,
				Type:                 models.TransferInitiationTypeTransfer,
				Amount:               big.NewInt(100),
				InitialAmount:        big.NewInt(100),
				Asset:                models.Asset("EUR/2"),
				RelatedAdjustments: []*models.TransferInitiationAdjustment{
					{
						TransferInitiationID: models.TransferInitiationID{
							Reference:   "ref1",
							ConnectorID: connectorDummyPay.ID,
						},
						Status: models.TransferInitiationStatusWaitingForValidation,
					},
				},
			},
		},
		{
			name: "nominal with status changed",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            true,
			},
			expectedTF: &models.TransferInitiation{
				ID: models.TransferInitiationID{
					Reference:   "ref1",
					ConnectorID: connectorDummyPay.ID,
				},
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				DestinationAccountID: destinationAccountID,
				SourceAccountID:      &sourceAccountID,
				ConnectorID:          connectorDummyPay.ID,
				Provider:             models.ConnectorProviderDummyPay,
				Type:                 models.TransferInitiationTypeTransfer,
				Amount:               big.NewInt(100),
				InitialAmount:        big.NewInt(100),
				Asset:                models.Asset("EUR/2"),
				RelatedAdjustments: []*models.TransferInitiationAdjustment{
					{
						TransferInitiationID: models.TransferInitiationID{
							Reference:   "ref1",
							ConnectorID: connectorDummyPay.ID,
						},
						Status: models.TransferInitiationStatusValidated,
					},
				},
			},
		},
		{
			name: "transfer with external account as destination",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationExternalAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedError: ErrValidation,
		},
		{
			name: "payout with internal account as destination",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypePayout.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedError: ErrValidation,
		},
		{
			name: "invalid connector id",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          "invalid",
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedError: ErrValidation,
		},
		{
			name: "invalid provider",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				Provider:             "invalid",
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedError: ErrValidation,
		},
		{
			name: "too many connectors list",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			listConnectorLength: 2,
			expectedError:       ErrValidation,
		},
		{
			name: "no connectors in connectors list",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			listConnectorLength: 0,
			expectedError:       ErrValidation,
		},
		{
			name: "connector not installed",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorBankingCircle.ID.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedError: ErrValidation,
		},
		{
			name: "error publishing",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			errorPublish:  true,
			expectedError: ErrPublish,
		},
		{
			name: "no payments handler found",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            true,
			},
			noPaymentsHandler: true,
			expectedError:     ErrValidation,
		},
		{
			name: "error in payments handler",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            true,
			},
			errorPaymentHandler: manager.ErrValidation,
			expectedError:       ErrValidation,
		},
		{
			name: "error in payments handler",
			req: &CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Description:          "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorDummyPay.ID.String(),
				Provider:             string(models.ConnectorProviderDummyPay),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            true,
			},
			errorPaymentHandler: manager.ErrConnectorNotFound,
			expectedError:       ErrValidation,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockPublisher{}
			s := &MockStore{}

			var errPublish error
			if tc.errorPublish {
				errPublish = errors.New("publish error")
			}

			var handlers map[models.ConnectorProvider]*ConnectorHandlers
			if !tc.noPaymentsHandler {
				handlers = map[models.ConnectorProvider]*ConnectorHandlers{
					models.ConnectorProviderDummyPay: {
						InitiatePaymentHandler: func(ctx context.Context, transfer *models.TransferInitiation) error {
							if tc.errorPaymentHandler != nil {
								return tc.errorPaymentHandler
							}

							return nil
						},
					},
				}
			}
			service := New(s.WithListConnectorsNB(tc.listConnectorLength), m.WithError(errPublish), messages.NewMessages(""), handlers)

			tf, err := service.CreateTransferInitiation(context.Background(), tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
				tc.expectedTF.CreatedAt = tf.CreatedAt
				require.Len(t, tf.RelatedAdjustments, 1)
				tc.expectedTF.RelatedAdjustments[0].CreatedAt = tf.RelatedAdjustments[0].CreatedAt
				tc.expectedTF.RelatedAdjustments[0].ID = tf.RelatedAdjustments[0].ID
				require.Equal(t, tc.expectedTF, tf)
			}
		})
	}
}

func TestUpdateTransferInitiationStatus(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                string
		transferID          string
		req                 *UpdateTransferInitiationStatusRequest
		errorPublish        bool
		errorPaymentHandler error
		noPaymentsHandler   bool
		expectedError       error
	}

	tfNotFoundID := models.TransferInitiationID{
		Reference:   "not_found",
		ConnectorID: connectorDummyPay.ID,
	}

	testCases := []testCase{
		{
			name: "nominal validated",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID: transferInitiationWaiting.ID.String(),
		},
		{
			name: "nominal rejected",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "REJECTED",
			},
			transferID: transferInitiationWaiting.ID.String(),
		},
		{
			name: "unknown status",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "INVALID",
			},
			transferID:    transferInitiationWaiting.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name: "invalid status waiting for validation",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "WAITING_FOR_VALIDATION",
			},
			transferID:    transferInitiationWaiting.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name: "invalid status failed",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "FAILED",
			},
			transferID:    transferInitiationWaiting.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name: "invalid status processed",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "PROCESSED",
			},
			transferID:    transferInitiationWaiting.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name: "invalid status processing",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "PROCESSING",
			},
			transferID:    transferInitiationWaiting.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name: "invalid transfer id",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:    "invalid",
			expectedError: ErrInvalidID,
		},
		{
			name: "transfer not found",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:    tfNotFoundID.String(),
			expectedError: storage.ErrNotFound,
		},
		{
			name: "previous transfer with wrong status",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:    transferInitiationFailed.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name: "error publishing",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:    transferInitiationWaiting.ID.String(),
			errorPublish:  true,
			expectedError: ErrPublish,
		},
		{
			name: "error publishing",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:        transferInitiationWaiting.ID.String(),
			noPaymentsHandler: true,
			expectedError:     ErrValidation,
		},
		{
			name: "error in payments handler",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:          transferInitiationWaiting.ID.String(),
			errorPaymentHandler: manager.ErrValidation,
			expectedError:       ErrValidation,
		},
		{
			name: "error in payments handler",
			req: &UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID:          transferInitiationWaiting.ID.String(),
			errorPaymentHandler: manager.ErrConnectorNotFound,
			expectedError:       ErrValidation,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockPublisher{}

			var errPublish error
			if tc.errorPublish {
				errPublish = errors.New("publish error")
			}

			var handlers map[models.ConnectorProvider]*ConnectorHandlers
			if !tc.noPaymentsHandler {
				handlers = map[models.ConnectorProvider]*ConnectorHandlers{
					models.ConnectorProviderDummyPay: {
						InitiatePaymentHandler: func(ctx context.Context, transfer *models.TransferInitiation) error {
							if tc.errorPaymentHandler != nil {
								return tc.errorPaymentHandler
							}

							return nil
						},
					},
				}
			}
			service := New(&MockStore{}, m.WithError(errPublish), messages.NewMessages(""), handlers)

			err := service.UpdateTransferInitiationStatus(context.Background(), tc.transferID, tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRetryTransferInitiation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                string
		transferID          string
		errorPublish        bool
		errorPaymentHandler error
		noPaymentsHandler   bool
		expectedError       error
	}

	testCases := []testCase{
		{
			name:       "nominal",
			transferID: transferInitiationFailed.ID.String(),
		},
		{
			name:          "invalid transfer id",
			transferID:    "invalid",
			expectedError: ErrInvalidID,
		},
		{
			name:          "invalid previous transfer status",
			transferID:    transferInitiationWaiting.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name:              "error publishing",
			transferID:        transferInitiationFailed.ID.String(),
			noPaymentsHandler: true,
			expectedError:     ErrValidation,
		},
		{
			name:                "error in payments handler",
			transferID:          transferInitiationFailed.ID.String(),
			errorPaymentHandler: manager.ErrValidation,
			expectedError:       ErrValidation,
		},
		{
			name:                "error in payments handler",
			transferID:          transferInitiationFailed.ID.String(),
			errorPaymentHandler: manager.ErrConnectorNotFound,
			expectedError:       ErrValidation,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockPublisher{}

			var errPublish error
			if tc.errorPublish {
				errPublish = errors.New("publish error")
			}

			var handlers map[models.ConnectorProvider]*ConnectorHandlers
			if !tc.noPaymentsHandler {
				handlers = map[models.ConnectorProvider]*ConnectorHandlers{
					models.ConnectorProviderDummyPay: {
						InitiatePaymentHandler: func(ctx context.Context, transfer *models.TransferInitiation) error {
							if tc.errorPaymentHandler != nil {
								return tc.errorPaymentHandler
							}

							return nil
						},
					},
				}
			}
			service := New(&MockStore{}, m.WithError(errPublish), messages.NewMessages(""), handlers)

			err := service.RetryTransferInitiation(context.Background(), tc.transferID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeleteTransferInitiation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		transferID    string
		errorPublish  bool
		expectedError error
	}

	testCases := []testCase{
		{
			name:       "nominal",
			transferID: transferInitiationWaiting.ID.String(),
		},
		{
			name:          "invalid transfer id",
			transferID:    "invalid",
			expectedError: ErrInvalidID,
		},
		{
			name:          "invalid previous transfer initiation status",
			transferID:    transferInitiationFailed.ID.String(),
			expectedError: ErrValidation,
		},
		{
			name:          "error publishing",
			transferID:    transferInitiationWaiting.ID.String(),
			errorPublish:  true,
			expectedError: ErrPublish,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			m := &MockPublisher{}

			var errPublish error
			if tc.errorPublish {
				errPublish = errors.New("publish error")
			}

			service := New(&MockStore{}, m.WithError(errPublish), messages.NewMessages(""), nil)

			err := service.DeleteTransferInitiation(context.Background(), tc.transferID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
