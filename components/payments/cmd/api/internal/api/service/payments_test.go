package service

import (
	"context"
	"errors"
	"testing"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

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

	service := New(&MockStore{})

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
