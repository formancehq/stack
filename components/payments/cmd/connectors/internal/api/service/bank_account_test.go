package service

import (
	"context"
	"testing"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateBankAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		req           *CreateBankAccountRequest
		errorPublish  bool
		expectedError error
	}

	connectorNotFound := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderCurrencyCloud,
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorBankingCircle.ID.String(),
				Name:          "test_nominal",
			},
		},
		{
			name: "invalid connector id",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   "invalid",
				Name:          "test_nominal",
			},
			expectedError: ErrValidation,
		},
		{
			name: "connector not found",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorNotFound.String(),
				Name:          "test_nominal",
			},
			expectedError: ErrValidation,
		},
		{
			name: "create bank account on another connector than BankingCircle",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal",
			},
			expectedError: ErrValidation,
		},
		{
			name: "publish error",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorBankingCircle.ID.String(),
				Name:          "test_nominal",
			},
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

			service := New(&MockStore{}, m.WithError(errPublish), nil)

			_, err := service.CreateBankAccount(context.Background(), tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
