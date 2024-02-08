package mangopay

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/formancehq/payments/cmd/connectors/internal/connectors"
	"github.com/formancehq/payments/cmd/connectors/internal/connectors/mangopay/client"
	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/storage"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/payments/internal/otel"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
)

func taskCreateExternalBankAccount(mangopayClient *client.Client, bankAccountID uuid.UUID) task.Task {
	return func(
		ctx context.Context,
		connectorID models.ConnectorID,
		taskID models.TaskID,
		storageReader storage.Reader,
		ingester ingestion.Ingester,
	) error {
		ctx, span := connectors.StartSpan(
			ctx,
			"mangopay.taskCreateExternalBankAccount",
			attribute.String("connectorID", connectorID.String()),
			attribute.String("taskID", taskID.String()),
			attribute.String("bankAccountID", bankAccountID.String()),
		)
		defer span.End()

		bankAccount, err := storageReader.GetBankAccount(ctx, bankAccountID, true)
		if err != nil {
			otel.RecordError(span, err)
			return err
		}

		if err := createExternalBankAccount(ctx, connectorID, mangopayClient, bankAccount, ingester); err != nil {
			otel.RecordError(span, err)
			return err
		}

		return nil
	}
}

func createExternalBankAccount(
	ctx context.Context,
	connectorID models.ConnectorID,
	mangopayClient *client.Client,
	bankAccount *models.BankAccount,
	ingester ingestion.Ingester,
) error {
	userID := models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayUserIDMetadataKey)
	if userID == "" {
		return fmt.Errorf("missing userID in bank account metadata")
	}

	ownerAddress := client.OwnerAddress{
		AddressLine1: models.ExtractNamespacedMetadata(bankAccount.Metadata, models.BankAccountOwnerAddressLine1MetadataKey),
		AddressLine2: models.ExtractNamespacedMetadata(bankAccount.Metadata, models.BankAccountOwnerAddressLine2MetadataKey),
		City:         models.ExtractNamespacedMetadata(bankAccount.Metadata, models.BankAccountOwnerCityMetadataKey),
		Region:       models.ExtractNamespacedMetadata(bankAccount.Metadata, models.BankAccountOwnerRegionMetadataKey),
		PostalCode:   models.ExtractNamespacedMetadata(bankAccount.Metadata, models.BankAccountOwnerPostalCodeMetadataKey),
		Country:      bankAccount.Country,
	}

	var mangopayBankAccount *client.BankAccount
	if bankAccount.IBAN != "" {
		req := &client.CreateIBANBankAccountRequest{
			OwnerName:    bankAccount.Name,
			OwnerAddress: &ownerAddress,
			IBAN:         bankAccount.IBAN,
			BIC:          bankAccount.SwiftBicCode,
			Tag:          models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayTagMetadataKey),
		}

		var err error
		mangopayBankAccount, err = mangopayClient.CreateIBANBankAccount(ctx, userID, req)
		if err != nil {
			return err
		}
	} else {
		switch bankAccount.Country {
		case "US":
			req := &client.CreateUSBankAccountRequest{
				OwnerName:          bankAccount.Name,
				OwnerAddress:       &ownerAddress,
				AccountNumber:      bankAccount.AccountNumber,
				ABA:                models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayABAMetadataKey),
				DepositAccountType: models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayDepositAccountTypeMetadataKey),
				Tag:                models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = mangopayClient.CreateUSBankAccount(ctx, userID, req)
			if err != nil {
				return err
			}

		case "CA":
			req := &client.CreateCABankAccountRequest{
				OwnerName:         bankAccount.Name,
				OwnerAddress:      &ownerAddress,
				AccountNumber:     bankAccount.AccountNumber,
				InstitutionNumber: models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayInstitutionNumberMetadataKey),
				BranchCode:        models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayBranchCodeMetadataKey),
				BankName:          models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayBankNameMetadataKey),
				Tag:               models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = mangopayClient.CreateCABankAccount(ctx, userID, req)
			if err != nil {
				return err
			}

		case "GB":
			req := &client.CreateGBBankAccountRequest{
				OwnerName:     bankAccount.Name,
				OwnerAddress:  &ownerAddress,
				AccountNumber: bankAccount.AccountNumber,
				SortCode:      models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopaySortCodeMetadataKey),
				Tag:           models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = mangopayClient.CreateGBBankAccount(ctx, userID, req)
			if err != nil {
				return err
			}

		default:
			req := &client.CreateOtherBankAccountRequest{
				OwnerName:     bankAccount.Name,
				OwnerAddress:  &ownerAddress,
				AccountNumber: bankAccount.AccountNumber,
				BIC:           bankAccount.SwiftBicCode,
				Country:       bankAccount.Country,
				Tag:           models.ExtractNamespacedMetadata(bankAccount.Metadata, client.MangopayTagMetadataKey),
			}

			var err error
			mangopayBankAccount, err = mangopayClient.CreateOtherBankAccount(ctx, userID, req)
			if err != nil {
				return err
			}
		}
	}

	if mangopayBankAccount != nil {
		buf, err := json.Marshal(mangopayBankAccount)
		if err != nil {
			return err
		}

		externalAccount := &models.Account{
			ID: models.AccountID{
				Reference:   mangopayBankAccount.ID,
				ConnectorID: connectorID,
			},
			CreatedAt:   time.Unix(mangopayBankAccount.CreationDate, 0),
			Reference:   mangopayBankAccount.ID,
			ConnectorID: connectorID,
			AccountName: mangopayBankAccount.OwnerName,
			Type:        models.AccountTypeExternal,
			Metadata: map[string]string{
				"user_id": userID,
			},
			RawData: buf,
		}

		if err := ingester.IngestAccounts(ctx, ingestion.AccountBatch{externalAccount}); err != nil {
			return err
		}

		if err := ingester.LinkBankAccountWithAccount(ctx, bankAccount, &externalAccount.ID); err != nil {
			return err
		}
	}

	return nil
}
