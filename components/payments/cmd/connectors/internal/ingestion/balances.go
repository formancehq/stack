package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

type BalanceBatch []*models.Balance

type BalanceIngesterFn func(ctx context.Context, batch BalanceBatch) error

func (fn BalanceIngesterFn) IngestBalances(ctx context.Context, batch BalanceBatch) error {
	return fn(ctx, batch)
}

func (i *DefaultIngester) IngestBalances(ctx context.Context, batch BalanceBatch, checkIfAccountExists bool) error {
	startingAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":       len(batch),
		"startingAt": startingAt,
	}).Debugf("Ingest balances batch")

	if err := i.repo.InsertBalances(ctx, batch, checkIfAccountExists); err != nil {
		return fmt.Errorf("error inserting balances: %w", err)
	}

	for _, balance := range batch {
		if err := i.publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				ctx,
				messages.NewEventSavedBalances(balance, i.provider),
			),
		); err != nil {
			logging.FromContext(ctx).Errorf("Publishing message: %w", err)
		}
	}

	endedAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":    len(batch),
		"endedAt": endedAt,
		"latency": endedAt.Sub(startingAt).String(),
	}).Debugf("Accounts batch ingested")

	return nil
}
