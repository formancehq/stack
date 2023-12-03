package service

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreatePool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		request       *CreatePoolRequest
		expectedError error
	}

	accountID := models.AccountID{
		Reference: "acc1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name: "nominal",
			request: &CreatePoolRequest{
				Name:       "pool1",
				AccountIDs: []string{accountID.String()},
			},
			expectedError: nil,
		},
		{
			name: "invalid accountID",
			request: &CreatePoolRequest{
				Name:       "pool1",
				AccountIDs: []string{"invalid"},
			},
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			pool, err := service.CreatePool(context.Background(), tc.request)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
				require.NotNil(t, pool)
			}
		})
	}
}

func TestGetPool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		poolID        string
		expectedError error
	}

	uuid1 := uuid.New()

	testCases := []testCase{
		{
			name:          "nominal",
			poolID:        uuid1.String(),
			expectedError: nil,
		},
		{
			name:          "invalid poolID",
			poolID:        "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.GetPool(context.Background(), tc.poolID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestAddAccountToPool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		poolID        string
		accountID     string
		expectedError error
	}

	uuid1 := uuid.New()
	accountID := models.AccountID{
		Reference: "acc1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:          "nominal",
			poolID:        uuid1.String(),
			accountID:     accountID.String(),
			expectedError: nil,
		},
		{
			name:          "invalid poolID",
			poolID:        "invalid",
			accountID:     accountID.String(),
			expectedError: ErrValidation,
		},
		{
			name:          "invalid accountID",
			poolID:        uuid1.String(),
			accountID:     "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := service.AddAccountToPool(context.Background(), tc.poolID, &AddAccountToPoolRequest{
				AccountID: tc.accountID,
			})
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestRemoveAccountFromPool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		poolID        string
		accountID     string
		expectedError error
	}

	uuid1 := uuid.New()
	accountID := models.AccountID{
		Reference: "acc1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:          "nominal",
			poolID:        uuid1.String(),
			accountID:     accountID.String(),
			expectedError: nil,
		},
		{
			name:          "invalid poolID",
			poolID:        "invalid",
			accountID:     accountID.String(),
			expectedError: ErrValidation,
		},
		{
			name:          "invalid accountID",
			poolID:        uuid1.String(),
			accountID:     "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := service.RemoveAccountFromPool(context.Background(), tc.poolID, tc.accountID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestDeletePool(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		poolID        string
		expectedError error
	}

	uuid1 := uuid.New()

	testCases := []testCase{
		{
			name:          "nominal",
			poolID:        uuid1.String(),
			expectedError: nil,
		},
		{
			name:          "invalid poolID",
			poolID:        "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			err := service.DeletePool(context.Background(), tc.poolID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}

}

func TestGetPoolBalance(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name          string
		poolID        string
		atTime        string
		expectedError error
	}

	uuid1 := uuid.New()

	testCases := []testCase{
		{
			name:   "nominal",
			poolID: uuid1.String(),
			atTime: "2021-01-01T00:00:00Z",
		},
		{
			name:          "invalid poolID",
			poolID:        "invalid",
			atTime:        "2021-01-01T00:00:00Z",
			expectedError: ErrValidation,
		},
		{
			name:          "invalid atTime",
			poolID:        uuid1.String(),
			atTime:        "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			expectedResponse := &GetPoolBalanceResponse{
				Balances: []*Balance{
					{
						Asset:  "EUR/2",
						Amount: big.NewInt(200),
					},
					{
						Asset:  "USD/2",
						Amount: big.NewInt(300),
					},
				},
			}

			balances, err := service.GetPoolBalance(context.Background(), tc.poolID, tc.atTime)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
				fmt.Println(balances)
				require.Equal(t, expectedResponse, balances)
			}
		})
	}
}
