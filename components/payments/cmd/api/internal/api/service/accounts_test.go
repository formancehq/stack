package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCreateAccout(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                 string
		request              *CreateAccountRequest
		isConnectorInstalled bool
		expectedError        error
	}

	connectorID := models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	testCases := []testCase{
		{
			name: "nominal",
			request: &CreateAccountRequest{
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
			isConnectorInstalled: true,
		},
		{
			name: "nominal without default asset",
			request: &CreateAccountRequest{
				Reference:   "test",
				ConnectorID: connectorID.String(),
				CreatedAt:   time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				AccountName: "test",
				Type:        "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			isConnectorInstalled: true,
		},
		{
			name: "connector not installed",
			request: &CreateAccountRequest{
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
			isConnectorInstalled: false,
			expectedError:        ErrValidation,
		},
		{
			name: "invalid connectorID",
			request: &CreateAccountRequest{
				Reference:    "test",
				ConnectorID:  "invalid",
				CreatedAt:    time.Date(2023, 11, 22, 8, 0, 0, 0, time.UTC),
				DefaultAsset: "USD/2",
				AccountName:  "test",
				Type:         "INTERNAL",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			expectedError:        ErrValidation,
			isConnectorInstalled: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			store := &MockStore{}
			service := New(store.WithIsConnectorInstalled(tc.isConnectorInstalled), &MockPublisher{}, messages.NewMessages(""))
			p, err := service.CreateAccount(context.Background(), tc.request)
			if tc.expectedError != nil {
				require.True(t, errors.Is(err, tc.expectedError))
			} else {
				require.NoError(t, err)
				require.NotNil(t, p)
			}
		})
	}
}

func TestGetAccount(t *testing.T) {
	t.Parallel()

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

	service := New(&MockStore{}, &MockPublisher{}, messages.NewMessages(""))

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
