package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
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

func TestCreateTransferInitiations(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.CreateTransferInitiationRequest
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
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
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "nominal without description",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "missing reference",
			req: &service.CreateTransferInitiationRequest{
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing destination account id",
			req: &service.CreateTransferInitiationRequest{
				Reference:       "ref1",
				ScheduledAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:     "test_nominal",
				SourceAccountID: sourceAccountID.String(),
				ConnectorID:     connectorID.String(),
				Provider:        models.ConnectorProviderDummyPay.String(),
				Type:            models.TransferInitiationTypeTransfer.String(),
				Amount:          big.NewInt(100),
				Asset:           "EUR/2",
				Validated:       false,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing source account id, should not end in error",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
		},
		{
			name: "wrong transfer initiation type",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 "invalid",
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing amount",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Asset:                "EUR/2",
				Validated:            false,
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing asset",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Validated:            false,
			},
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
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "service error not found",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "service error validation",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error invalid ID",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name: "service error publish",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
			},
			serviceError:       service.ErrPublish,
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name: "service error other errors",
			req: &service.CreateTransferInitiationRequest{
				Reference:            "ref1",
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				ConnectorID:          connectorID.String(),
				Provider:             models.ConnectorProviderDummyPay.String(),
				Type:                 models.TransferInitiationTypeTransfer.String(),
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Validated:            false,
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

			createTransferInitiationResponse := models.TransferInitiation{
				ID: models.TransferInitiationID{
					Reference:   "ref1",
					ConnectorID: connectorID,
				},
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          "test_nominal",
				Type:                 models.TransferInitiationTypeTransfer,
				SourceAccountID:      &sourceAccountID,
				DestinationAccountID: destinationAccountID,
				Provider:             models.ConnectorProviderDummyPay,
				ConnectorID:          connectorID,
				Amount:               big.NewInt(100),
				Asset:                "EUR/2",
				Metadata: map[string]string{
					"foo": "bar",
				},
				RelatedAdjustments: []*models.TransferInitiationAdjustment{
					{
						ID: uuid.New(),
						TransferInitiationID: models.TransferInitiationID{
							Reference:   "ref1",
							ConnectorID: connectorID,
						},
						CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
						Status:    models.TransferInitiationStatusProcessing,
					},
				},
			}

			expectedCreateTransferInitiationResponse := &transferInitiationResponse{
				ID: models.TransferInitiationID{
					Reference:   "ref1",
					ConnectorID: connectorID,
				}.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				ScheduledAt:          time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Description:          createTransferInitiationResponse.Description,
				SourceAccountID:      createTransferInitiationResponse.SourceAccountID.String(),
				DestinationAccountID: createTransferInitiationResponse.DestinationAccountID.String(),
				ConnectorID:          createTransferInitiationResponse.ConnectorID.String(),
				Type:                 createTransferInitiationResponse.Type.String(),
				Amount:               createTransferInitiationResponse.Amount,
				Asset:                createTransferInitiationResponse.Asset.String(),
				Status:               models.TransferInitiationStatusProcessing.String(),
				Metadata:             createTransferInitiationResponse.Metadata,
			}

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					CreateTransferInitiation(gomock.Any(), testCase.req).
					Return(&createTransferInitiationResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					CreateTransferInitiation(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/transfer-initiations", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[transferInitiationResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedCreateTransferInitiationResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestUpdateTransferInitiationStatus(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.UpdateTransferInitiationStatusRequest
		transferID         string
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	transferID := models.TransferInitiationID{
		Reference: "ref1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			transferID: transferID.String(),
		},
		{
			name:               "missing body",
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name: "service error duplicate key",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "service error not found",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			serviceError:       storage.ErrNotFound,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "service error validation",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			serviceError:       service.ErrValidation,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error invalid ID",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			serviceError:       service.ErrInvalidID,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name: "service error publish",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			serviceError:       service.ErrPublish,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name: "service error other errors",
			req: &service.UpdateTransferInitiationStatusRequest{
				Status: "VALIDATED",
			},
			serviceError:       errors.New("some error"),
			transferID:         transferID.String(),
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

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					UpdateTransferInitiationStatus(gomock.Any(), testCase.transferID, testCase.req).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					UpdateTransferInitiationStatus(gomock.Any(), testCase.transferID, testCase.req).
					Return(testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transfer-initiations/%s/status", testCase.transferID), bytes.NewReader(body))
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

func TestRetryTransferInitiation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		transferID         string
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	transferID := models.TransferInitiationID{
		Reference: "ref1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:       "nominal",
			transferID: transferID.String(),
		},
		{
			name:               "service error duplicate key",
			serviceError:       storage.ErrDuplicateKeyValue,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "service error not found",
			serviceError:       storage.ErrNotFound,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "service error validation",
			serviceError:       service.ErrValidation,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "service error invalid ID",
			serviceError:       service.ErrInvalidID,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "service error publish",
			serviceError:       service.ErrPublish,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name:               "service error other errors",
			serviceError:       errors.New("some error"),
			transferID:         transferID.String(),
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

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					RetryTransferInitiation(gomock.Any(), testCase.transferID).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					RetryTransferInitiation(gomock.Any(), testCase.transferID).
					Return(testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/transfer-initiations/%s/retry", testCase.transferID), nil)
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

func TestDeleteTransferInitiation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		transferID         string
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	transferID := models.TransferInitiationID{
		Reference: "ref1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:       "nominal",
			transferID: transferID.String(),
		},
		{
			name:               "service error duplicate key",
			serviceError:       storage.ErrDuplicateKeyValue,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "service error not found",
			serviceError:       storage.ErrNotFound,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "service error validation",
			serviceError:       service.ErrValidation,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "service error invalid ID",
			serviceError:       service.ErrInvalidID,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "service error publish",
			serviceError:       service.ErrPublish,
			transferID:         transferID.String(),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
		{
			name:               "service error other errors",
			serviceError:       errors.New("some error"),
			transferID:         transferID.String(),
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

			backend, mockService := newServiceTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					DeleteTransferInitiation(gomock.Any(), testCase.transferID).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					DeleteTransferInitiation(gomock.Any(), testCase.transferID).
					Return(testCase.serviceError)
			}

			router := httpRouter(logging.Testing(), backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/transfer-initiations/%s", testCase.transferID), nil)
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
