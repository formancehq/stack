package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	sharedapi "github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/bun/bunpaginate"
	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/query"
	"github.com/formancehq/payments/cmd/api/internal/api/service"
	"github.com/formancehq/payments/cmd/api/internal/storage"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCreateAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.CreateAccountRequest
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
			name: "nomimal",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "no default asset, but should still pass",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "missing reference",
			req: &service.CreateAccountRequest{
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing connectorID",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing createdAt",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "createdAt zero",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Time{},
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing accountName",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing type",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "invalid type",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "unknown",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "err validation from backend",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			req: &service.CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  connectorID.String(),
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
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
				testCase.expectedStatusCode = http.StatusOK
			}

			createAccountResponse := &models.Account{
				ID: models.AccountID{
					Reference:   testCase.req.Reference,
					ConnectorID: connectorID,
				},
				ConnectorID:  connectorID,
				CreatedAt:    testCase.req.CreatedAt,
				Reference:    testCase.req.Reference,
				DefaultAsset: models.Asset(testCase.req.DefaultAsset),
				AccountName:  testCase.req.AccountName,
				Type:         models.AccountType(testCase.req.Type),
				Metadata: map[string]string{
					"foo": "bar",
				},
				PoolAccounts: make([]*models.PoolAccounts, 0),
			}

			expectedCreateAccountResponse := &accountResponse{
				ID:              createAccountResponse.ID.String(),
				Reference:       createAccountResponse.Reference,
				CreatedAt:       createAccountResponse.CreatedAt,
				ConnectorID:     createAccountResponse.ConnectorID.String(),
				Provider:        createAccountResponse.ConnectorID.Provider.String(),
				DefaultCurrency: createAccountResponse.DefaultAsset.String(),
				DefaultAsset:    createAccountResponse.DefaultAsset.String(),
				AccountName:     createAccountResponse.AccountName,
				Type:            createAccountResponse.Type.String(),
				Metadata: map[string]string{
					"foo": "bar",
				},
				Pools: make([]uuid.UUID, 0),
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					CreateAccount(gomock.Any(), testCase.req).
					Return(createAccountResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					CreateAccount(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), false)

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[accountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedCreateAccountResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestListAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListAccountsQuery
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name: "nomimal",
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
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
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
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
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
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
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			expectedQuery: storage.NewListAccountsQuery(
				storage.NewPaginatedQueryOptions(storage.AccountQuery{}).
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

			accounts := []models.Account{
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
				},
			}

			listAccountsResponse := &bunpaginate.Cursor[models.Account]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     accounts,
			}

			expectedAccountsResponse := []*accountResponse{
				{
					ID:              accounts[0].ID.String(),
					Reference:       accounts[0].Reference,
					CreatedAt:       accounts[0].CreatedAt,
					ConnectorID:     accounts[0].ConnectorID.String(),
					Provider:        accounts[0].ConnectorID.Provider.String(),
					DefaultCurrency: accounts[0].DefaultAsset.String(),
					DefaultAsset:    accounts[0].DefaultAsset.String(),
					AccountName:     accounts[0].AccountName,
					Type:            accounts[0].Type.String(),
					Pools:           []uuid.UUID{},
					Metadata:        accounts[0].Metadata,
				},
				{
					ID:              accounts[1].ID.String(),
					Reference:       accounts[1].Reference,
					CreatedAt:       accounts[1].CreatedAt,
					ConnectorID:     accounts[1].ConnectorID.String(),
					Provider:        accounts[1].ConnectorID.Provider.String(),
					DefaultCurrency: accounts[1].DefaultAsset.String(),
					DefaultAsset:    accounts[1].DefaultAsset.String(),
					AccountName:     accounts[1].AccountName,
					Pools:           []uuid.UUID{},
					// Type is converted to external when it is external formance
					Type: string(models.AccountTypeExternal),
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListAccounts(gomock.Any(), testCase.expectedQuery).
					Return(listAccountsResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListAccounts(gomock.Any(), testCase.expectedQuery).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), false)

			req := httptest.NewRequest(http.MethodGet, "/accounts", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*accountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedAccountsResponse, resp.Cursor.Data)
				require.Equal(t, listAccountsResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listAccountsResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listAccountsResponse.Next, resp.Cursor.Next)
				require.Equal(t, listAccountsResponse.Previous, resp.Cursor.Previous)
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

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{
				Debug: testing.Verbose(),
			}, auth.NewNoAuth(), false)

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
