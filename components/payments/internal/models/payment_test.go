package models

import (
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/go-libs/pointer"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestFromPSPPaymentToPayment(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	connectorID := ConnectorID{
		Reference: uuid.New(),
		Provider:  "test",
	}

	t.Run("parent reference is empty", func(t *testing.T) {
		t.Parallel()

		pspPayment := PSPPayment{
			ParentReference:        "",
			Reference:              "test1",
			CreatedAt:              now.UTC(),
			Type:                   PAYMENT_TYPE_PAYOUT,
			Amount:                 big.NewInt(100),
			Asset:                  "USD/2",
			Scheme:                 PAYMENT_SCHEME_OTHER,
			Status:                 PAYMENT_STATUS_CANCELLED,
			SourceAccountReference: pointer.For("acc"),
			Metadata: map[string]string{
				"foo": "bar",
			},
			Raw: []byte(`{}`),
		}

		actual := FromPSPPaymentToPayment(pspPayment, connectorID)

		pid := PaymentID{
			PaymentReference: PaymentReference{
				Reference: "test1",
				Type:      PAYMENT_TYPE_PAYOUT,
			},
			ConnectorID: connectorID,
		}
		expected := Payment{
			ID:            pid,
			ConnectorID:   connectorID,
			Reference:     "test1",
			CreatedAt:     now.UTC(),
			Type:          PAYMENT_TYPE_PAYOUT,
			InitialAmount: big.NewInt(100),
			Amount:        big.NewInt(100),
			Asset:         "USD/2",
			Scheme:        PAYMENT_SCHEME_OTHER,
			Status:        PAYMENT_STATUS_CANCELLED,
			SourceAccountID: &AccountID{
				Reference:   "acc",
				ConnectorID: connectorID,
			},
			Metadata: map[string]string{
				"foo": "bar",
			},
			Adjustments: []PaymentAdjustment{
				{
					ID: PaymentAdjustmentID{
						PaymentID: pid,
						Reference: "test1",
						CreatedAt: now.UTC(),
						Status:    PAYMENT_STATUS_CANCELLED,
					},
					PaymentID: pid,
					Reference: "test1",
					CreatedAt: now.UTC(),
					Status:    PAYMENT_STATUS_CANCELLED,
					Amount:    big.NewInt(100),
					Asset:     pointer.For("USD/2"),
					Metadata: map[string]string{
						"foo": "bar",
					},
					Raw: []byte(`{}`),
				},
			},
		}

		comparePayment(t, expected, actual)
	})

	t.Run("parent reference is not empty", func(t *testing.T) {
		t.Parallel()

		pspPayment := PSPPayment{
			ParentReference:             "parent_reference",
			Reference:                   "test1",
			CreatedAt:                   now.UTC(),
			Type:                        PAYMENT_TYPE_TRANSFER,
			Amount:                      big.NewInt(150),
			Asset:                       "EUR/2",
			Scheme:                      PAYMENT_SCHEME_OTHER,
			Status:                      PAYMENT_STATUS_SUCCEEDED,
			DestinationAccountReference: pointer.For("acc"),
			Metadata: map[string]string{
				"foo": "bar",
			},
			Raw: []byte(`{}`),
		}

		actual := FromPSPPaymentToPayment(pspPayment, connectorID)

		pid := PaymentID{
			PaymentReference: PaymentReference{
				Reference: "parent_reference",
				Type:      PAYMENT_TYPE_TRANSFER,
			},
			ConnectorID: connectorID,
		}
		expected := Payment{
			ID:            pid,
			ConnectorID:   connectorID,
			Reference:     "parent_reference",
			CreatedAt:     now.UTC(),
			Type:          PAYMENT_TYPE_TRANSFER,
			InitialAmount: big.NewInt(150),
			Amount:        big.NewInt(150),
			Asset:         "EUR/2",
			Scheme:        PAYMENT_SCHEME_OTHER,
			Status:        PAYMENT_STATUS_SUCCEEDED,
			DestinationAccountID: &AccountID{
				Reference:   "acc",
				ConnectorID: connectorID,
			},
			Metadata: map[string]string{
				"foo": "bar",
			},
			Adjustments: []PaymentAdjustment{
				{
					ID: PaymentAdjustmentID{
						PaymentID: pid,
						Reference: "test1",
						CreatedAt: now.UTC(),
						Status:    PAYMENT_STATUS_SUCCEEDED,
					},
					PaymentID: pid,
					Reference: "test1",
					CreatedAt: now.UTC(),
					Status:    PAYMENT_STATUS_SUCCEEDED,
					Amount:    big.NewInt(150),
					Asset:     pointer.For("EUR/2"),
					Metadata: map[string]string{
						"foo": "bar",
					},
					Raw: []byte(`{}`),
				},
			},
		}

		comparePayment(t, expected, actual)
	})
}

func comparePayment(t *testing.T, expected, actual Payment) {
	require.Equal(t, expected.ID, actual.ID)
	require.Equal(t, expected.ConnectorID, actual.ConnectorID)
	require.Equal(t, expected.Reference, actual.Reference)
	require.Equal(t, expected.CreatedAt, actual.CreatedAt)
	require.Equal(t, expected.Type, actual.Type)
	require.Equal(t, expected.InitialAmount, actual.InitialAmount)
	require.Equal(t, expected.Amount, actual.Amount)
	require.Equal(t, expected.Asset, actual.Asset)
	require.Equal(t, expected.Scheme, actual.Scheme)
	require.Equal(t, expected.Status, actual.Status)

	switch {
	case expected.SourceAccountID == nil && actual.SourceAccountID == nil:
	case expected.SourceAccountID != nil && actual.SourceAccountID != nil:
		require.Equal(t, *expected.SourceAccountID, *actual.SourceAccountID)
	default:
		require.Fail(t, "source account id mismatch")
	}

	switch {
	case expected.DestinationAccountID == nil && actual.DestinationAccountID == nil:
	case expected.DestinationAccountID != nil && actual.DestinationAccountID != nil:
		require.Equal(t, *expected.DestinationAccountID, *actual.DestinationAccountID)
	default:
		require.Fail(t, "destination account id mismatch")
	}

	require.Equal(t, len(expected.Metadata), len(actual.Metadata))
	for k, v := range expected.Metadata {
		_, ok := actual.Metadata[k]
		require.True(t, ok)
		require.Equal(t, v, actual.Metadata[k])
	}

	compareAdjustments(t, expected.Adjustments, actual.Adjustments)
}

func compareAdjustments(t *testing.T, expected, actual []PaymentAdjustment) {
	require.Equal(t, len(expected), len(actual))
	for i := range expected {
		require.Equal(t, expected[i].ID, actual[i].ID)
		require.Equal(t, expected[i].PaymentID, actual[i].PaymentID)
		require.Equal(t, expected[i].Reference, actual[i].Reference)
		require.Equal(t, expected[i].CreatedAt, actual[i].CreatedAt)
		require.Equal(t, expected[i].Status, actual[i].Status)
		require.Equal(t, expected[i].Amount, actual[i].Amount)
		require.Equal(t, expected[i].Asset, actual[i].Asset)

		require.Equal(t, len(expected[i].Metadata), len(actual[i].Metadata))
		for k, v := range expected[i].Metadata {
			_, ok := actual[i].Metadata[k]
			require.True(t, ok)
			require.Equal(t, v, actual[i].Metadata[k])
		}
	}
}
