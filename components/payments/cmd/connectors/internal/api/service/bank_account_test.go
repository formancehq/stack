package service

import (
	"context"
	"testing"
	"time"

	manager "github.com/formancehq/payments/cmd/connectors/internal/api/connectors_manager"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/messages"
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
			name: "nominal without connectorID",
			req: &CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
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
						BankAccountHandler: func(ctx context.Context, connectorID models.ConnectorID, bankAccount *models.BankAccount) error {
							if tc.errorBankAccountCreateHandler != nil {
								return tc.errorBankAccountCreateHandler
							}

							return nil
						},
					},
				}
			}

			service := New(&MockStore{}, &MockPublisher{}, messages.NewMessages(""), handlers)

			_, err := service.CreateBankAccount(context.Background(), tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestForwardBankAccountToConnector(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                           string
		bankAccountID                  string
		req                            *ForwardBankAccountToConnectorRequest
		withBankAccountRelatedAccounts []*models.BankAccountRelatedAccount
		expectedError                  error
		noBankAccountForwardHandler    bool
		errorBankAccountForwardHandler error
	}

	connectorNotFound := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderCurrencyCloud,
	}

	var ErrOther = errors.New("other error")
	bankAccountID := uuid.New()
	testCases := []testCase{
		{
			name:          "nominal",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
		},
		{
			name:          "already forwarded to connector",
			bankAccountID: bankAccountID.String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			withBankAccountRelatedAccounts: []*models.BankAccountRelatedAccount{
				{
					ID:            uuid.New(),
					CreatedAt:     time.Now().UTC(),
					BankAccountID: bankAccountID,
					ConnectorID:   connectorDummyPay.ID,
					AccountID: models.AccountID{
						Reference:   "test",
						ConnectorID: connectorDummyPay.ID,
					},
				},
			},
			expectedError: ErrValidation,
		},
		{
			name: "empty bank account id",
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			expectedError: ErrInvalidID,
		},
		{
			name:          "invalid bank account id",
			bankAccountID: "invalid",
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			expectedError: ErrInvalidID,
		},
		{
			name:          "missing connectorID",
			bankAccountID: uuid.New().String(),
			req:           &ForwardBankAccountToConnectorRequest{},
			expectedError: ErrValidation,
		},
		{
			name:          "invalid connector id",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: "invalid",
			},
			expectedError: ErrValidation,
		},
		{
			name:          "connector not found",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorNotFound.String(),
			},
			expectedError: ErrValidation,
		},
		{
			name:          "no connector handler",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			noBankAccountForwardHandler: true,
			expectedError:               ErrValidation,
		},
		{
			name:          "connector handler error validation",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			errorBankAccountForwardHandler: manager.ErrValidation,
			expectedError:                  ErrValidation,
		},
		{
			name:          "connector handler error connector not found",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			errorBankAccountForwardHandler: manager.ErrConnectorNotFound,
			expectedError:                  ErrValidation,
		},
		{
			name:          "connector handler other error",
			bankAccountID: uuid.New().String(),
			req: &ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorDummyPay.ID.String(),
			},
			errorBankAccountForwardHandler: ErrOther,
			expectedError:                  ErrOther,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var handlers map[models.ConnectorProvider]*ConnectorHandlers
			if !tc.noBankAccountForwardHandler {
				handlers = map[models.ConnectorProvider]*ConnectorHandlers{
					models.ConnectorProviderDummyPay: {
						BankAccountHandler: func(ctx context.Context, connectorID models.ConnectorID, bankAccount *models.BankAccount) error {
							if tc.errorBankAccountForwardHandler != nil {
								return tc.errorBankAccountForwardHandler
							}

							return nil
						},
					},
				}
			}

			store := &MockStore{}
			service := New(store.WithBankAccountRelatedAccounts(tc.withBankAccountRelatedAccounts), &MockPublisher{}, messages.NewMessages(""), handlers)

			_, err := service.ForwardBankAccountToConnector(context.Background(), tc.bankAccountID, tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestUpdateBankAccountMetadata(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		bankAccountID string
		req           *UpdateBankAccountMetadataRequest
		storageError  error
		expectedError error
	}

	testCases := []testCase{
		{
			name:          "nominal",
			bankAccountID: uuid.New().String(),
			req: &UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name:          "err not found from storage",
			bankAccountID: uuid.New().String(),
			req: &UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			storageError:  storage.ErrNotFound,
			expectedError: storage.ErrNotFound,
		},
		{
			name: "empty bank account id",
			req: &UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedError: ErrInvalidID,
		},
		{
			name:          "invalid bank account id",
			bankAccountID: "invalid",
			req: &UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedError: ErrInvalidID,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var handlers map[models.ConnectorProvider]*ConnectorHandlers

			store := &MockStore{}
			if tc.storageError != nil {
				store = store.WithError(tc.storageError)
			}
			service := New(store, &MockPublisher{}, messages.NewMessages(""), handlers)

			err := service.UpdateBankAccountMetadata(context.Background(), tc.bankAccountID, tc.req)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
