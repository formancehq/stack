package atlar

import (
	"context"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	"github.com/get-momo/atlar-v1-go-client/client/counterparties"
	"github.com/get-momo/atlar-v1-go-client/client/external_accounts"
)

func FetchExternalAccountTask(config Config, client *atlar_client.Rest, externalAccountID string) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		accountsBatch := ingestion.AccountBatch{}

		getExternalAccountParams := external_accounts.GetV1ExternalAccountsIDParams{
			Context: ctx,
			ID:      externalAccountID,
		}
		externalAccountResponse, err := client.ExternalAccounts.GetV1ExternalAccountsID(&getExternalAccountParams)
		if err != nil {
			return err
		}

		getCounterpartyParams := counterparties.GetV1CounterpartiesIDParams{
			Context: ctx,
			ID:      externalAccountResponse.Payload.CounterpartyID,
		}
		counterpartyResponse, err := client.Counterparties.GetV1CounterpartiesID(&getCounterpartyParams)
		if err != nil {
			return err
		}

		newAccount, err := ExternalAccountFromAtlarData(connectorID, externalAccountResponse.Payload, counterpartyResponse.Payload)
		if err != nil {
			return err
		}

		accountsBatch = append(accountsBatch, newAccount)

		if err := ingester.IngestAccounts(ctx, accountsBatch); err != nil {
			return err
		}

		return nil
	}
}
