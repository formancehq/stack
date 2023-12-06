package service

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreatePayment(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                 string
		request              *CreatePaymentRequest
		isConnectorInstalled bool
		expectedError        error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	sourceAccountID := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	destinationAccountID := models.AccountID{
		Reference:   "acc2",
		ConnectorID: connectorID,
	}

	testCases := []testCase{
		{
			name: "nominal",
			request: &CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			isConnectorInstalled: true,
		},
		{
			name: "connector not installed",
			request: &CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			isConnectorInstalled: false,
			expectedError:        ErrValidation,
		},
		{
			name: "nominal without source or destination account ids",
			request: &CreatePaymentRequest{
				Reference:   "test",
				ConnectorID: connectorID.String(),
				CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:      big.NewInt(100),
				Type:        string(models.PaymentTypeTransfer),
				Status:      string(models.PaymentStatusSucceeded),
				Scheme:      string(models.PaymentSchemeOther),
				Asset:       "EUR/2",
			},
			isConnectorInstalled: true,
		},
		{
			name: "invalid connectorID",
			request: &CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          "invalid",
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedError:        ErrValidation,
			isConnectorInstalled: true,
		},
		{
			name: "invalid source account id",
			request: &CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      "invalid",
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedError:        ErrValidation,
			isConnectorInstalled: true,
		},
		{
			name: "invalid destination account id",
			request: &CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: "invalid",
			},
			expectedError:        ErrValidation,
			isConnectorInstalled: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &MockStore{}
			service := New(store.WithIsConnectorInstalled(tc.isConnectorInstalled), &MockPublisher{})
			p, err := service.CreatePayment(context.Background(), tc.request)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
				require.NotNil(t, p)
			}
		})
	}
}

func TestGetPayment(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		paymentID     string
		expectedError error
	}

	paymentID := models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "p1",
			Type:      models.PaymentTypePayIn,
		},
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:          "nominal",
			paymentID:     paymentID.String(),
			expectedError: nil,
		},
		{
			name:          "invalid paymentID",
			paymentID:     "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{}, &MockPublisher{})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.GetPayment(context.Background(), tc.paymentID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
