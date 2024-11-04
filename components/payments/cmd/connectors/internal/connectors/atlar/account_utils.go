package atlar

import (
	"encoding/json"
	"fmt"

	"github.com/formancehq/payments/internal/models"
	"github.com/formancehq/stack/libs/go-libs/metadata"
	atlar_models "github.com/get-momo/atlar-v1-go-client/models"
)

type AtlarExternalAccountAndCounterparty struct {
	ExternalAccount atlar_models.ExternalAccount `json:"externalAccount" yaml:"externalAccount" bson:"externalAccount"`
	Counterparty    atlar_models.Counterparty    `json:"counterparty" yaml:"counterparty" bson:"counterparty"`
}

func ExternalAccountFromAtlarData(
	connectorID models.ConnectorID,
	externalAccount *atlar_models.ExternalAccount,
	counterparty *atlar_models.Counterparty,
) (*models.Account, error) {
	raw, err := json.Marshal(AtlarExternalAccountAndCounterparty{ExternalAccount: *externalAccount, Counterparty: *counterparty})
	if err != nil {
		return nil, err
	}

	createdAt, err := ParseAtlarTimestamp(externalAccount.Created)
	if err != nil {
		return nil, fmt.Errorf("failed to parse opening date: %w", err)
	}

	return &models.Account{
		ID: models.AccountID{
			Reference:   externalAccount.ID,
			ConnectorID: connectorID,
		},
		CreatedAt:   createdAt,
		Reference:   externalAccount.ID,
		ConnectorID: connectorID,
		// DefaultAsset: left empty because the information is not provided by Atlar,
		AccountName: counterparty.Name, // TODO: is that okay? External accounts do not have a name at Atlar.
		Type:        models.AccountTypeExternal,
		Metadata:    extractExternalAccountAndCounterpartyMetadata(externalAccount, counterparty),
		RawData:     raw,
	}, nil
}

func ExtractAccountMetadata(account *atlar_models.Account, bank *atlar_models.ThirdParty) metadata.Metadata {
	result := metadata.Metadata{}
	result = result.Merge(ComputeAccountMetadataBool("fictive", account.Fictive))
	result = result.Merge(ComputeAccountMetadata("bank/id", bank.ID))
	result = result.Merge(ComputeAccountMetadata("bank/name", bank.Name))
	result = result.Merge(ComputeAccountMetadata("bank/bic", account.Bank.Bic))
	result = result.Merge(IdentifiersToMetadata(account.Identifiers))
	result = result.Merge(ComputeAccountMetadata("alias", account.Alias))
	result = result.Merge(ComputeAccountMetadata("owner/name", account.Owner.Name))
	return result
}

func IdentifiersToMetadata(identifiers []*atlar_models.AccountIdentifier) metadata.Metadata {
	result := metadata.Metadata{}
	for _, i := range identifiers {
		result = result.Merge(ComputeAccountMetadata(
			fmt.Sprintf("identifier/%s/%s", *i.Market, *i.Type),
			*i.Number,
		))
		if *i.Type == "IBAN" {
			result = result.Merge(ComputeAccountMetadata(
				fmt.Sprintf("identifier/%s", *i.Type),
				*i.Number,
			))
		}
	}
	return result
}

func extractExternalAccountAndCounterpartyMetadata(externalAccount *atlar_models.ExternalAccount, counterparty *atlar_models.Counterparty) metadata.Metadata {
	result := metadata.Metadata{}
	result = result.Merge(ComputeAccountMetadata("bank/id", externalAccount.Bank.ID))
	result = result.Merge(ComputeAccountMetadata("bank/name", externalAccount.Bank.Name))
	result = result.Merge(ComputeAccountMetadata("bank/bic", externalAccount.Bank.Bic))
	result = result.Merge(IdentifiersToMetadata(externalAccount.Identifiers))
	result = result.Merge(ComputeAccountMetadata("owner/name", counterparty.Name))
	result = result.Merge(ComputeAccountMetadata("owner/type", counterparty.PartyType))
	result = result.Merge(ComputeAccountMetadata("owner/contact/email", counterparty.ContactDetails.Email))
	result = result.Merge(ComputeAccountMetadata("owner/contact/phone", counterparty.ContactDetails.Phone))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/streetName", counterparty.ContactDetails.Address.StreetName))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/streetNumber", counterparty.ContactDetails.Address.StreetNumber))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/city", counterparty.ContactDetails.Address.City))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/postalCode", counterparty.ContactDetails.Address.PostalCode))
	result = result.Merge(ComputeAccountMetadata("owner/contact/address/country", counterparty.ContactDetails.Address.Country))
	return result
}
