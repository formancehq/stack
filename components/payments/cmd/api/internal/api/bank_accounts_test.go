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
	"github.com/formancehq/stack/libs/go-libs/api"
	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/query"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListBankAccounts(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name               string
		queryParams        url.Values
		pageSize           int
		expectedQuery      storage.ListBankAccountQuery
		expectedStatusCode int
		expectedErrorCode  string
		serviceError       error
	}

	testCases := []testCase{
		{
			name: "nomimal",
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too low",
			queryParams: url.Values{
				"pageSize": {"0"},
			},
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
					WithPageSize(15),
			),
			pageSize: 15,
		},
		{
			name: "page size too high",
			queryParams: url.Values{
				"pageSize": {"100000"},
			},
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
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
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
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
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
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
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       service.ErrValidation,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrValidation,
		},
		{
			name: "ErrNotFound from storage",
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrNotFound,
			expectedStatusCode: http.StatusNotFound,
			expectedErrorCode:  ErrNotFound,
		},
		{
			name: "ErrDuplicateKeyValue from storage",
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       storage.ErrDuplicateKeyValue,
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrUniqueReference,
		},
		{
			name: "other storage errors from storage",
			expectedQuery: storage.NewListBankAccountQuery(
				storage.NewPaginatedQueryOptions(storage.BankAccountQuery{}).
					WithPageSize(15),
			),
			serviceError:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedErrorCode:  sharedapi.ErrorInternal,
		},
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if testCase.expectedStatusCode == 0 {
				testCase.expectedStatusCode = http.StatusOK
			}

			b1ID := uuid.New()
			b2ID := uuid.New()

			bankAccounts := []models.BankAccount{
				{
					ID:            b1ID,
					CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
					Name:          "ba1",
					AccountNumber: "0112345678",
					IBAN:          "FR7630006000011234567890189",
					SwiftBicCode:  "HBUKGB4B",
					Country:       "FR",
					RelatedAccounts: []*models.BankAccountRelatedAccount{
						{
							ID:            uuid.New(),
							BankAccountID: b1ID,
							ConnectorID:   connectorID,
							AccountID: models.AccountID{
								Reference:   "acc1",
								ConnectorID: connectorID,
							},
						},
					},
				},
				{
					ID:            b2ID,
					CreatedAt:     time.Date(2023, 11, 23, 8, 0, 0, 0, time.UTC),
					Name:          "ba2",
					AccountNumber: "0112345679",
					IBAN:          "FR7630006000011234567890188",
					SwiftBicCode:  "ABCDGB4B",
					Country:       "DE",
					RelatedAccounts: []*models.BankAccountRelatedAccount{
						{
							ID:            uuid.New(),
							BankAccountID: b2ID,
							ConnectorID:   connectorID,
							AccountID: models.AccountID{
								Reference:   "acc2",
								ConnectorID: connectorID,
							},
						},
					},
				},
			}
			listBankAccountsResponse := &api.Cursor[models.BankAccount]{
				PageSize: testCase.pageSize,
				HasMore:  false,
				Previous: "",
				Next:     "",
				Data:     bankAccounts,
			}

			expectedBankAccountsResponse := []*bankAccountResponse{
				{
					ID:          bankAccounts[0].ID.String(),
					Name:        bankAccounts[0].Name,
					CreatedAt:   bankAccounts[0].CreatedAt,
					Country:     bankAccounts[0].Country,
					ConnectorID: bankAccounts[0].RelatedAccounts[0].ConnectorID.String(),
					AccountID:   bankAccounts[0].RelatedAccounts[0].AccountID.String(),
					Provider:    bankAccounts[0].RelatedAccounts[0].ConnectorID.Provider.String(),
					RelatedAccounts: []*bankAccountRelatedAccountsResponse{
						{
							ID:          bankAccounts[0].RelatedAccounts[0].ID.String(),
							AccountID:   bankAccounts[0].RelatedAccounts[0].AccountID.String(),
							ConnectorID: bankAccounts[0].RelatedAccounts[0].ConnectorID.String(),
							Provider:    bankAccounts[0].RelatedAccounts[0].ConnectorID.Provider.String(),
						},
					},
				},
				{
					ID:          bankAccounts[1].ID.String(),
					Name:        bankAccounts[1].Name,
					CreatedAt:   bankAccounts[1].CreatedAt,
					Country:     bankAccounts[1].Country,
					ConnectorID: bankAccounts[1].RelatedAccounts[0].ConnectorID.String(),
					AccountID:   bankAccounts[1].RelatedAccounts[0].AccountID.String(),
					Provider:    bankAccounts[1].RelatedAccounts[0].ConnectorID.Provider.String(),
					RelatedAccounts: []*bankAccountRelatedAccountsResponse{
						{
							ID:          bankAccounts[1].RelatedAccounts[0].ID.String(),
							AccountID:   bankAccounts[1].RelatedAccounts[0].AccountID.String(),
							ConnectorID: bankAccounts[1].RelatedAccounts[0].ConnectorID.String(),
							Provider:    bankAccounts[1].RelatedAccounts[0].ConnectorID.Provider.String(),
						},
					},
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					ListBankAccounts(gomock.Any(), testCase.expectedQuery).
					Return(listBankAccountsResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					ListBankAccounts(gomock.Any(), testCase.expectedQuery).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, "/bank-accounts", nil)
			rec := httptest.NewRecorder()
			req.URL.RawQuery = testCase.queryParams.Encode()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[*bankAccountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedBankAccountsResponse, resp.Cursor.Data)
				require.Equal(t, listBankAccountsResponse.PageSize, resp.Cursor.PageSize)
				require.Equal(t, listBankAccountsResponse.HasMore, resp.Cursor.HasMore)
				require.Equal(t, listBankAccountsResponse.Next, resp.Cursor.Next)
				require.Equal(t, listBankAccountsResponse.Previous, resp.Cursor.Previous)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}

func TestGetBankAccount(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                    string
		bankAccountUUID         string
		expectedBankAccountUUID uuid.UUID
		expectedStatusCode      int
		expectedErrorCode       string
		serviceError            error
	}

	uuid1 := uuid.New()
	testCases := []testCase{
		{
			name:                    "nomimal",
			bankAccountUUID:         uuid1.String(),
			expectedBankAccountUUID: uuid1,
		},
		{
			name:               "invalid uuid",
			bankAccountUUID:    "invalid",
			expectedStatusCode: http.StatusBadRequest,
			expectedErrorCode:  ErrInvalidID,
		},
		{
			name:                    "err validation from backend",
			bankAccountUUID:         uuid1.String(),
			expectedBankAccountUUID: uuid1,
			serviceError:            service.ErrValidation,
			expectedStatusCode:      http.StatusBadRequest,
			expectedErrorCode:       ErrValidation,
		},
		{
			name:                    "ErrNotFound from storage",
			bankAccountUUID:         uuid1.String(),
			expectedBankAccountUUID: uuid1,
			serviceError:            storage.ErrNotFound,
			expectedStatusCode:      http.StatusNotFound,
			expectedErrorCode:       ErrNotFound,
		},
		{
			name:                    "ErrDuplicateKeyValue from storage",
			bankAccountUUID:         uuid1.String(),
			expectedBankAccountUUID: uuid1,
			serviceError:            storage.ErrDuplicateKeyValue,
			expectedStatusCode:      http.StatusBadRequest,
			expectedErrorCode:       ErrUniqueReference,
		},
		{
			name:                    "other storage errors from storage",
			bankAccountUUID:         uuid1.String(),
			expectedBankAccountUUID: uuid1,
			serviceError:            errors.New("some error"),
			expectedStatusCode:      http.StatusInternalServerError,
			expectedErrorCode:       sharedapi.ErrorInternal,
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
			getBankAccountResponse := &models.BankAccount{
				ID:            uuid1,
				CreatedAt:     time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				Name:          "ba1",
				AccountNumber: "13719713158835300",
				IBAN:          "FR7630006000011234567890188",
				SwiftBicCode:  "ABCDGB4B",
				Country:       "FR",
				RelatedAccounts: []*models.BankAccountRelatedAccount{
					{
						ID:            uuid.New(),
						BankAccountID: uuid1,
						ConnectorID:   connectorID,
						AccountID: models.AccountID{
							Reference:   "acc1",
							ConnectorID: connectorID,
						},
					},
				},
			}

			expectedBankAccountResponse := &bankAccountResponse{
				ID:            getBankAccountResponse.ID.String(),
				Name:          getBankAccountResponse.Name,
				CreatedAt:     getBankAccountResponse.CreatedAt,
				Country:       getBankAccountResponse.Country,
				ConnectorID:   getBankAccountResponse.RelatedAccounts[0].ConnectorID.String(),
				Provider:      getBankAccountResponse.RelatedAccounts[0].ConnectorID.Provider.String(),
				AccountID:     getBankAccountResponse.RelatedAccounts[0].AccountID.String(),
				Iban:          "FR76*******************0188",
				AccountNumber: "13************300",
				SwiftBicCode:  "ABCDGB4B",
				RelatedAccounts: []*bankAccountRelatedAccountsResponse{
					{
						ID:          getBankAccountResponse.RelatedAccounts[0].ID.String(),
						AccountID:   getBankAccountResponse.RelatedAccounts[0].AccountID.String(),
						ConnectorID: getBankAccountResponse.RelatedAccounts[0].ConnectorID.String(),
						Provider:    getBankAccountResponse.RelatedAccounts[0].ConnectorID.Provider.String(),
					},
				},
			}

			backend, mockService := newTestingBackend(t)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				mockService.EXPECT().
					GetBankAccount(gomock.Any(), testCase.expectedBankAccountUUID, true).
					Return(getBankAccountResponse, nil)
			}
			if testCase.serviceError != nil {
				mockService.EXPECT().
					GetBankAccount(gomock.Any(), testCase.expectedBankAccountUUID, true).
					Return(nil, testCase.serviceError)
			}

			router := httpRouter(backend, logging.Testing(), sharedapi.ServiceInfo{}, auth.NewNoAuth())

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/bank-accounts/%s", testCase.bankAccountUUID), nil)
			rec := httptest.NewRecorder()

			router.ServeHTTP(rec, req)

			require.Equal(t, testCase.expectedStatusCode, rec.Code)
			if testCase.expectedStatusCode < 300 && testCase.expectedStatusCode >= 200 {
				var resp sharedapi.BaseResponse[bankAccountResponse]
				sharedapi.Decode(t, rec.Body, &resp)
				require.Equal(t, expectedBankAccountResponse, resp.Data)
			} else {
				err := sharedapi.ErrorResponse{}
				sharedapi.Decode(t, rec.Body, &err)
				require.EqualValues(t, testCase.expectedErrorCode, err.ErrorCode)
			}
		})
	}
}
