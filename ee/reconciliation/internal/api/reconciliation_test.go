package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/reconciliation/internal/api/service"
	"github.com/formancehq/reconciliation/internal/models"
	"github.com/formancehq/reconciliation/internal/storage"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestReconciliation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		policyID           string
		req                *service.ReconciliationRequest
		res                *models.Reconciliation
		invalidBody        bool
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	policyID := uuid.New()
	testCases := []testCase{
		{
			name:     "nominal",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
		},
		{
			name:     "missing body",
			policyID: policyID.String(),
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:     "invalid body",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			invalidBody:        true,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:     "missing at ledger",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:     "zero time.Time ledger",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:     "missing at payments",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                 uuid.New(),
				PolicyID:           policyID,
				CreatedAt:          time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:             models.ReconciliationOK,
				PaymentsBalances:   map[string]*big.Int{},
				LedgerBalances:     map[string]*big.Int{},
				Error:              "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:     "zero time.Time ledger",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Time{},
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:     "service error validation",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:     "service error invalid id",
			policyID: "invalid",
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:     "service error invalid ID",
			policyID: "invalid",
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:     "storage error not found",
			policyID: "invalid",
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  sharedapi.ErrorCodeNotFound,
		},
		{
			name:     "service error other error",
			policyID: policyID.String(),
			req: &service.ReconciliationRequest{
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &models.Reconciliation{
				ID:                   uuid.New(),
				PolicyID:             policyID,
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				PaymentsBalances:     map[string]*big.Int{},
				LedgerBalances:       map[string]*big.Int{},
				Error:                "",
			},
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			expectedReconciliationResponse := &reconciliationResponse{
				ID:                   testCase.res.ID.String(),
				PolicyID:             testCase.res.PolicyID.String(),
				CreatedAt:            testCase.res.CreatedAt,
				ReconciledAtLedger:   testCase.res.ReconciledAtLedger,
				ReconciledAtPayments: testCase.res.ReconciledAtPayments,
				Status:               testCase.res.Status.String(),
				PaymentsBalances:     testCase.res.PaymentsBalances,
				LedgerBalances:       testCase.res.LedgerBalances,
				Error:                testCase.res.Error,
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					Reconciliation(gomock.Any(), testCase.policyID, testCase.req).
					Return(testCase.res, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					Reconciliation(gomock.Any(), testCase.policyID, testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := newRouter(backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.invalidBody {
				body = []byte("invalid")
			} else if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/policies/%s/reconciliation", testCase.policyID), bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[reconciliationResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedReconciliationResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestGetReconciliation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		id                 uuid.UUID
		serviceError       error
		expectedStatusCode int
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name: "nominal",
			id:   uuid.New(),
		},
		{
			name:               "service error validation",
			id:                 uuid.New(),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "service error invalid id",
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "storage error not found",
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  sharedapi.ErrorCodeNotFound,
		},
		{
			name:               "service error other error",
			id:                 uuid.New(),
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, tc := range testCases {
		testCase := tc
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			getReconciliationResponse := &models.Reconciliation{
				ID:                   testCase.id,
				PolicyID:             uuid.New(),
				CreatedAt:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtLedger:   time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				ReconciledAtPayments: time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
				Status:               models.ReconciliationOK,
				LedgerBalances: map[string]*big.Int{
					"USD/2": big.NewInt(100),
					"EUR/2": big.NewInt(200),
				},
				PaymentsBalances: map[string]*big.Int{
					"USD/2": big.NewInt(100),
					"EUR/2": big.NewInt(200),
				},
				Error: "",
			}

			expectedReconciliationResponse := &reconciliationResponse{
				ID:                   getReconciliationResponse.ID.String(),
				PolicyID:             getReconciliationResponse.PolicyID.String(),
				CreatedAt:            getReconciliationResponse.CreatedAt,
				ReconciledAtLedger:   getReconciliationResponse.ReconciledAtLedger,
				ReconciledAtPayments: getReconciliationResponse.ReconciledAtPayments,
				Status:               getReconciliationResponse.Status.String(),
				PaymentsBalances:     getReconciliationResponse.PaymentsBalances,
				LedgerBalances:       getReconciliationResponse.LedgerBalances,
				Error:                getReconciliationResponse.Error,
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetReconciliation(gomock.Any(), testCase.id.String()).
					Return(getReconciliationResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetReconciliation(gomock.Any(), testCase.id.String()).
					Return(nil, testCase.serviceError)
			}

			router := newRouter(backend, sharedapi.ServiceInfo{}, auth.NewNoAuth(), nil)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/reconciliations/%s", testCase.id.String()), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[reconciliationResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedReconciliationResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
