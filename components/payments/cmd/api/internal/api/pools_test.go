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

func TestCreatePool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.CreatePoolRequest
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	accountID := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	uuid1 := uuid.New()

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.CreatePoolRequest{
				Name:       "test",
				AccountIDs: []string{accountID.String()},
			},
		},
		{
			name: "no accounts",
			req: &service.CreatePoolRequest{
				Name:       "test",
				AccountIDs: []string{},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "missing name",
			req: &service.CreatePoolRequest{
				Name:       "",
				AccountIDs: []string{accountID.String()},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "err validation from backend",
			req: &service.CreatePoolRequest{
				Name:       "test",
				AccountIDs: []string{accountID.String()},
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			req: &service.CreatePoolRequest{
				Name:       "test",
				AccountIDs: []string{accountID.String()},
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			req: &service.CreatePoolRequest{
				Name:       "test",
				AccountIDs: []string{accountID.String()},
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			req: &service.CreatePoolRequest{
				Name:       "test",
				AccountIDs: []string{accountID.String()},
			},
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			createPoolResponse := &models.Pool{
				ID:   uuid1,
				Name: testCase.req.Name,
				PoolAccounts: []*models.PoolAccounts{
					{
						PoolID:    uuid1,
						AccountID: accountID,
					},
				},
			}

			accounts := make([]string, len(createPoolResponse.PoolAccounts))
			for i := range createPoolResponse.PoolAccounts {
				accounts[i] = createPoolResponse.PoolAccounts[i].AccountID.String()
			}
			expectedCreatePoolResponse := &poolResponse{
				ID:       createPoolResponse.ID.String(),
				Name:     createPoolResponse.Name,
				Accounts: accounts,
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					CreatePool(gomock.Any(), testCase.req).
					Return(createPoolResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					CreatePool(gomock.Any(), testCase.req).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, "/pools", bytes.NewReader(body))
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[poolResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedCreatePoolResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestAddAccountToPool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		req                *service.AddAccountToPoolRequest
		poolID             string
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	accountID := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	uuid1 := uuid.New()

	testCases := []testCase{
		{
			name: "nominal",
			req: &service.AddAccountToPoolRequest{
				AccountID: accountID.String(),
			},
			poolID: uuid1.String(),
		},
		{
			name: "missing accountID",
			req: &service.AddAccountToPoolRequest{
				AccountID: "",
			},
			poolID:             uuid1.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "missing body",
			poolID:             uuid1.String(),
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrMissingOrInvalidBody,
		},
		{
			name:   "err validation from backend",
			poolID: uuid1.String(),
			req: &service.AddAccountToPoolRequest{
				AccountID: accountID.String(),
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:   "ErrNotFound from storage",
			poolID: uuid1.String(),
			req: &service.AddAccountToPoolRequest{
				AccountID: accountID.String(),
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:   "ErrDuplicateKeyValue from storage",
			poolID: uuid1.String(),
			req: &service.AddAccountToPoolRequest{
				AccountID: accountID.String(),
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:   "other storage errors from storage",
			poolID: uuid1.String(),
			req: &service.AddAccountToPoolRequest{
				AccountID: accountID.String(),
			},
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					AddAccountToPool(gomock.Any(), testCase.poolID, testCase.req).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					AddAccountToPool(gomock.Any(), testCase.poolID, testCase.req).
					Return(testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			var body []byte
			if testCase.req != nil {
				var err error
				body, err = json.Marshal(testCase.req)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/pools/%s/accounts", testCase.poolID), bytes.NewReader(body))
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

func TestRemoveAccountFromPool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		poolID             string
		accountID          string
		serviceError       error
		expectedStatusCode int
		expectedErrorCode  string
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
	accountID := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}
	uuid1 := uuid.New()

	testCases := []testCase{
		{
			name:      "nominal",
			poolID:    uuid1.String(),
			accountID: accountID.String(),
		},
		{
			name:               "err validation from backend",
			poolID:             uuid1.String(),
			accountID:          accountID.String(),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			poolID:             uuid1.String(),
			accountID:          accountID.String(),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			poolID:             uuid1.String(),
			accountID:          accountID.String(),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			poolID:             uuid1.String(),
			accountID:          accountID.String(),
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
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					RemoveAccountFromPool(gomock.Any(), testCase.poolID, testCase.accountID).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					RemoveAccountFromPool(gomock.Any(), testCase.poolID, testCase.accountID).
					Return(testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/pools/%s/accounts/%s", testCase.poolID, testCase.accountID), nil)
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

func TestListPools(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListPoolsQuery
		expectedStatusCode int
		serviceError       error
		expectedErrorCode  string
	}

	testCases := []testCase{
		{
			name: "nomimal",
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
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
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
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
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
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
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
					WithPageSize(15),
			),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			expectedQuery: storage.NewListPoolsQuery(
				storage.NewPaginatedQueryOptions(storage.PoolQuery{}).
					WithPageSize(15),
			),
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	poolID1 := uuid.New()
	poolID2 := uuid.New()

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	accountID1 := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	accountID2 := models.AccountID{
		Reference:   "acc2",
		ConnectorID: connectorID,
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			pools := []models.Pool{
				{
					ID:   poolID1,
					Name: "test1",
					PoolAccounts: []*models.PoolAccounts{
						{
							PoolID:    poolID1,
							AccountID: accountID1,
						},
					},
				},
				{
					ID:   poolID2,
					Name: "test2",
					PoolAccounts: []*models.PoolAccounts{
						{
							PoolID:    poolID2,
							AccountID: accountID1,
						},
						{
							PoolID:    poolID2,
							AccountID: accountID2,
						},
					},
				},
			}
			listPoolsResponse := &api.Cursor[models.Pool]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     pools,
			}

			accounts1 := make([]string, len(pools[0].PoolAccounts))
			for i := range pools[0].PoolAccounts {
				accounts1[i] = pools[0].PoolAccounts[i].AccountID.String()
			}

			accounts2 := make([]string, len(pools[1].PoolAccounts))
			for i := range pools[1].PoolAccounts {
				accounts2[i] = pools[1].PoolAccounts[i].AccountID.String()
			}
			expectedListPoolsResponse := []*poolResponse{
				{
					ID:       pools[0].ID.String(),
					Name:     pools[0].Name,
					Accounts: accounts1,
				},
				{
					ID:       pools[1].ID.String(),
					Name:     pools[1].Name,
					Accounts: accounts2,
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListPools(gomock.Any(), testCase.expectedQuery).
					Return(listPoolsResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListPools(gomock.Any(), testCase.expectedQuery).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, "/pools", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*poolResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedListPoolsResponse, resp.Cursor.Data)
				require.Equal(t, listPoolsResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listPoolsResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listPoolsResponse.Next, resp.Cursor.Next)
				require.Equal(t, listPoolsResponse.Previous, resp.Cursor.Previous)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestGetPool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		poolID             string
		serviceError       error
		expectedPoolID     uuid.UUID
		expectedStatusCode int
		expectedErrorCode  string
	}

	uuid1 := uuid.New()
	testCases := []testCase{
		{
			name:           "nomimal",
			poolID:         uuid1.String(),
			expectedPoolID: uuid1,
		},
		{
			name:               "err validation from backend",
			poolID:             uuid1.String(),
			expectedPoolID:     uuid1,
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			poolID:             uuid1.String(),
			expectedPoolID:     uuid1,
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			poolID:             uuid1.String(),
			expectedPoolID:     uuid1,
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			poolID:             uuid1.String(),
			expectedPoolID:     uuid1,
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	accountID1 := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			getPoolResponse := &models.Pool{
				ID:   uuid1,
				Name: "test1",
				PoolAccounts: []*models.PoolAccounts{
					{
						PoolID:    uuid1,
						AccountID: accountID1,
					},
				},
			}

			accounts := make([]string, len(getPoolResponse.PoolAccounts))
			for i := range getPoolResponse.PoolAccounts {
				accounts[i] = getPoolResponse.PoolAccounts[i].AccountID.String()
			}
			expectedPoolResponse := &poolResponse{
				ID:       uuid1.String(),
				Name:     getPoolResponse.Name,
				Accounts: accounts,
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetPool(gomock.Any(), testCase.poolID).
					Return(getPoolResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetPool(gomock.Any(), testCase.poolID).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/pools/%s", testCase.poolID), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[poolResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPoolResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestGetPoolBalance(t *testing.T) {
	t.Parallel()

	uuid1 := uuid.New()
	type testCase struct {
		name               string
		queryParams        url.Values
		poolID             string
		serviceError       error
		expectedStatusCode int
		expectedErrorCode  string
	}

	atTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC).UTC()
	testCases := []testCase{
		{
			name:   "nominal",
			poolID: uuid1.String(),
			queryParams: url.Values{
				"at": {atTime.Format(time.RFC3339)},
			},
		},
		{
			name:   "missing at",
			poolID: uuid1.String(),

			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:   "err validation from backend",
			poolID: uuid1.String(),
			queryParams: url.Values{
				"at": {atTime.Format(time.RFC3339)},
			},
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:   "ErrNotFound from storage",
			poolID: uuid1.String(),
			queryParams: url.Values{
				"at": {atTime.Format(time.RFC3339)},
			},
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:   "ErrDuplicateKeyValue from storage",
			poolID: uuid1.String(),
			queryParams: url.Values{
				"at": {atTime.Format(time.RFC3339)},
			},
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:   "other storage errors from storage",
			poolID: uuid1.String(),
			queryParams: url.Values{
				"at": {atTime.Format(time.RFC3339)},
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

			getPoolBalanceResponse := &service.GetPoolBalanceResponse{
				Balances: []*service.Balance{
					{
						Amount: big.NewInt(100),
						Asset:  "EUR/2",
					},
					{
						Amount: big.NewInt(12000),
						Asset:  "USD/2",
					},
				},
			}

			expectedPoolBalancesResponse := &poolBalancesResponse{
				Balances: []*poolBalanceResponse{
					{
						Amount: getPoolBalanceResponse.Balances[0].Amount,
						Asset:  getPoolBalanceResponse.Balances[0].Asset,
					},
					{
						Amount: getPoolBalanceResponse.Balances[1].Amount,
						Asset:  getPoolBalanceResponse.Balances[1].Asset,
					},
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetPoolBalance(gomock.Any(), testCase.poolID, atTime.Format(time.RFC3339)).
					Return(getPoolBalanceResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetPoolBalance(gomock.Any(), testCase.poolID, atTime.Format(time.RFC3339)).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/pools/%s/balances", testCase.poolID), nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[poolBalancesResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedPoolBalancesResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}

}

func TestDeletePool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		poolID             string
		serviceError       error
		expectedStatusCode int
		expectedErrorCode  string
	}

	uuid1 := uuid.New()
	testCases := []testCase{
		{
			name:   "nominal",
			poolID: uuid1.String(),
		},
		{
			name:               "err validation from backend",
			poolID:             uuid1.String(),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name:               "ErrNotFound from storage",
			poolID:             uuid1.String(),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name:               "ErrDuplicateKeyValue from storage",
			poolID:             uuid1.String(),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name:               "other storage errors from storage",
			poolID:             uuid1.String(),
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
				testCase.expectedStatusCode = http.StatusNoContent
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					DeletePool(gomock.Any(), testCase.poolID).
					Return(nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					DeletePool(gomock.Any(), testCase.poolID).
					Return(testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/pools/%s", testCase.poolID), nil)
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
