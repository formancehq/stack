package api

import (
	"bytes"
	"encoding/json"
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
			name: "missing connectorId",
			req: &service.CreateBankAccountRequest{
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
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

			createBankAccountResponse := models.BankAccount{
				ID:            uuid.New(),
				CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				ConnectorID:   connectorID,
				Name:          "test_nominal",
				AccountNumber: "0112345678",
				IBAN:          "FR7630006000011234567890189",
				SwiftBicCode:  "HBUKGB4B",
				Country:       "FR",
			}

			expectedCreateBankAccountResponse := &bankAccountResponse{
				ID:          createBankAccountResponse.ID.String(),
				CreatedAt:   createBankAccountResponse.CreatedAt,
				Country:     createBankAccountResponse.Country,
				ConnectorID: createBankAccountResponse.ConnectorID.String(),
				AccountID:   createBankAccountResponse.AccountID.String(),
				Provider:    createBankAccountResponse.ConnectorID.Provider.String(),
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
