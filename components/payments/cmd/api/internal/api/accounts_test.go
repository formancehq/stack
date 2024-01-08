package api

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.PaginatorQuery
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
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

			listAccountsResponse := []*models.Account{
				{
					ID: models.AccountID{Reference: "acc1", ConnectorID: models.ConnectorID{Reference: uuid.New(), Provider: models.ConnectorProviderDummyPay}},
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
					CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference: "acc1",
					Type:      models.AccountTypeInternal,
					Metadata: map[string]string{
						"foo": "bar",
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				},
				{
					ID: models.AccountID{Reference: "acc2", ConnectorID: models.ConnectorID{Reference: uuid.New(), Provider: models.ConnectorProviderDummyPay}},
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
					CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference: "acc2",
					Type:      models.AccountTypeExternalFormance,
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				},
			}

			expectedAccountsResponse := []*accountResponse{
				{
					ID:              listAccountsResponse[0].ID.String(),
					Reference:       listAccountsResponse[0].Reference,
					CreatedAt:       listAccountsResponse[0].CreatedAt,
					ConnectorID:     listAccountsResponse[0].ConnectorID.String(),
					Provider:        listAccountsResponse[0].ConnectorID.Provider.String(),
					DefaultCurrency: listAccountsResponse[0].DefaultAsset.String(),
					DefaultAsset:    listAccountsResponse[0].DefaultAsset.String(),
					AccountName:     listAccountsResponse[0].AccountName,
					Type:            listAccountsResponse[0].Type.String(),
					Pools:           []uuid.UUID{},
					Metadata:        listAccountsResponse[0].Metadata,
				},
				{
					ID:              listAccountsResponse[1].ID.String(),
					Reference:       listAccountsResponse[1].Reference,
					CreatedAt:       listAccountsResponse[1].CreatedAt,
					ConnectorID:     listAccountsResponse[1].ConnectorID.String(),
					Provider:        listAccountsResponse[1].ConnectorID.Provider.String(),
					DefaultCurrency: listAccountsResponse[1].DefaultAsset.String(),
					DefaultAsset:    listAccountsResponse[1].DefaultAsset.String(),
					AccountName:     listAccountsResponse[1].AccountName,
					Pools:           []uuid.UUID{},
					// Type is converted to external when it is external formance
					Type: string(models.AccountTypeExternal),
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
					ListAccounts(gomock.Any(), testCase.expectedQuery).
					Return(listAccountsResponse, expectedPaginationDetails, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListAccounts(gomock.Any(), testCase.expectedQuery).
					Return(nil, storage.PaginationDetails{}, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*accountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedAccountsResponse, resp.Cursor.Data)
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

func TestGetAccount(t *testing.T) {
	t.Parallel()

	accountID1 := models.AccountID{
		Reference: "acc1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}
	accountID2 := models.AccountID{
		Reference: "acc2",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	type testCase struct {
		name               string
		accountID          string
		serviceError       error
		expectedAccountID  models.AccountID
		expectedStatusCode int
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name:              "nomimal acc1",
			accountID:         accountID1.String(),
			expectedAccountID: accountID1,
		},
		{
			name:              "nomimal acc2",
			accountID:         accountID2.String(),
			expectedAccountID: accountID2,
		},
		{
			name:               "err validation from backend",
			accountID:          accountID1.String(),
			expectedAccountID:  accountID1,
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			accountID:          accountID1.String(),
			expectedAccountID:  accountID1,
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			accountID:          accountID1.String(),
			expectedAccountID:  accountID1,
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			accountID:          accountID1.String(),
			expectedAccountID:  accountID1,
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

			var getAccountResponse *models.Account
			var expectedAccountsResponse *accountResponse
			if testCase.expectedAccountID == accountID1 {
				getAccountResponse = &models.Account{
					ID: models.AccountID{Reference: "acc1", ConnectorID: models.ConnectorID{Reference: uuid.New(), Provider: models.ConnectorProviderDummyPay}},
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
					CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference: "acc1",
					Type:      models.AccountTypeInternal,
					Metadata: map[string]string{
						"foo": "bar",
					},
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				}

				expectedAccountsResponse = &accountResponse{
					ID:              getAccountResponse.ID.String(),
					Reference:       getAccountResponse.Reference,
					CreatedAt:       getAccountResponse.CreatedAt,
					ConnectorID:     getAccountResponse.ConnectorID.String(),
					Provider:        getAccountResponse.ConnectorID.Provider.String(),
					DefaultCurrency: getAccountResponse.DefaultAsset.String(),
					DefaultAsset:    getAccountResponse.DefaultAsset.String(),
					AccountName:     getAccountResponse.AccountName,
					Metadata:        getAccountResponse.Metadata,
					Pools:           []uuid.UUID{},
					Type:            getAccountResponse.Type.String(),
				}
			} else {
				getAccountResponse = &models.Account{
					ID: models.AccountID{Reference: "acc2", ConnectorID: models.ConnectorID{Reference: uuid.New(), Provider: models.ConnectorProviderDummyPay}},
					ConnectorID: models.ConnectorID{
						Reference: uuid.New(),
						Provider:  models.ConnectorProviderDummyPay,
					},
					CreatedAt: time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Reference: "acc2",
					Type:      models.AccountTypeExternalFormance,
					Connector: &models.Connector{
						Provider: models.ConnectorProviderDummyPay,
					},
				}
				expectedAccountsResponse = &accountResponse{
					ID:              getAccountResponse.ID.String(),
					Reference:       getAccountResponse.Reference,
					CreatedAt:       getAccountResponse.CreatedAt,
					ConnectorID:     getAccountResponse.ConnectorID.String(),
					Provider:        getAccountResponse.ConnectorID.Provider.String(),
					DefaultCurrency: getAccountResponse.DefaultAsset.String(),
					DefaultAsset:    getAccountResponse.DefaultAsset.String(),
					AccountName:     getAccountResponse.AccountName,
					Pools:           []uuid.UUID{},
					// Type is converted to external when it is external formance
					Type: models.AccountTypeExternal.String(),
				}
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetAccount(gomock.Any(), testCase.expectedAccountID.String()).
					Return(getAccountResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetAccount(gomock.Any(), testCase.expectedAccountID.String()).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s", testCase.accountID), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[accountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedAccountsResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
