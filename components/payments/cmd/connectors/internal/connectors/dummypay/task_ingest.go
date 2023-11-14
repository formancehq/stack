package dummypay

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const taskKeyIngest = "ingest"

var (
	paymentsAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "payments"))...)
)

// newTaskIngest returns a new task descriptor for the ingest task.
func newTaskIngest(filePath string) TaskDescriptor {
	return TaskDescriptor{
		Name:     "Ingest payments from read files",
		Key:      taskKeyIngest,
		FileName: filePath,
	}
}

// taskIngest ingests a payment file.
func taskIngest(config Config, descriptor TaskDescriptor, fs fs) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), paymentsAttrs)
		}()

		ingestionPayload, err := parseIngestionPayload(connectorID, config, descriptor, fs)
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return err
		}

		// Ingest the payment into the system.
		err = ingester.IngestPayments(ctx, connectorID, ingestionPayload, struct{}{})
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, paymentsAttrs)
			return fmt.Errorf("failed to ingest file '%s': %w", descriptor.FileName, err)
		}
		metricsRegistry.ConnectorObjects().Add(ctx, 1, paymentsAttrs)

		return nil
	}
}

func parseIngestionPayload(connectorID models.ConnectorID, config Config, descriptor TaskDescriptor, fs fs) (ingestion.PaymentBatch, error) {
	// Open the file.
	file, err := fs.Open(filepath.Join(config.Directory, descriptor.FileName))
	if err != nil {
		return nil, fmt.Errorf("failed to open file '%s': %w", descriptor.FileName, err)
	}

	defer file.Close()

	var paymentElement payment

	// Decode the JSON file.
	err = json.NewDecoder(file).Decode(&paymentElement)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file '%s': %w", descriptor.FileName, err)
	}

	ingestionPayload := ingestion.PaymentBatch{ingestion.PaymentBatchElement{
		Payment: &models.Payment{
			ID: models.PaymentID{
				PaymentReference: models.PaymentReference{
					Reference: paymentElement.Reference,
					Type:      paymentElement.Type,
				},
				ConnectorID: connectorID,
			},
			Reference:   paymentElement.Reference,
			ConnectorID: connectorID,
			Amount:      paymentElement.Amount,
			Type:        paymentElement.Type,
			Status:      paymentElement.Status,
			Scheme:      paymentElement.Scheme,
			Asset:       paymentElement.Asset,
			RawData:     paymentElement.RawData,
		},
		Update: true,
	}}

	return ingestionPayload, nil
}
