package ingestion

import (
	"context"
	"encoding/json"
	"math/big"
	"testing"
	"time"

	"github.com/formancehq/payments/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	connectorID = models.ConnectorID{
		Reference: uuid.New(),
		Provider:  models.ConnectorProviderDummyPay,
	}

	p1 = &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: "p1",
				Type:      models.PaymentTypePayIn,
			},
			ConnectorID: connectorID,
		},
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 4, 55, 0, 0, time.UTC),
		Reference:   "p1",
		Amount:      big.NewInt(100),
		Type:        models.PaymentTypePayIn,
		Status:      models.PaymentStatusCancelled,
		Scheme:      models.PaymentSchemeA2A,
		Asset:       models.Asset("USD/2"),
	}

	p2 = &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: "p2",
				Type:      models.PaymentTypeTransfer,
			},
			ConnectorID: connectorID,
		},
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 4, 54, 0, 0, time.UTC),
		Reference:   "p2",
		Amount:      big.NewInt(150),
		Type:        models.PaymentTypeTransfer,
		Status:      models.PaymentStatusSucceeded,
		Scheme:      models.PaymentSchemeApplePay,
		Asset:       models.Asset("EUR/2"),
	}

	p3 = &models.Payment{
		ID: models.PaymentID{
			PaymentReference: models.PaymentReference{
				Reference: "p3",
				Type:      models.PaymentTypePayOut,
			},
			ConnectorID: connectorID,
		},
		ConnectorID: connectorID,
		CreatedAt:   time.Date(2023, 11, 14, 4, 53, 0, 0, time.UTC),
		Reference:   "p3",
		Amount:      big.NewInt(200),
		Type:        models.PaymentTypePayOut,
		Status:      models.PaymentStatusPending,
		Scheme:      models.PaymentSchemeCardMasterCard,
		Asset:       models.Asset("USD/2"),
	}
)

type paymentMessagePayload struct {
	Paylaod struct {
		ID string `json:"id"`
	} `json:"payload"`
}

func TestIngestPayments(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name                        string
		batch                       PaymentBatch
		paymentIDsNotModified       []models.PaymentID
		requiredPublishedPaymentIDs []models.PaymentID
	}

	testCases := []testCase{
		{
			name: "nominal",
			batch: PaymentBatch{
				{
					Payment: p1,
				},
				{
					Payment: p2,
				},
				{
					Payment: p3,
				},
			},
			paymentIDsNotModified:       []models.PaymentID{},
			requiredPublishedPaymentIDs: []models.PaymentID{p1.ID, p2.ID, p3.ID},
		},
		{
			name: "only one payment upserted, should publish only one message",
			batch: PaymentBatch{
				{
					Payment: p1,
				},
				{
					Payment: p2,
				},
				{
					Payment: p3,
				},
			},
			paymentIDsNotModified:       []models.PaymentID{p1.ID, p2.ID},
			requiredPublishedPaymentIDs: []models.PaymentID{p3.ID},
		},
		{
			name: "all payments are not modified, should not publish any message",
			batch: PaymentBatch{
				{
					Payment: p1,
				},
				{
					Payment: p2,
				},
				{
					Payment: p3,
				},
			},
			paymentIDsNotModified:       []models.PaymentID{p1.ID, p2.ID, p3.ID},
			requiredPublishedPaymentIDs: []models.PaymentID{},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			publisher := NewMockPublisher()

			ingester := NewDefaultIngester(
				models.ConnectorProviderDummyPay,
				nil,
				NewMockStore().WithPaymentIDsNotModified(tc.paymentIDsNotModified),
				publisher,
			)

			err := ingester.IngestPayments(context.Background(), connectorID, tc.batch, nil)
			publisher.Close()
			require.NoError(t, err)

			require.Len(t, publisher.messages, len(tc.requiredPublishedPaymentIDs))
			i := 0
			for msg := range publisher.messages {
				var payload paymentMessagePayload
				require.NoError(t, json.Unmarshal(msg.Payload, &payload))
				require.Equal(t, tc.requiredPublishedPaymentIDs[i].String(), payload.Paylaod.ID)
				i++
			}
		})
	}
}
