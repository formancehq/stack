package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/api/service"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/internal/models"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateBankAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.CreateBankAccountRequest
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	acc1 := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
		},
		{
			name: "nominal without connectorID",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				Name:          "test_nominal",
			},
		},
		{
			name:               "no body",
			req:                nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name: "missing AccountNumber and Iban",
			req: &service.CreateBankAccountRequest{
				SwiftBicCode: "HBUKGB4B",
				Country:      "FR",
				ConnectorID:  connectorID.String(),
				Name:         "test_nominal",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing name",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing country",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},

		{
			name: "service error duplicate key",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "service error not found",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "service error validation",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error invalid ID",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name: "service error publish",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			serviceError:       service.ErrPublish,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name: "service error other errors",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
				ConnectorID:   connectorID.String(),
				Name:          "test_nominal",
			},
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			bankAccountID := uuid.New()
			createBankAccountResponse := models.BankAccount{
				ID:            bankAccountID,
				CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Name:          "test_nominal",
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
			}

			if testCase.req != nil && testCase.req.ConnectorID != "" {
				createBankAccountResponse.RelatedAccounts = []*models.BankAccountRelatedAccount{
					{
						ID:            uuid.New(),
						BankAccountID: bankAccountID,
						ConnectorID:   connectorID,
						AccountID:     acc1,
					},
				}
			}

			expectedCreateBankAccountResponse := &bankAccountResponse{
				ID:        createBankAccountResponse.ID.String(),
				Name:      createBankAccountResponse.Name,
				CreatedAt: createBankAccountResponse.CreatedAt,
				Country:   createBankAccountResponse.Country,
			}

			if testCase.req != nil && testCase.req.ConnectorID != "" {
				expectedCreateBankAccountResponse.ConnectorID = createBankAccountResponse.RelatedAccounts[0].ConnectorID.String()
				expectedCreateBankAccountResponse.AccountID = createBankAccountResponse.RelatedAccounts[0].AccountID.String()
				expectedCreateBankAccountResponse.Provider = createBankAccountResponse.RelatedAccounts[0].ConnectorID.Provider.String()
				expectedCreateBankAccountResponse.RelatedAccounts = []*bankAccountRelatedAccountsResponse{
					{
						ID:          createBankAccountResponse.RelatedAccounts[0].ID.String(),
						AccountID:   createBankAccountResponse.RelatedAccounts[0].AccountID.String(),
						ConnectorID: createBankAccountResponse.RelatedAccounts[0].ConnectorID.String(),
						Provider:    createBankAccountResponse.RelatedAccounts[0].ConnectorID.Provider.String(),
					},
				}
			}

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					CreateBankAccount(gomock.Any(), testCase.req).
					Return(&createBankAccountResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					CreateBankAccount(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/bank-accounts", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[bankAccountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedCreateBankAccountResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestForwardBankAccountToConnector(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.ForwardBankAccountToConnectorRequest
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	acc1 := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
		},
		{
			name:               "nominal without connectorID",
			req:                &service.ForwardBankAccountToConnectorRequest{},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "no body",
			req:                nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name: "service error duplicate key",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "service error not found",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "service error validation",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error invalid ID",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name: "service error publish",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
			serviceError:       service.ErrPublish,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name: "service error other errors",
			req: &service.ForwardBankAccountToConnectorRequest{
				ConnectorID: connectorID.String(),
			},
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			bankAccountID := uuid.New()
			forwardBankAccountResponse := models.BankAccount{
				ID:            bankAccountID,
				CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Name:          "test_nominal",
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
			}

			if testCase.req != nil && testCase.req.ConnectorID != "" {
				forwardBankAccountResponse.RelatedAccounts = []*models.BankAccountRelatedAccount{
					{
						ID:            uuid.New(),
						BankAccountID: bankAccountID,
						ConnectorID:   connectorID,
						AccountID:     acc1,
					},
				}
			}

			expectedForwardBankAccountResponse := &bankAccountResponse{
				ID:        forwardBankAccountResponse.ID.String(),
				Name:      forwardBankAccountResponse.Name,
				CreatedAt: forwardBankAccountResponse.CreatedAt,
				Country:   forwardBankAccountResponse.Country,
			}

			if testCase.req != nil && testCase.req.ConnectorID != "" {
				expectedForwardBankAccountResponse.ConnectorID = forwardBankAccountResponse.RelatedAccounts[0].ConnectorID.String()
				expectedForwardBankAccountResponse.AccountID = forwardBankAccountResponse.RelatedAccounts[0].AccountID.String()
				expectedForwardBankAccountResponse.Provider = forwardBankAccountResponse.RelatedAccounts[0].ConnectorID.Provider.String()
				expectedForwardBankAccountResponse.RelatedAccounts = []*bankAccountRelatedAccountsResponse{
					{
						ID:          forwardBankAccountResponse.RelatedAccounts[0].ID.String(),
						AccountID:   forwardBankAccountResponse.RelatedAccounts[0].AccountID.String(),
						ConnectorID: forwardBankAccountResponse.RelatedAccounts[0].ConnectorID.String(),
						Provider:    forwardBankAccountResponse.RelatedAccounts[0].ConnectorID.Provider.String(),
					},
				}
			}

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ForwardBankAccountToConnector(gomock.Any(), bankAccountID.String(), testCase.req).
					Return(&forwardBankAccountResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ForwardBankAccountToConnector(gomock.Any(), bankAccountID.String(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/bank-accounts/%s/forward", bankAccountID), bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[bankAccountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedForwardBankAccountResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestUpdateBankAccountMetadata(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.UpdateBankAccountMetadataRequest
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name:               "empty metadata",
			req:                &service.UpdateBankAccountMetadataRequest{},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "no body",
			req:                nil,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name: "service error duplicate key",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "service error not found",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "service error validation",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error invalid ID",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name: "service error publish",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       service.ErrPublish,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name: "service error other errors",
			req: &service.UpdateBankAccountMetadataRequest{
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusNoContent
			}

			bankAccountID := uuid.New()

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					UpdateBankAccountMetadata(gomock.Any(), bankAccountID.String(), testCase.req).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					UpdateBankAccountMetadata(gomock.Any(), bankAccountID.String(), testCase.req).
					Return(testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPatch, fmt.Sprintf("/bank-accounts/%s/metadata", bankAccountID), bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode >= 300 || testCase.expectedStatusCode < 200 {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
