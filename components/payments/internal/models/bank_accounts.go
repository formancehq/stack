package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	bankAccountOwnerNamespace = formanceMetadataSpecNamespace + "owner/"

	BankAccountOwnerAddressLine1MetadataKey = bankAccountOwnerNamespace + "addressLine1"
	BankAccountOwnerAddressLine2MetadataKey = bankAccountOwnerNamespace + "addressLine2"
	BankAccountOwnerCityMetadataKey         = bankAccountOwnerNamespace + "city"
	BankAccountOwnerRegionMetadataKey       = bankAccountOwnerNamespace + "region"
	BankAccountOwnerPostalCodeMetadataKey   = bankAccountOwnerNamespace + "postalCode"
)

type BankAccount struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	Name      string    `json:"name"`

	AccountNumber *string `json:"accountNumber"`
	IBAN          *string `json:"iban"`
	SwiftBicCode  *string `json:"swiftBicCode"`
	Country       *string `json:"country"`

	Metadata map[string]string `json:"metadata"`

	RelatedAccounts []BankAccountRelatedAccount `json:"relatedAccounts"`
}

func (a *BankAccount) Offuscate() error {
	if a.IBAN != nil {
		length := len(*a.IBAN)
		if length < 8 {
			return errors.New("IBAN is not valid")
		}

		*a.IBAN = (*a.IBAN)[:4] + strings.Repeat("*", length-8) + (*a.IBAN)[length-4:]
	}

	if a.AccountNumber != nil {
		length := len(*a.AccountNumber)
		if length < 5 {
			return errors.New("Account number is not valid")
		}

		*a.AccountNumber = (*a.AccountNumber)[:2] + strings.Repeat("*", length-5) + (*a.AccountNumber)[length-3:]
	}

	return nil
}
