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
	"github.com/formancehq/stack/libs/go-libs/api"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListTransferInitiations(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListTransferInitiationsQuery
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name: "nomimal",
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
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
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
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
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
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
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
					WithPageSize(15),
			),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			expectedQuery: storage.NewListTransferInitiationsQuery(
				storage.NewPaginatedQueryOptions(storage.TransferInitiationQuery{}).
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

			tfs := []models.TransferInitiation{
				{
					ID: models.TransferInitiationID{
						Reference:   "t1",
						ConnectorID: connectorID,
					},
					CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					ScheduledAt: time.Date(2023, 11, 22, 8, 30, 0, 0, time.UTC),
					Description: "test1",
					Type:        models.TransferInitiationTypePayout,
					SourceAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					DestinationAccountID: models.AccountID{
						Reference:   "acc2",
						ConnectorID: connectorID,
					},
					Provider:    models.ConnectorProviderDummyPay,
					ConnectorID: connectorID,
					Amount:      big.NewInt(100),
					Asset:       models.Asset("EUR/2"),
					RelatedAdjustments: []*models.TransferInitiationAdjustment{
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t1",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 8, 45, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusProcessed,
						},
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t1",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 8, 40, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusWaitingForValidation,
						},
					},
					RelatedPayments: []*models.TransferInitiationPayment{
						{
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t1",
								ConnectorID: connectorID,
							},
							PaymentID: models.PaymentID{
								PaymentReference: models.PaymentReference{
									Reference: "p1",
									Type:      models.PaymentTypePayIn,
								},
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 8, 30, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusProcessed,
							Error:     "",
						},
					},
					Metadata: map[string]string{
						"foo": "bar",
					},
				},
				{
					ID: models.TransferInitiationID{
						Reference:   "t2",
						ConnectorID: connectorID,
					},
					CreatedAt:   time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
					ScheduledAt: time.Date(2023, 11, 22, 9, 30, 0, 0, time.UTC),
					Description: "test2",
					Type:        models.TransferInitiationTypeTransfer,
					SourceAccountID: &models.AccountID{
						Reference:   "acc3",
						ConnectorID: connectorID,
					},
					DestinationAccountID: models.AccountID{
						Reference:   "acc4",
						ConnectorID: connectorID,
					},
					Provider:    models.ConnectorProviderDummyPay,
					ConnectorID: connectorID,
					Amount:      big.NewInt(2000),
					Asset:       models.Asset("USD/2"),
					RelatedAdjustments: []*models.TransferInitiationAdjustment{
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t2",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 9, 45, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusFailed,
							Error:     "error",
						},
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t2",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 9, 40, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusWaitingForValidation,
						},
					},
				},
			}
			listTFsResponse := &api.Cursor[models.TransferInitiation]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     tfs,
			}

			expectedTFsResponse := []*transferInitiationResponse{
				{
					ID:                   tfs[0].ID.String(),
					Reference:            tfs[0].ID.Reference,
					CreatedAt:            tfs[0].CreatedAt,
					ScheduledAt:          tfs[0].ScheduledAt,
					Description:          tfs[0].Description,
					SourceAccountID:      tfs[0].SourceAccountID.String(),
					DestinationAccountID: tfs[0].DestinationAccountID.String(),
					Provider:             tfs[0].Provider.String(),
					Type:                 tfs[0].Type.String(),
					Amount:               tfs[0].Amount,
					Asset:                tfs[0].Asset.String(),
					Status:               models.TransferInitiationStatusProcessed.String(),
					ConnectorID:          tfs[0].ConnectorID.String(),
					Error:                "",
					Metadata:             tfs[0].Metadata,
				},
				{
					ID:                   tfs[1].ID.String(),
					Reference:            tfs[1].ID.Reference,
					CreatedAt:            tfs[1].CreatedAt,
					ScheduledAt:          tfs[1].ScheduledAt,
					Description:          tfs[1].Description,
					SourceAccountID:      tfs[1].SourceAccountID.String(),
					DestinationAccountID: tfs[1].DestinationAccountID.String(),
					Provider:             tfs[1].Provider.String(),
					Type:                 tfs[1].Type.String(),
					Amount:               tfs[1].Amount,
					Asset:                tfs[1].Asset.String(),
					ConnectorID:          tfs[1].ConnectorID.String(),
					Status:               models.TransferInitiationStatusFailed.String(),
					Error:                "error",
					Metadata:             tfs[1].Metadata,
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListTransferInitiations(gomock.Any(), testCase.expectedQuery).
					Return(listTFsResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListTransferInitiations(gomock.Any(), testCase.expectedQuery).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, "/transfer-initiations", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*transferInitiationResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedTFsResponse, resp.Cursor.Data)
				require.Equal(t, listTFsResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listTFsResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listTFsResponse.Next, resp.Cursor.Next)
				require.Equal(t, listTFsResponse.Previous, resp.Cursor.Previous)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestGetTransferInitiation(t *testing.T) {
	t.Parallel()

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
	tfID1 := models.TransferInitiationID{
		Reference:   "t1",
		ConnectorID: connectorID,
	}
	tfID2 := models.TransferInitiationID{
		Reference:   "t2",
		ConnectorID: connectorID,
	}

	type testCase struct {
		name               string
		tfID               string
		expectedTFID       models.TransferInitiationID
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name:         "nomimal acc1",
			tfID:         tfID1.String(),
			expectedTFID: tfID1,
		},
		{
			name:         "nomimal acc2",
			tfID:         tfID2.String(),
			expectedTFID: tfID2,
		},
		{
			name:               "invalid tf ID",
			tfID:               "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:               "err validation from backend",
			tfID:               tfID1.String(),
			expectedTFID:       tfID1,
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			tfID:               tfID1.String(),
			expectedTFID:       tfID1,
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			tfID:               tfID1.String(),
			expectedTFID:       tfID1,
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			tfID:               tfID1.String(),
			expectedTFID:       tfID1,
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

			var getTransferInitiationResponse *models.TransferInitiation
			var expectedTransferInitiationResponse *readTransferInitiationResponse
			if testCase.expectedTFID == tfID1 {
				getTransferInitiationResponse = &models.TransferInitiation{
					ID: models.TransferInitiationID{
						Reference:   "t1",
						ConnectorID: connectorID,
					},
					CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					ScheduledAt: time.Date(2023, 11, 22, 8, 30, 0, 0, time.UTC),
					Description: "test1",
					Type:        models.TransferInitiationTypePayout,
					SourceAccountID: &models.AccountID{
						Reference:   "acc1",
						ConnectorID: connectorID,
					},
					DestinationAccountID: models.AccountID{
						Reference:   "acc2",
						ConnectorID: connectorID,
					},
					Provider:    models.ConnectorProviderDummyPay,
					ConnectorID: connectorID,
					Amount:      big.NewInt(100),
					Asset:       models.Asset("EUR/2"),
					RelatedAdjustments: []*models.TransferInitiationAdjustment{
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t1",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 8, 45, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusProcessed,
						},
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t1",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 8, 40, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusWaitingForValidation,
						},
					},
					RelatedPayments: []*models.TransferInitiationPayment{
						{
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t1",
								ConnectorID: connectorID,
							},
							PaymentID: models.PaymentID{
								PaymentReference: models.PaymentReference{
									Reference: "p1",
									Type:      models.PaymentTypePayIn,
								},
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 8, 30, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusProcessed,
							Error:     "",
						},
					},
					Metadata: map[string]string{
						"foo": "bar",
					},
				}

				expectedTransferInitiationResponse = &readTransferInitiationResponse{
					transferInitiationResponse: transferInitiationResponse{
						ID:                   getTransferInitiationResponse.ID.String(),
						Reference:            getTransferInitiationResponse.ID.Reference,
						CreatedAt:            getTransferInitiationResponse.CreatedAt,
						ScheduledAt:          getTransferInitiationResponse.ScheduledAt,
						Description:          getTransferInitiationResponse.Description,
						SourceAccountID:      getTransferInitiationResponse.SourceAccountID.String(),
						DestinationAccountID: getTransferInitiationResponse.DestinationAccountID.String(),
						Provider:             getTransferInitiationResponse.Provider.String(),
						Type:                 getTransferInitiationResponse.Type.String(),
						Amount:               getTransferInitiationResponse.Amount,
						ConnectorID:          getTransferInitiationResponse.ConnectorID.String(),
						Asset:                getTransferInitiationResponse.Asset.String(),
						Status:               models.TransferInitiationStatusProcessed.String(),
						Error:                "",
						Metadata:             getTransferInitiationResponse.Metadata,
					},
					RelatedPayments: []*transferInitiationPaymentsResponse{
						{
							PaymentID: getTransferInitiationResponse.RelatedPayments[0].PaymentID.String(),
							CreatedAt: getTransferInitiationResponse.RelatedPayments[0].CreatedAt,
							Status:    getTransferInitiationResponse.RelatedPayments[0].Status.String(),
							Error:     getTransferInitiationResponse.RelatedPayments[0].Error,
						},
					},
					RelatedAdjustments: []*transferInitiationAdjustmentsResponse{
						{
							AdjustmentID: getTransferInitiationResponse.RelatedAdjustments[0].ID.String(),
							CreatedAt:    getTransferInitiationResponse.RelatedAdjustments[0].CreatedAt,
							Status:       getTransferInitiationResponse.RelatedAdjustments[0].Status.String(),
							Error:        getTransferInitiationResponse.RelatedAdjustments[0].Error,
							Metadata:     getTransferInitiationResponse.RelatedAdjustments[0].Metadata,
						},
						{
							AdjustmentID: getTransferInitiationResponse.RelatedAdjustments[1].ID.String(),
							CreatedAt:    getTransferInitiationResponse.RelatedAdjustments[1].CreatedAt,
							Status:       getTransferInitiationResponse.RelatedAdjustments[1].Status.String(),
							Error:        getTransferInitiationResponse.RelatedAdjustments[1].Error,
							Metadata:     getTransferInitiationResponse.RelatedAdjustments[1].Metadata,
						},
					},
				}
			} else {
				getTransferInitiationResponse = &models.TransferInitiation{
					ID: models.TransferInitiationID{
						Reference:   "t2",
						ConnectorID: connectorID,
					},
					CreatedAt:   time.Date(2023, 11, 22, 9, 0, 0, 0, time.UTC),
					ScheduledAt: time.Date(2023, 11, 22, 9, 30, 0, 0, time.UTC),
					Description: "test2",
					Type:        models.TransferInitiationTypeTransfer,
					SourceAccountID: &models.AccountID{
						Reference:   "acc3",
						ConnectorID: connectorID,
					},
					DestinationAccountID: models.AccountID{
						Reference:   "acc4",
						ConnectorID: connectorID,
					},
					Provider:    models.ConnectorProviderDummyPay,
					ConnectorID: connectorID,
					Amount:      big.NewInt(2000),
					Asset:       models.Asset("USD/2"),
					RelatedAdjustments: []*models.TransferInitiationAdjustment{
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t2",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 9, 45, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusFailed,
							Error:     "error",
						},
						{
							ID: uuid.New(),
							TransferInitiationID: models.TransferInitiationID{
								Reference:   "t2",
								ConnectorID: connectorID,
							},
							CreatedAt: time.Date(2023, 11, 22, 9, 40, 0, 0, time.UTC),
							Status:    models.TransferInitiationStatusWaitingForValidation,
						},
					},
				}
				expectedTransferInitiationResponse = &readTransferInitiationResponse{
					transferInitiationResponse: transferInitiationResponse{
						ID:                   getTransferInitiationResponse.ID.String(),
						Reference:            getTransferInitiationResponse.ID.Reference,
						CreatedAt:            getTransferInitiationResponse.CreatedAt,
						ScheduledAt:          getTransferInitiationResponse.ScheduledAt,
						Description:          getTransferInitiationResponse.Description,
						SourceAccountID:      getTransferInitiationResponse.SourceAccountID.String(),
						DestinationAccountID: getTransferInitiationResponse.DestinationAccountID.String(),
						Provider:             getTransferInitiationResponse.Provider.String(),
						Type:                 getTransferInitiationResponse.Type.String(),
						Amount:               getTransferInitiationResponse.Amount,
						ConnectorID:          getTransferInitiationResponse.ConnectorID.String(),
						Asset:                getTransferInitiationResponse.Asset.String(),
						Status:               models.TransferInitiationStatusFailed.String(),
						Error:                "error",
						Metadata:             getTransferInitiationResponse.Metadata,
					},
					RelatedAdjustments: []*transferInitiationAdjustmentsResponse{
						{
							AdjustmentID: getTransferInitiationResponse.RelatedAdjustments[0].ID.String(),
							CreatedAt:    getTransferInitiationResponse.RelatedAdjustments[0].CreatedAt,
							Status:       getTransferInitiationResponse.RelatedAdjustments[0].Status.String(),
							Error:        getTransferInitiationResponse.RelatedAdjustments[0].Error,
							Metadata:     getTransferInitiationResponse.RelatedAdjustments[0].Metadata,
						},
						{
							AdjustmentID: getTransferInitiationResponse.RelatedAdjustments[1].ID.String(),
							CreatedAt:    getTransferInitiationResponse.RelatedAdjustments[1].CreatedAt,
							Status:       getTransferInitiationResponse.RelatedAdjustments[1].Status.String(),
							Error:        getTransferInitiationResponse.RelatedAdjustments[1].Error,
							Metadata:     getTransferInitiationResponse.RelatedAdjustments[1].Metadata,
						},
					},
				}
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ReadTransferInitiation(gomock.Any(), testCase.expectedTFID).
					Return(getTransferInitiationResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ReadTransferInitiation(gomock.Any(), testCase.expectedTFID).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/transfer-initiations/%s", testCase.tfID), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[readTransferInitiationResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedTransferInitiationResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
