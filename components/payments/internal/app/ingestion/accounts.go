package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/app/messages"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"

	"github.com/formancehq/payments/internal/app/models"
)

type AccountBatchElement struct {
	Account *models.Account
}

type AccountBatch []AccountBatchElement

type AccountIngesterFn func(ctx context.Context, batch AccountBatch, commitState any) error

func (fn AccountIngesterFn) IngestAccounts(ctx context.Context, batch AccountBatch, commitState any) error {
	return fn(ctx, batch, commitState)
}

func (i *DefaultIngester) IngestAccounts(ctx context.Context, batch AccountBatch) error {
	startingAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":       len(batch),
		"startingAt": startingAt,
	}).Debugf("Ingest accounts batch")

	accounts := make([]*models.Account, len(batch))
	for batchIdx := range batch {
		accounts[batchIdx] = batch[batchIdx].Account
	}

	if err := i.repo.UpsertAccounts(ctx, i.provider, accounts); err != nil {
		return fmt.Errorf("error upserting accounts: %w", err)
	}

	err := i.publisher.Publish(events.TopicPayments,
		publish.NewMessage(ctx, messages.NewEventSavedAccounts(accounts)))
	if err != nil {
		logging.FromContext(ctx).Errorf("Publishing message: %w", err)
	}

	endedAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":    len(batch),
		"endedAt": endedAt,
		"latency": endedAt.Sub(startingAt).String(),
	}).Debugf("Accounts batch ingested")

	return nil
}
