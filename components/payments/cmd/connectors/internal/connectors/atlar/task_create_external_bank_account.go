package atlar

import (
	"context"
	"errors"
	"fmt"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/atlar/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"go.opentelemetry.io/otel/attribute"
)

func CreateExternalBankAccountTask(config Config, client *client.Client, newExternalBankAccount *models.BankAccount) task.Task {
	return func(
		ctx context.Context,
		taskID models.TaskID,
		connectorID models.ConnectorID,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"atlar.taskCreateExternalBankAccount",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("bankAccount.name", newExternalBankAccount.Name),
			attribute.String("bankAccount.id", newExternalBankAccount.ID.String()),
		)
		defer span.End()

		err := validateExternalBankAccount(newExternalBankAccount)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		externalAccountID, err := createExternalBankAccount(ctx, client, newExternalBankAccount)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}
		if externalAccountID == nil {
			err := errors.New("no external account id returned")
			otel.RecordError(span, err)
			return err
		}

		err = ingestExternalAccountFromAtlar(
			ctx,
			connectorID,
			ingester,
			client,
			newExternalBankAccount,
			*externalAccountID,
		)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

// TODO: validation (also metadata) needs to return a 400
func validateExternalBankAccount(newExternalBankAccount *models.BankAccount) error {
	_, err := ExtractNamespacedMetadata(newExternalBankAccount.Metadata, "owner/name")
	if err != nil {
		return fmt.Errorf("required metadata field %sowner/name is missing", atlarMetadataSpecNamespace)
	}
	ownerType, err := ExtractNamespacedMetadata(newExternalBankAccount.Metadata, "owner/type")
	if err != nil {
		return fmt.Errorf("required metadata field %sowner/type is missing", atlarMetadataSpecNamespace)
	}
	if *ownerType != "INDIVIDUAL" && *ownerType != "COMPANY" {
		return fmt.Errorf("metadata field %sowner/type needs to be one of [ INDIVIDUAL COMPANY ]", atlarMetadataSpecNamespace)
	}

	return nil
}

func createExternalBankAccount(ctx context.Context, client *client.Client, newExternalBankAccount *models.BankAccount) (*string, error) {
	return client.CreateCounterParty(ctx, newExternalBankAccount)
}

func ingestExternalAccountFromAtlar(
	ctx context.Context,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	client *client.Client,
	formanceBankAccount *models.BankAccount,
	externalAccountID string,
) error {
	accountsBatch := ingestion.AccountBatch{}

	externalAccountResponse, err := client.GetV1ExternalAccountsID(ctx, externalAccountID)
	if err != nil {
		return err
	}

	counterpartyResponse, err := client.GetV1CounterpartiesID(ctx, externalAccountResponse.Payload.CounterpartyID)
	if err != nil {
		return err
	}

	newAccount, err := ExternalAccountFromAtlarData(connectorID, externalAccountResponse.Payload, counterpartyResponse.Payload)
	if err != nil {
		return err
	}

	accountsBatch = append(accountsBatch, newAccount)

	err = ingester.IngestAccounts(ctx, accountsBatch)
	if err != nil {
		return err
	}

	if err := ingester.LinkBankAccountWithAccount(ctx, formanceBankAccount, &newAccount.ID); err != nil {
		return err
	}

	return nil
}
