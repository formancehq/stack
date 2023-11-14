package modulr

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/internal/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors/modulr/client"
	"github.com/formancehq/payments/cmd/connectors/internal/task"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

var (
	beneficiariesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "beneficiaries"))...)
)

func taskFetchBeneficiaries(logger logging.Logger, client *client.Client) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskNameFetchBeneficiaries)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), beneficiariesAttrs)
		}()

		beneficiaries, err := client.GetBeneficiaries()
		if err != nil {
			metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, beneficiariesAttrs)
			return err
		}

		if err := ingestBeneficiariesAccountsBatch(ctx, connectorID, ingester, metricsRegistry, beneficiaries); err != nil {
			return err
		}

		return nil
	}
}

func ingestBeneficiariesAccountsBatch(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	metricsRegistry metrics.MetricsRegistry,
	beneficiaries []*client.Beneficiary,
) error {
	accountsBatch := ingestion.AccountBatch{}

	for _, beneficiary := range beneficiaries {
		raw, err := json.Marshal(beneficiary)
		if err != nil {
			return err
		}

		openingDate, err := time.Parse("2006-01-02T15:04:05.999999999+0000", beneficiary.Created)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, &models.Account{
			ID: models.AccountID{
				Reference:   beneficiary.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:   openingDate,
			Reference:   beneficiary.ID,
			ConnectorID: connectorID,
			AccountName: beneficiary.Name,
			Type:        models.AccountTypeExternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
		metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, beneficiariesAttrs)
		return err
	}
	metricsRegistry.ConnectorObjects().Add(ctx, int64(len(accountsBatch)), beneficiariesAttrs)

	return nil
}
