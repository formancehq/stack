package service

import (
	"context"
	"errors"
	"testing"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestGetAccount(t *testing.T) {
	type testCase struct {
		name          string
		accountID     string
		expectedError error
	}

	accountID := models.AccountID{
		Reference: "a1",
		ConnectorID: models.ConnectorID{
			Reference: uuid.New(),
			Provider:  models.ConnectorProviderDummyPay,
		},
	}

	testCases := []testCase{
		{
			name:          "nominal",
			accountID:     accountID.String(),
			expectedError: nil,
		},
		{
			name:          "invalid accountID",
			accountID:     "invalid",
			expectedError: ErrValidation,
		},
	}

	service := New(&MockStore{})

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := service.GetAccount(context.Background(), tc.accountID)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
			}
		})
	}
}
