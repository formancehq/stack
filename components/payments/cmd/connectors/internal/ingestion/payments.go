package ingestion

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/internal/messages"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/pkg/events"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/publish"
)

type PaymentBatchElement struct {
	Payment    *models.Payment
	Adjustment *models.Adjustment
	Metadata   *models.Metadata
	Update     bool
}

type PaymentBatch []PaymentBatchElement

type IngesterFn func(ctx context.Context, batch PaymentBatch, commitState any) error

func (fn IngesterFn) IngestPayments(ctx context.Context, batch PaymentBatch, commitState any) error {
	return fn(ctx, batch, commitState)
}

func (i *DefaultIngester) IngestPayments(
	ctx context.Context,
	connectorID models.ConnectorID,
	batch PaymentBatch,
	commitState any,
) error {
	startingAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":       len(batch),
		"startingAt": startingAt,
	}).Debugf("Ingest batch")

	var allPayments []*models.Payment //nolint:prealloc // length is unknown

	for batchIdx := range batch {
		payment := batch[batchIdx].Payment

		if payment == nil {
			continue
		}

		allPayments = append(allPayments, payment)
	}

	idsInserted, err := i.store.UpsertPayments(ctx, allPayments)
	if err != nil {
		return fmt.Errorf("error upserting payments: %w", err)
	}

	idsInsertedMap := make(map[string]struct{}, len(idsInserted))
	for idx := range idsInserted {
		idsInsertedMap[idsInserted[idx].String()] = struct{}{}
	}

	taskState, err := json.Marshal(commitState)
	if err != nil {
		return fmt.Errorf("error marshaling task state: %w", err)
	}

	if err = i.store.UpdateTaskState(ctx, connectorID, i.descriptor, taskState); err != nil {
		return fmt.Errorf("error updating task state: %w", err)
	}

	for paymentIdx := range allPayments {
		_, ok := idsInsertedMap[allPayments[paymentIdx].ID.String()]
		if !ok {
			// No need to publish an event for an already existing payment
			continue
		}
		err = i.publisher.Publish(events.TopicPayments,
			publish.NewMessage(ctx, messages.NewEventSavedPayments(i.provider, allPayments[paymentIdx])))
		if err != nil {
			logging.FromContext(ctx).Errorf("Publishing message: %w", err)

			continue
		}
	}

	endedAt := time.Now()

	logging.FromContext(ctx).WithFields(map[string]interface{}{
		"size":    len(batch),
		"endedAt": endedAt,
		"latency": endedAt.Sub(startingAt).String(),
	}).Debugf("Batch ingested")

	return nil
}
