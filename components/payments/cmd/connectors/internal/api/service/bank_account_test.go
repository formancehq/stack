package service

import (
	"context"
	"testing"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func TestCreateBankAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                          string
		req                           *CreateBankAccountRequest
		expectedError                 error
		noBankAccountCreateHandler    bool
		errorBankAccountCreateHandler error
	}

	connectorNotFound := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderCurrencyCloud,
	}

	var ErrOther = errors.New("other error")
	testCases := []testCase{
		{
			name: "nominal",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal",
			},
		},
		{
			name: "nominal with metadata",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal_metadata",
				Metadata:      map[string]string{"test": "metadata"},
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
			name: "no connector handler",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal",
			},
			noBankAccountCreateHandler: true,
			expectedError:              ErrValidation,
		},
		{
			name: "connector handler error validation",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal",
			},
			errorBankAccountCreateHandler: manager.ErrValidation,
			expectedError:                 ErrValidation,
		},
		{
			name: "connector handler error connector not found",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal",
			},
			errorBankAccountCreateHandler: manager.ErrConnectorNotFound,
			expectedError:                 ErrValidation,
		},
		{
			name: "connector handler other error",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorDummyPay.ID.String(),
				Name:          "test_nominal",
			},
			errorBankAccountCreateHandler: ErrOther,
			expectedError:                 ErrOther,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var handlers map[models.ConnectorProvider]*ConnectorHandlers
			if !tc.noBankAccountCreateHandler {
				handlers = map[models.ConnectorProvider]*ConnectorHandlers{
					models.ConnectorProviderDummyPay: {
						BankAccountHandler: func(ctx context.Context, bankAccount *models.BankAccount) error {
							if tc.errorBankAccountCreateHandler != nil {
								return tc.errorBankAccountCreateHandler
							}

							return nil
						},
					},
				}
			}

			service := New(&MockStore{}, &MockPublisher{}, handlers)

			_, err := service.CreateBankAccount(context.Background(), tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
