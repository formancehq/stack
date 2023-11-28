package api

import (
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
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestPayments(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.PaginatorQuery
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name:          "nomimal",
			expectedQuery: storage.NewPaginatorQuery(15, nil, nil),
			pageSize:      15,
		},
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewPaginatorQuery(15, nil, nil),
			pageSize:      15,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewPaginatorQuery(100, nil, nil),
			pageSize:      100,
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
			expectedQuery: storage.NewPaginatorQuery(15, nil, query.Match("source_account_id", "acc1")),
			pageSize:      15,
		},
		{
			name: "valid sorter",
			queryParams: url.Values{
				"sort": {"source_account_id:asc"},
			},
			expectedQuery: storage.NewPaginatorQuery(15, storage.Sorter{}.Add("source_account_id", storage.SortOrderAsc), nil),
			pageSize:      15,
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
			name: "valid cursor",
			queryParams: url.Values{
				"cursor": {cursor},
			},
			expectedQuery: storage.NewPaginatorQuery(15, nil, nil).
				WithToken(cursor).
				WithCursor(storage.NewBaseCursor("test", nil, false)),
			pageSize: 15,
		},
		{
			name:               "err validation from backend",
			expectedQuery:      storage.NewPaginatorQuery(15, nil, nil),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			expectedQuery:      storage.NewPaginatorQuery(15, nil, nil),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			expectedQuery:      storage.NewPaginatorQuery(15, nil, nil),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			expectedQuery:      storage.NewPaginatorQuery(15, nil, nil),
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

			listPaymentsResponse := []*models.Payment{
				{
					ID: models.PaymentID{
						PaymentReference: models.PaymentReference{
							Reference: "p1",
							Type:      models.PaymentTypePayIn,
						},
						ConnectorID: connectorID,
					},
					ConnectorID: connectorID,
					CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference:   "p1",
					Amount:      big.NewInt(100),
					Type:        models.PaymentTypePayIn,
					Status:      models.PaymentStatusPending,
					Scheme:      models.PaymentSchemeCardMasterCard,
					Asset:       models.Asset("USD/2"),
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
					ConnectorID: connectorID,
					CreatedAt:   time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
					Reference:   "p2",
					Amount:      big.NewInt(1000),
					Type:        models.PaymentTypePayOut,
					Status:      models.PaymentStatusSucceeded,
					Scheme:      models.PaymentSchemeCardVisa,
					Asset:       models.Asset("EUR/2"),
					DestinationAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				},
			}

			expectedPaymentsResponse := []*paymentResponse{
				{
					ID:                   listPaymentsResponse[0].ID.String(),
					Reference:            listPaymentsResponse[0].Reference,
					SourceAccountID:      listPaymentsResponse[0].SourceAccountID.String(),
					DestinationAccountID: listPaymentsResponse[0].DestinationAccountID.String(),
					Type:                 listPaymentsResponse[0].Type.String(),
					Provider:             listPaymentsResponse[0].Connector.Provider,
					ConnectorID:          listPaymentsResponse[0].ConnectorID.String(),
					Status:               listPaymentsResponse[0].Status,
					InitialAmount:        listPaymentsResponse[0].Amount,
					Scheme:               listPaymentsResponse[0].Scheme,
					Asset:                listPaymentsResponse[0].Asset.String(),
					CreatedAt:            listPaymentsResponse[0].CreatedAt,
					Adjustments:          make([]paymentAdjustment, len(listPaymentsResponse[0].Adjustments)),
				},
				{
					ID:                   listPaymentsResponse[1].ID.String(),
					Reference:            listPaymentsResponse[1].Reference,
					SourceAccountID:      listPaymentsResponse[1].SourceAccountID.String(),
					DestinationAccountID: listPaymentsResponse[1].DestinationAccountID.String(),
					Type:                 listPaymentsResponse[1].Type.String(),
					Provider:             listPaymentsResponse[1].Connector.Provider,
					ConnectorID:          listPaymentsResponse[1].ConnectorID.String(),
					Status:               listPaymentsResponse[1].Status,
					InitialAmount:        listPaymentsResponse[1].Amount,
					Scheme:               listPaymentsResponse[1].Scheme,
					Asset:                listPaymentsResponse[1].Asset.String(),
					CreatedAt:            listPaymentsResponse[1].CreatedAt,
					Adjustments:          make([]paymentAdjustment, len(listPaymentsResponse[0].Adjustments)),
				},
			}
			expectedPaginationDetails := storage.PaginationDetails{
				PageSize:     testCase.pageSize,
				HasMore:      false,
				PreviousPage: "",
				NextPage:     "",
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListPayments(gomock.Any(), testCase.expectedQuery).
					Return(listPaymentsResponse, expectedPaginationDetails, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListPayments(gomock.Any(), testCase.expectedQuery).
					Return(nil, storage.PaginationDetails{}, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{})

			req := httptest.NewRequest(http.MethodGet, "/payments", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*paymentResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPaymentsResponse, resp.Cursor.Data)
				require.Equal(t, expectedPaginationDetails.PageSize, resp.Cursor.PageSize)
				require.Equal(t, expectedPaginationDetails.HasMore, resp.Cursor.HasMore)
				require.Equal(t, expectedPaginationDetails.NextPage, resp.Cursor.Next)
				require.Equal(t, expectedPaginationDetails.PreviousPage, resp.Cursor.Previous)
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
					ConnectorID: connectorID,
					CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference:   "p1",
					Amount:      big.NewInt(100),
					Type:        models.PaymentTypePayIn,
					Status:      models.PaymentStatusPending,
					Scheme:      models.PaymentSchemeCardMasterCard,
					Asset:       models.Asset("USD/2"),
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
					InitialAmount:        getPaymentResponse.Amount,
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
					ConnectorID: connectorID,
					CreatedAt:   time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
					Reference:   "p2",
					Amount:      big.NewInt(1000),
					Type:        models.PaymentTypePayOut,
					Status:      models.PaymentStatusSucceeded,
					Scheme:      models.PaymentSchemeCardVisa,
					Asset:       models.Asset("EUR/2"),
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
					InitialAmount:        getPaymentResponse.Amount,
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

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{})

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
