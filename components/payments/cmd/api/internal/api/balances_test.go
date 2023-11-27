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

func TestGetBalances(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		accountID          string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.BalanceQuery
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithLimit(10).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithFrom(time.Date(2023, 11, 20, 6, 0, 0, 0, time.UTC)).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
				"to":       []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
				"to":       []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(100, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, query.Match("account_id", "acc1"))).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "valid sorter",
			queryParams: url.Values{
				"sort": {"account_id:asc"},
				"to":   []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, storage.Sorter{}.Add("account_id", storage.SortOrderAsc), nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			name: "valid cursor",
			queryParams: url.Values{
				"cursor": {cursor},
				"to":     []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewBalanceQuery(
				storage.NewPaginatorQuery(15, nil, nil).
					WithToken(cursor).
					WithCursor(storage.NewBaseCursor("test", nil, false)),
			).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
			pageSize:  15,
			accountID: accountIDString,
		},
		{
			name: "err validation from backend",
			queryParams: url.Values{
				"to": []string{time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)},
			},
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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
			expectedQuery: storage.NewBalanceQuery(storage.NewPaginatorQuery(15, nil, nil)).
				WithAccountID(&accountID).
				WithTo(time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC)),
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

			listBalancesResponse := []*models.Balance{
				{
					AccountID:     accountID,
					Asset:         "EUR/2",
					Balance:       big.NewInt(100),
					CreatedAt:     time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC),
					LastUpdatedAt: time.Date(2023, 11, 23, 9, 0, 0, 0, time.UTC),
					ConnectorID:   connectorID,
				},
			}

			expectedBalancessResponse := []*balancesResponse{
				{
					AccountID:     listBalancesResponse[0].AccountID.String(),
					CreatedAt:     listBalancesResponse[0].CreatedAt,
					LastUpdatedAt: listBalancesResponse[0].LastUpdatedAt,
					Currency:      listBalancesResponse[0].Asset.String(),
					Asset:         listBalancesResponse[0].Asset.String(),
					Balance:       listBalancesResponse[0].Balance,
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
					ListBalances(gomock.Any(), testCase.expectedQuery).
					Return(listBalancesResponse, expectedPaginationDetails, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListBalances(gomock.Any(), testCase.expectedQuery).
					Return(nil, storage.PaginationDetails{}, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{})

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/accounts/%s/balances", testCase.accountID), nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*balancesResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedBalancessResponse, resp.Cursor.Data)
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
