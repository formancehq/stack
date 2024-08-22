package api

import (
	"bytes"
	"encoding/json"
	"errors"
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

func TestCreatePolicy(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.CreatePolicyRequest
		invalidBody        bool
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.CreatePolicyRequest{
				Name:           "test",
				LedgerName:     "test",
				LedgerQuery:    map[string]interface{}{},
				PaymentsPoolID: "00000000-0000-0000-0000-000000000000",
			},
		},
		{
			name:               "missing body",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:               "invalid body",
			invalidBody:        true,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name: "service error validation",
			req: &service.CreatePolicyRequest{
				Name:           "test",
				LedgerName:     "test",
				LedgerQuery:    map[string]interface{}{},
				PaymentsPoolID: "00000000-0000-0000-0000-000000000000",
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "service error invalid id",
			req: &service.CreatePolicyRequest{
				Name:           "test",
				LedgerName:     "test",
				LedgerQuery:    map[string]interface{}{},
				PaymentsPoolID: "00000000-0000-0000-0000-000000000000",
			},
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name: "storage error not found",
			req: &service.CreatePolicyRequest{
				Name:           "test",
				LedgerName:     "test",
				LedgerQuery:    map[string]interface{}{},
				PaymentsPoolID: "00000000-0000-0000-0000-000000000000",
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  sharedapi.ErrorCodeNotFound,
		},
		{
			name: "service error other error",
			req: &service.CreatePolicyRequest{
				Name:           "test",
				LedgerName:     "test",
				LedgerQuery:    map[string]interface{}{},
				PaymentsPoolID: "00000000-0000-0000-0000-000000000000",
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
				testCase.expectedStatusCode = http.StatusCreated
			}

			var policyServiceResponse models.Policy
			if testCase.req != nil {
				policyServiceResponse = models.Policy{
					ID:             uuid.New(),
					CreatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Name:           testCase.req.Name,
					LedgerName:     testCase.req.LedgerName,
					LedgerQuery:    testCase.req.LedgerQuery,
					PaymentsPoolID: uuid.MustParse(testCase.req.PaymentsPoolID),
				}
			}

			expectedPolicyResponse := &policyResponse{
				ID:             policyServiceResponse.ID.String(),
				Name:           policyServiceResponse.Name,
				CreatedAt:      policyServiceResponse.CreatedAt,
				LedgerName:     policyServiceResponse.LedgerName,
				LedgerQuery:    policyServiceResponse.LedgerQuery,
				PaymentsPoolID: policyServiceResponse.PaymentsPoolID.String(),
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					CreatePolicy(gomock.Any(), testCase.req).
					Return(&policyServiceResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					CreatePolicy(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := newRouter(backend, sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), nil)

			var body []byte
			if testCase.invalidBody {
				body = []byte("invalid")
			} else if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/policies", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[policyResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPolicyResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestDeletePolicy(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		policyID           string
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name:     "nominal",
			policyID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:               "service error validation",
			policyID:           "00000000-0000-0000-0000-000000000000",
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "service error invalid id",
			policyID:           "invalid",
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "storage error not found",
			policyID:           "invalid",
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  sharedapi.ErrorCodeNotFound,
		},
		{
			name:               "service error other error",
			policyID:           "00000000-0000-0000-0000-000000000000",
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
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					DeletePolicy(gomock.Any(), testCase.policyID).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					DeletePolicy(gomock.Any(), testCase.policyID).
					Return(testCase.serviceError)
			}

			router := newRouter(backend, sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), nil)

			req := httptest.NewRequest(http.MethodDelete, "/policies/"+testCase.policyID, nil)
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

func TestGetPolicy(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		policyID           string
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name:     "nominal",
			policyID: "00000000-0000-0000-0000-000000000000",
		},
		{
			name:               "service error validation",
			policyID:           "00000000-0000-0000-0000-000000000000",
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "service error invalid id",
			policyID:           "00000000-0000-0000-0000-000000000000",
			serviceError:       service.ErrInvalidID,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "storage error not found",
			policyID:           "00000000-0000-0000-0000-000000000000",
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  sharedapi.ErrorCodeNotFound,
		},
		{
			name:               "service error other error",
			policyID:           "00000000-0000-0000-0000-000000000000",
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

			var policyServiceResponse models.Policy
			if testCase.policyID != "" {
				policyServiceResponse = models.Policy{
					ID:             uuid.MustParse(testCase.policyID),
					CreatedAt:      time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
					Name:           "test",
					LedgerName:     "test",
					LedgerQuery:    map[string]interface{}{},
					PaymentsPoolID: uuid.New(),
				}
			}

			expectedPolicyResponse := &policyResponse{
				ID:             policyServiceResponse.ID.String(),
				Name:           policyServiceResponse.Name,
				CreatedAt:      policyServiceResponse.CreatedAt,
				LedgerName:     policyServiceResponse.LedgerName,
				LedgerQuery:    policyServiceResponse.LedgerQuery,
				PaymentsPoolID: policyServiceResponse.PaymentsPoolID.String(),
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetPolicy(gomock.Any(), testCase.policyID).
					Return(&policyServiceResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetPolicy(gomock.Any(), testCase.policyID).
					Return(nil, testCase.serviceError)
			}

			router := newRouter(backend, sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), nil)

			req := httptest.NewRequest(http.MethodGet, "/policies/"+testCase.policyID, nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[policyResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPolicyResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
