package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/api"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreatePayments(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.CreatePaymentRequest
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
			name: "nomimal",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
		},
		{
			name: "no source account id, but should still pass",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				DestinationAccountID: destinationAccountID.String(),
			},
		},
		{
			name: "no destination account id, but should still pass",
			req: &service.CreatePaymentRequest{
				Reference:       "test",
				ConnectorID:     connectorID.String(),
				CreatedAt:       time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:          big.NewInt(100),
				Type:            string(models.PaymentTypeTransfer),
				Status:          string(models.PaymentStatusSucceeded),
				Scheme:          string(models.PaymentSchemeOther),
				Asset:           "EUR/2",
				SourceAccountID: sourceAccountID.String(),
			},
		},
		{
			name: "missing reference",
			req: &service.CreatePaymentRequest{
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing createdAt",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "created at to zero",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Time{},
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing connectorID",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing amount",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing type",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "invalid type",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 "invalid",
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing status",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "invalid status",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               "invalid",
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing scheme",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "invalid scheme",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               "invalid",
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing asset",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "invalid asset",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "invalid",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "err validation from backend",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			req: &service.CreatePaymentRequest{
				Reference:            "test",
				ConnectorID:          connectorID.String(),
				CreatedAt:            time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Amount:               big.NewInt(100),
				Type:                 string(models.PaymentTypeTransfer),
				Status:               string(models.PaymentStatusSucceeded),
				Scheme:               string(models.PaymentSchemeOther),
				Asset:                "EUR/2",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
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

			createPaymentResponse := &models.Payment{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: testCase.req.Reference,
						Type:      models.PaymentTypeTransfer,
					},
					ConnectorID: connectorID,
				},
				ConnectorID:          connectorID,
				CreatedAt:            testCase.req.CreatedAt,
				Reference:            testCase.req.Reference,
				Amount:               testCase.req.Amount,
				InitialAmount:        testCase.req.Amount,
				Type:                 models.PaymentTypeTransfer,
				Status:               models.PaymentStatusSucceeded,
				Scheme:               models.PaymentSchemeOther,
				Asset:                models.Asset("EUR/2"),
				SourceAccountID:      &sourceAccountID,
				DestinationAccountID: &destinationAccountID,
			}

			expectedCreatePaymentResponse := &paymentResponse{
				ID: models.PaymentID{
					PaymentReference: models.PaymentReference{
						Reference: "test",
						Type:      models.PaymentTypeTransfer,
					},
					ConnectorID: connectorID,
				}.String(),
				Reference:            "test",
				SourceAccountID:      sourceAccountID.String(),
				DestinationAccountID: destinationAccountID.String(),
				Type:                 string(createPaymentResponse.Type),
				Provider:             createPaymentResponse.ConnectorID.Provider,
				ConnectorID:          createPaymentResponse.ConnectorID.String(),
				Status:               createPaymentResponse.Status,
				InitialAmount:        createPaymentResponse.Amount,
				Amount:               createPaymentResponse.Amount,
				Scheme:               createPaymentResponse.Scheme,
				Asset:                createPaymentResponse.Asset.String(),
				CreatedAt:            createPaymentResponse.CreatedAt,
				Adjustments:          make([]paymentAdjustment, len(createPaymentResponse.Adjustments)),
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					CreatePayment(gomock.Any(), testCase.req).
					Return(createPaymentResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					CreatePayment(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/payments", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[paymentResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedCreatePaymentResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestPayments(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListPaymentsQuery
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name: "nomimal",
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(100),
			),
			pageSize: 100,
		},
		{
			name: "with invalid page size",
			queryParams: url.Values{
				"pageSize": {"nan"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "invalid query builder json",
			queryParams: url.Values{
				"query": {"invalid"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "valid query builder json",
			queryParams: url.Values{
				"query": {"{\"$match\": {\"source_account_id\": \"acc1\"}}"},
			},
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15).
					WithQueryBuilder(query.Match("source_account_id", "acc1")),
			),
			pageSize: 15,
		},
		{
			name: "valid sorter",
			queryParams: url.Values{
				"sort": {"source_account_id:asc"},
			},
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15).
					WithSorter(storage.Sorter{}.Add("source_account_id", storage.SortOrderAsc)),
			),
			pageSize: 15,
		},
		{
			name: "invalid sorter",
			queryParams: url.Values{
				"sort": {"source_account_id:invalid"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "err validation from backend",
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15),
			),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			expectedQuery: storage.NewListPaymentsQuery(
				storage.NewPaginatedQueryOptions(storage.PaymentQuery{}).
					WithPageSize(15),
			),
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

			connectorID := models.ConnectorID{
				Reference: uuid.New(),
				Provider:  models.ConnectorProviderDummyPay,
			}

			payments := []models.Payment{
				{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: "p1",
							Type:      models.PaymentTypePayIn,
						},
						ConnectorID: connectorID,
					},
					ConnectorID:   connectorID,
					CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference:     "p1",
					Amount:        big.NewInt(100),
					InitialAmount: big.NewInt(1000),
					Type:          models.PaymentTypePayIn,
					Status:        models.PaymentStatusPending,
					Scheme:        models.PaymentSchemeCardMasterCard,
					Asset:         models.Asset("USD/2"),
					SourceAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				},
				{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: "p2",
							Type:      models.PaymentTypePayOut,
						},
						ConnectorID: connectorID,
					},
					ConnectorID:   connectorID,
					CreatedAt:     time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
					Reference:     "p2",
					Amount:        big.NewInt(1000),
					InitialAmount: big.NewInt(10000),
					Type:          models.PaymentTypePayOut,
					Status:        models.PaymentStatusSucceeded,
					Scheme:        models.PaymentSchemeCardVisa,
					Asset:         models.Asset("EUR/2"),
					DestinationAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				},
			}
			listPaymentsResponse := &api.Cursor[models.Payment]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     payments,
			}

			expectedPaymentsResponse := []*paymentResponse{
				{
					ID:                   payments[0].ID.String(),
					Reference:            payments[0].Reference,
					SourceAccountID:      payments[0].SourceAccountID.String(),
					DestinationAccountID: payments[0].DestinationAccountID.String(),
					Type:                 payments[0].Type.String(),
					Provider:             payments[0].Connector.Provider,
					ConnectorID:          payments[0].ConnectorID.String(),
					Status:               payments[0].Status,
					InitialAmount:        payments[0].InitialAmount,
					Amount:               payments[0].Amount,
					Scheme:               payments[0].Scheme,
					Asset:                payments[0].Asset.String(),
					CreatedAt:            payments[0].CreatedAt,
					Adjustments:          make([]paymentAdjustment, len(payments[0].Adjustments)),
				},
				{
					ID:                   payments[1].ID.String(),
					Reference:            payments[1].Reference,
					SourceAccountID:      payments[1].SourceAccountID.String(),
					DestinationAccountID: payments[1].DestinationAccountID.String(),
					Type:                 payments[1].Type.String(),
					Provider:             payments[1].Connector.Provider,
					ConnectorID:          payments[1].ConnectorID.String(),
					Status:               payments[1].Status,
					InitialAmount:        payments[1].InitialAmount,
					Amount:               payments[1].Amount,
					Scheme:               payments[1].Scheme,
					Asset:                payments[1].Asset.String(),
					CreatedAt:            payments[1].CreatedAt,
					Adjustments:          make([]paymentAdjustment, len(payments[0].Adjustments)),
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListPayments(gomock.Any(), testCase.expectedQuery).
					Return(listPaymentsResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListPayments(gomock.Any(), testCase.expectedQuery).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, "/payments", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*paymentResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPaymentsResponse, resp.Cursor.Data)
				require.Equal(t, listPaymentsResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listPaymentsResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listPaymentsResponse.Next, resp.Cursor.Next)
				require.Equal(t, listPaymentsResponse.Previous, resp.Cursor.Previous)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestGetPayment(t *testing.T) {
	t.Parallel()

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
	paymentID1 := models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "p1",
			Type:      models.PaymentTypePayIn,
		},
		ConnectorID: connectorID,
	}
	paymentID2 := models.PaymentID{
		PaymentReference: models.PaymentReference{
			Reference: "p2",
			Type:      models.PaymentTypePayOut,
		},
		ConnectorID: connectorID,
	}

	type testCase struct {
		name               string
		paymentID          string
		expectedPaymentID  models.PaymentID
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name:              "nomimal p1",
			paymentID:         paymentID1.String(),
			expectedPaymentID: paymentID1,
		},
		{
			name:              "nomimal p2",
			paymentID:         paymentID2.String(),
			expectedPaymentID: paymentID2,
		},
		{
			name:               "err validation from backend",
			paymentID:          paymentID1.String(),
			expectedPaymentID:  paymentID1,
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			paymentID:          paymentID1.String(),
			expectedPaymentID:  paymentID1,
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			paymentID:          paymentID1.String(),
			expectedPaymentID:  paymentID1,
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			paymentID:          paymentID1.String(),
			expectedPaymentID:  paymentID1,
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

			var getPaymentResponse *models.Payment
			var expectedPaymentResponse *paymentResponse
			if testCase.expectedPaymentID == paymentID1 {
				getPaymentResponse = &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: "p1",
							Type:      models.PaymentTypePayIn,
						},
						ConnectorID: connectorID,
					},
					ConnectorID:   connectorID,
					CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference:     "p1",
					Amount:        big.NewInt(100),
					InitialAmount: big.NewInt(1000),
					Type:          models.PaymentTypePayIn,
					Status:        models.PaymentStatusPending,
					Scheme:        models.PaymentSchemeCardMasterCard,
					Asset:         models.Asset("USD/2"),
					SourceAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				}

				expectedPaymentResponse = &paymentResponse{
					ID:                   getPaymentResponse.ID.String(),
					Reference:            getPaymentResponse.Reference,
					SourceAccountID:      getPaymentResponse.SourceAccountID.String(),
					DestinationAccountID: getPaymentResponse.DestinationAccountID.String(),
					Type:                 getPaymentResponse.Type.String(),
					Provider:             getPaymentResponse.Connector.Provider,
					ConnectorID:          getPaymentResponse.ConnectorID.String(),
					Status:               getPaymentResponse.Status,
					InitialAmount:        getPaymentResponse.InitialAmount,
					Amount:               getPaymentResponse.Amount,
					Scheme:               getPaymentResponse.Scheme,
					Asset:                getPaymentResponse.Asset.String(),
					CreatedAt:            getPaymentResponse.CreatedAt,
					Adjustments:          make([]paymentAdjustment, len(getPaymentResponse.Adjustments)),
				}
			} else {
				getPaymentResponse = &models.Payment{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: "p2",
							Type:      models.PaymentTypePayOut,
						},
						ConnectorID: connectorID,
					},
					ConnectorID:   connectorID,
					CreatedAt:     time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
					Reference:     "p2",
					Amount:        big.NewInt(1000),
					InitialAmount: big.NewInt(10000),
					Type:          models.PaymentTypePayOut,
					Status:        models.PaymentStatusSucceeded,
					Scheme:        models.PaymentSchemeCardVisa,
					Asset:         models.Asset("EUR/2"),
					DestinationAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				}
				expectedPaymentResponse = &paymentResponse{
					ID:                   getPaymentResponse.ID.String(),
					Reference:            getPaymentResponse.Reference,
					SourceAccountID:      getPaymentResponse.SourceAccountID.String(),
					DestinationAccountID: getPaymentResponse.DestinationAccountID.String(),
					Type:                 getPaymentResponse.Type.String(),
					Provider:             getPaymentResponse.Connector.Provider,
					ConnectorID:          getPaymentResponse.ConnectorID.String(),
					Status:               getPaymentResponse.Status,
					InitialAmount:        getPaymentResponse.InitialAmount,
					Amount:               getPaymentResponse.Amount,
					Scheme:               getPaymentResponse.Scheme,
					Asset:                getPaymentResponse.Asset.String(),
					CreatedAt:            getPaymentResponse.CreatedAt,
					Adjustments:          make([]paymentAdjustment, len(getPaymentResponse.Adjustments)),
				}
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetPayment(gomock.Any(), testCase.expectedPaymentID.String()).
					Return(getPaymentResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetPayment(gomock.Any(), testCase.expectedPaymentID.String()).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/payments/%s", testCase.paymentID), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[paymentResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPaymentResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
