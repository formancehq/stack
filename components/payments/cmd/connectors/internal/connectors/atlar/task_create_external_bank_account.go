package atlar

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/formancehq/payments/cmd/connectors/internal/ingestion"
	"github.com/formancehq/payments/cmd/connectors/internal/metrics"
	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/logging"
	atlar_client "github.com/get-momo/atlar-v1-go-client/client"
	"github.com/get-momo/atlar-v1-go-client/client/counterparties"
	"github.com/get-momo/atlar-v1-go-client/client/external_accounts"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
)

func CreateExternalBankAccountTask(config Config, client *atlar_client.Rest, newExternalBankAccount *models.BankAccount) task.Task {
	return func(
		ctx context.Context,
		logger logging.Logger,
		connectorID models.ConnectorID,
		resolver task.StateResolver,
		scheduler task.Scheduler,
		ingester ingestion.Ingester,
		metricsRegistry metrics.MetricsRegistry,
	) error {
		err := validateExternalBankAccount(newExternalBankAccount)
		if err != nil {
			return err
		}

		externalAccountID, err := createExternalBankAccount(ctx, logger, config, newExternalBankAccount)
		if err != nil {
			return err
		}
		if externalAccountID == nil {
			return errors.New("no external account id returned")
		}

		err = ingestExternalAccountFromAtlar(
			ctx,
			logger,
			connectorID,
			ingester,
			client,
			*externalAccountID,
		)
		if err != nil {
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

func createExternalBankAccount(ctx context.Context, logger logging.Logger, config Config, newExternalBankAccount *models.BankAccount) (*string, error) {
	client := createAtlarClient(&config)

	// TODO: make sure an account with that IBAN does not already exist (Atlar API v2 needed, v1 lacks the filters)
	// alternatively we could query the local DB

	createCounterpartyRequest := atlar_models.CreateCounterpartyRequest{
		Name:      ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/name"),
		PartyType: *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/type"),
		ContactDetails: &atlar_models.ContactDetails{
			Email: *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/email"),
			Phone: *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/phone"),
			Address: &atlar_models.Address{
				StreetName:   *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/address/streetName"),
				StreetNumber: *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/address/streetNumber"),
				City:         *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/address/city"),
				PostalCode:   *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/address/postalCode"),
				Country:      *ExtractNamespacedMetadataIgnoreEmpty(newExternalBankAccount.Metadata, "owner/contact/address/country"),
			},
		},
		ExternalAccounts: []*atlar_models.CreateEmbeddedExternalAccountRequest{
			{
				// ExternalID could cause problems when synchronizing with Accounts[type=external]
				Bank: &atlar_models.UpdatableBank{
					Bic: newExternalBankAccount.SwiftBicCode,
				},
				Identifiers: extractAtlarAccountIdentifiersFromBankAccount(newExternalBankAccount),
			},
		},
	}
	postCounterpartiesParams := counterparties.PostV1CounterpartiesParams{
		Context:      ctx,
		Counterparty: &createCounterpartyRequest,
	}
	postCounterpartiesResponse, err := client.Counterparties.PostV1Counterparties(&postCounterpartiesParams)
	if err != nil {
		return nil, err
	}
	logger.WithContext(ctx).Debug("external bank account has been created")

	if len(postCounterpartiesResponse.Payload.ExternalAccounts) != 1 {
		// should never occur, but when in case it happens it's nice to have an error to search for
		return nil, errors.New("counterparty was not created with exactly one account")
	}

	externalAccountID := postCounterpartiesResponse.Payload.ExternalAccounts[0].ID

	return &externalAccountID, nil
}

func ingestExternalAccountFromAtlar(
	ctx context.Context,
	logger logging.Logger,
	connectorID models.ConnectorID,
	ingester ingestion.Ingester,
	client *atlar_client.Rest,
	externalAccountID string,
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
	logger.WithContext(ctx).Info("Got external Account from Atlar", newAccount)

	accountsBatch = append(accountsBatch, newAccount)

	err = ingester.IngestAccounts(ctx, accountsBatch)
	if err != nil {
		return err
	}

	return nil
}

func extractAtlarAccountIdentifiersFromBankAccount(bankAccount *models.BankAccount) []*atlar_models.AccountIdentifier {
	ownerName := bankAccount.Metadata[atlarMetadataSpecNamespace+"owner/name"]
	ibanType := "IBAN"
	accountIdentifiers := []*atlar_models.AccountIdentifier{{
		HolderName: &ownerName,
		Market:     &bankAccount.Country,
		Type:       &ibanType,
		Number:     &bankAccount.IBAN,
	}}
	for k := range bankAccount.Metadata {
		// check whether the key has format com.atlar.spec/identifier/<market>/<type>
		identifierData, err := metadataToIdentifierData(k, bankAccount.Metadata[k])
		if err != nil {
			// matadata does not describe an identifier
			continue
		}
		if identifierData.Market == bankAccount.Country && identifierData.Type == "IBAN" {
			// avoid duplicate identifiers
			continue
		}
		accountIdentifiers = append(accountIdentifiers, &atlar_models.AccountIdentifier{
			HolderName: &ownerName,
			Market:     &identifierData.Market,
			Type:       &identifierData.Type,
			Number:     &identifierData.Number,
		})
	}
	return accountIdentifiers
}

type IdentifierData struct {
	Market string
	Type   string
	Number string
}

var identifierMetadataRegex = regexp.MustCompile(`^com\.atlar\.spec/identifier/([^/]+)/([^/]+)$`)

func metadataToIdentifierData(key, value string) (*IdentifierData, error) {
	// Find matches in the input string
	matches := identifierMetadataRegex.FindStringSubmatch(key)
	if matches == nil {
		return nil, errors.New("input does not match the expected format")
	}

	// Extract values from the matched groups
	return &IdentifierData{
		Market: matches[1],
		Type:   matches[2],
		Number: value,
	}, nil
}
