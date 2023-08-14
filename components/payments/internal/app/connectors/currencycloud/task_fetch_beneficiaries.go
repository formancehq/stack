package currencycloud

import (
	"context"
	"encoding/json"
	"time"

	"github.com/formancehq/payments/internal/app/connectors/currencycloud/client"
	"github.com/formancehq/payments/internal/app/ingestion"
	"github.com/formancehq/payments/internal/app/metrics"
	"github.com/formancehq/payments/internal/app/models"
	"github.com/formancehq/payments/internal/app/task"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	beneficiariesAttrs = metric.WithAttributes(append(connectorAttrs, attribute.String(metrics.ObjectAttributeKey, "beneficiaries"))...)
)

func taskFetchBeneficiaries(
	logger logging.Logger,
	client *client.Client,
) task.Task {
	return func(
		ctx context.Context,
		ingester ingestion.Ingester,
		scheduler task.Scheduler,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		logger.Info(taskFetchBeneficiaries)

		now := time.Now()
		defer func() {
			metricsRegistry.ConnectorObjectsLatency().Record(ctx, time.Since(now).Milliseconds(), beneficiariesAttrs)
		}()

		page := 1
		for {
			if page < 0 {
				break
			}

			pagedBeneficiaries, nextPage, err := client.GetBeneficiaries(ctx, page)
			if err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, beneficiariesAttrs)
				return err
			}

			page = nextPage

			if err := ingestBeneficiariesAccountsBatch(ctx, ingester, pagedBeneficiaries); err != nil {
				metricsRegistry.ConnectorObjectsErrors().Add(ctx, 1, beneficiariesAttrs)
				return err
			}
			metricsRegistry.ConnectorObjects().Add(ctx, int64(len(pagedBeneficiaries)), beneficiariesAttrs)
		}

		return nil
	}
}

func ingestBeneficiariesAccountsBatch(
	ctx context.Context,
	ingester ingestion.Ingester,
	beneficiaries []*client.Beneficiary,
) error {
	batch := ingestion.AccountBatch{}
	for _, beneficiary := range beneficiaries {
		raw, err := json.Marshal(beneficiary)
		if err != nil {
			return err
		}

		batch = append(batch, &models.Account{
			ID: models.AccountID{
				Reference: beneficiary.ID,
				Provider:  models.ConnectorProviderCurrencyCloud,
			},
			// Moneycorp does not send the opening date of the account
			CreatedAt:   beneficiary.CreatedAt,
			Reference:   beneficiary.ID,
			Provider:    models.ConnectorProviderCurrencyCloud,
			AccountName: beneficiary.Name,
			Type:        models.AccountTypeExternal,
			RawData:     raw,
		})
	}

	if err := ingester.IngestAccounts(ctx, batch); err != nil {
		return err
	}

	return nil
}
