package api

import (
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
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

func TestGetBalances(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		accountID          string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListBalancesQuery
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}
	accountID := models.AccountID{
		Reference:   "acc1",
		ConnectorID: connectorID,
	}
	accountIDString := accountID.String()
	testCases := []testCase{
		{
			name: "nomimal",
			queryParams: url.Values{
				"to": []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name:               "with invalid accountID",
			accountID:          "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "with valid limit",
			queryParams: url.Values{
				"limit": {"10"},
				"to":    []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(10),
			),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "with invalid limit",
			queryParams: url.Values{
				"limit": {"nan"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "with from and to",
			queryParams: url.Values{
				"from": []string{time.Date(2023, 11, 20, 6, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
				"to":   []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithFrom(time.Date(2023, 11, 20, 6, 0, 0, 0, time.UTC)).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			accountID: accountIDString,
		},
		{
			name: "with invalid from",
			queryParams: url.Values{
				"from": []string{"invalid"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "with invalid to",
			queryParams: url.Values{
				"to": []string{"invalid"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "page size too low, should use the default value",
			queryParams: url.Values{
				"pageSize": {"0"},
				"to":       []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
				"to":       []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(100),
			),
			pageSize:  100,
			accountID: accountIDString,
		},
		{
			name: "with invalid page size",
			queryParams: url.Values{
				"pageSize": {"nan"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "invalid query builder json",
			queryParams: url.Values{
				"query": {"invalid"},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "valid query builder json",
			queryParams: url.Values{
				"query": {"{\"$match\": {\"account_id\": \"acc1\"}}"},
				"to":    []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15).
					WithQueryBuilder(query.Match("account_id", "acc1")),
			),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "valid sorter",
			queryParams: url.Values{
				"sort": {"account_id:asc"},
				"to":   []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15).
					WithSorter(storage.Sorter{}.Add("account_id", storage.SortOrderAsc)),
			),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "invalid sorter",
			queryParams: url.Values{
				"sort": {"account_id:invalid"},
				"to":   []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "err validation from backend",
			queryParams: url.Values{
				"to": []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
			accountID:          accountIDString,
		},
		{
			name: "ErrNotFound from storage",
			queryParams: url.Values{
				"to": []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
			accountID:          accountIDString,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			queryParams: url.Values{
				"to": []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
			accountID:          accountIDString,
		},
		{
			name: "other storage errors from storage",
			queryParams: url.Values{
				"to": []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewListBalancesQuery(
				storage.NewPaginatedQueryOptions(
					storage.NewBalanceQuery().
						WithAccountID(&accountID).
						WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
				).WithPageSize(15),
			),
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
			accountID:          accountIDString,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			balances := []models.Balance{
				{
					AccountID:     accountID,
					Asset:         "EUR/2",
					Balance:       big.NewInt(100),
					CreatedAt:     time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC),
					LastUpdatedAt: time.Date(2023, 11, 23, 9, 0, 0, 0, time.UTC),
					ConnectorID:   connectorID,
				},
			}

			listBalancesResponse := &api.Cursor[models.Balance]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     balances,
			}

			if limit, ok := testCase.queryParams["limit"]; ok {
				testCase.pageSize, _ = strconv.Atoi(limit[0])
			}

			expectedBalancessResponse := []*balancesResponse{
				{
					AccountID:     balances[0].AccountID.String(),
					CreatedAt:     balances[0].CreatedAt,
					LastUpdatedAt: balances[0].LastUpdatedAt,
					Currency:      balances[0].Asset.String(),
					Asset:         balances[0].Asset.String(),
					Balance:       balances[0].Balance,
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListBalances(gomock.Any(), testCase.expectedQuery).
					Return(listBalancesResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListBalances(gomock.Any(), testCase.expectedQuery).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/balances", testCase.accountID), nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*balancesResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedBalancessResponse, resp.Cursor.Data)
				require.Equal(t, listBalancesResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listBalancesResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listBalancesResponse.Next, resp.Cursor.Next)
				require.Equal(t, listBalancesResponse.Previous, resp.Cursor.Previous)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
