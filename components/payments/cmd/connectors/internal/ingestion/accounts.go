package ingestion

import (
	"context"
	"fmt"
	"time"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
)

type AccountBatch []*models.Account

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

	idsInserted, err := i.store.UpsertAccounts(ctx, batch)
	if err != nil {
		return fmt.Errorf("error upserting accounts: %w", err)
	}

	idsInsertedMap := make(map[string]struct{}, len(idsInserted))
	for idx := range idsInserted {
		idsInsertedMap[idsInserted[idx].String()] = struct{}{}
	}

	for accountIdx := range batch {
		_, ok := idsInsertedMap[batch[accountIdx].ID.String()]
		if !ok {
			// No need to publish an event for an already existing payment
			continue
		}

		if err := i.publisher.Publish(
			events.TopicPayments,
			publish.NewMessage(
				ctx,
				i.messages.NewEventSavedAccounts(i.provider, batch[accountIdx]),
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
