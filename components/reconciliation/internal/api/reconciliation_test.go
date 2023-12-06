package api

import (
	"bytes"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/formancehq/reconciliation/internal/api/service"
	"github.com/formancehq/reconciliation/internal/models"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	gomock "github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestReconciliation(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.ReconciliationRequest
		res                *service.ReconciliationResponse
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.ReconciliationRequest{
				LedgerName:                   "default",
				LedgerAggregatedBalanceQuery: map[string]interface{}{},
				PaymentPoolID:                uuid.New().String(),
				At:                           time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
		},
		{
			name: "missing ledger name",
			req: &service.ReconciliationRequest{
				LedgerAggregatedBalanceQuery: map[string]interface{}{},
				PaymentPoolID:                uuid.New().String(),
				At:                           time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing ledger aggregated balance query",
			req: &service.ReconciliationRequest{
				LedgerName:    "default",
				PaymentPoolID: uuid.New().String(),
				At:            time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing payment pool id",
			req: &service.ReconciliationRequest{
				LedgerName:                   "default",
				LedgerAggregatedBalanceQuery: map[string]interface{}{},
				At:                           time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing at",
			req: &service.ReconciliationRequest{
				LedgerName:                   "default",
				LedgerAggregatedBalanceQuery: map[string]interface{}{},
				PaymentPoolID:                uuid.New().String(),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "zero time.Time",
			req: &service.ReconciliationRequest{
				LedgerName:                   "default",
				LedgerAggregatedBalanceQuery: map[string]interface{}{},
				PaymentPoolID:                uuid.New().String(),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error validation",
			req: &service.ReconciliationRequest{
				LedgerName:                   "default",
				LedgerAggregatedBalanceQuery: map[string]interface{}{},
				PaymentPoolID:                uuid.New().String(),
				At:                           time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			res: &service.ReconciliationResponse{
				Status:          models.ReconciliationOK,
				PaymentBalances: map[string]*big.Int{},
				LedgerBalances:  map[string]*big.Int{},
				Error:           "",
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
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
				Status:          testCase.res.Status.String(),
				PaymentBalances: testCase.res.PaymentBalances,
				LedgerBalances:  testCase.res.LedgerBalances,
				Error:           testCase.res.Error,
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					Reconciliation(gomock.Any(), testCase.req).
					Return(testCase.res, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					Reconciliation(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := newRouter(backend, sharedapi.ServiceInfo{}, nil)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/reconciliation", bytes.NewReader(body))
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
