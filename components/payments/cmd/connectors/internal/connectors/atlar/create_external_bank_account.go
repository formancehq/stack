package atlar

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/formancehq/payments/cmd/connectors/internal/task"
	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/contextutil"
	"github.com/get-momo/atlar-v1-go-client/client/counterparties"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
)

func createExternalBankAccount(ctx task.ConnectorContext, newExternalBankAccount *models.BankAccount, config Config) error {
	err := validateExternalBankAccount(newExternalBankAccount)
	if err != nil {
		return err
	}

	client := createAtlarClient(&config)
	detachedCtx, _ := contextutil.Detached(ctx.Context())
	// TODO: make sure an account with that IBAN does not already exist

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
		Context:      detachedCtx,
		Counterparty: &createCounterpartyRequest,
	}
	_, err = client.Counterparties.PostV1Counterparties(&postCounterpartiesParams)
	if err != nil {
		return err
	}

	// createExternalAccountRequest := atlar_models.CreateExternalAccountRequest{
	// 	CounterpartyID: createCounterpartiesReponse.Payload.ID,
	// }
	// postExternalAccountsParams := external_accounts.PostV1ExternalAccountsParams{
	// 	Context:         detachedCtx,
	// 	ExternalAccount: &createExternalAccountRequest,
	// }

	// _, err = client.ExternalAccounts.PostV1ExternalAccounts(&postExternalAccountsParams)
	// if err != nil {
	// 	return err
	// }

	return nil
}

// TODO: validation (also metadata) needs to return a 400
func validateExternalBankAccount(newExternalBankAccount *models.BankAccount) error {
	// if newExternalBankAccount.SwiftBicCode == "" {
	// 	return errors.New("swiftBicCode must be provided")
	// }
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

	// hasIdentifier := false
	// for k := range newExternalBankAccount.Metadata {
	// 	// check whether the key has format com.atlar.spec/identifier/<market>/<type>
	// 	_, err := metadataToIdentifierData(k, newExternalBankAccount.Metadata[k])
	// 	if err == nil {
	// 		hasIdentifier = true
	// 		break
	// 	}
	// }
	// if !hasIdentifier {
	// 	return fmt.Errorf("at least one metadata field in the form of %sidentifier/<market>/<type> with a value of the identifier value is needed (e.g. %sidentifier/DE/IBAN with an IBAN as a value)", atlarMetadataSpecNamespace, atlarMetadataSpecNamespace)
	// }

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
	// Define a regular expression for matching the format

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
