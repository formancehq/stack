package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BankAccount struct {
	bun.BaseModel `bun:"accounts.bank_account"`

	ID            uuid.UUID `bun:",pk,nullzero"`
	CreatedAt     time.Time `bun:",nullzero"`
	Name          string
	AccountNumber string `bun:"decrypted_account_number,scanonly"`
	IBAN          string `bun:"decrypted_iban,scanonly"`
	SwiftBicCode  string `bun:"decrypted_swift_bic_code,scanonly"`
	Country       string `bun:"country"`
	Metadata      map[string]string

	RelatedAccounts []*BankAccountRelatedAccount `bun:"rel:has-many,join:id=bank_account_id"`
}

func (a *BankAccount) Offuscate() error {
	if a.IBAN != "" {
		length := len(a.IBAN)
		if length < 8 {
			return errors.New("IBAN is not valid")
		}

		a.IBAN = a.IBAN[:4] + strings.Repeat("*", length-8) + a.IBAN[length-4:]
	}

	if a.AccountNumber != "" {
		length := len(a.AccountNumber)
		if length < 5 {
			return errors.New("Account number is not valid")
		}

		a.AccountNumber = a.AccountNumber[:2] + strings.Repeat("*", length-5) + a.AccountNumber[length-3:]
	}

	return nil
}

type BankAccountRelatedAccount struct {
	bun.BaseModel `bun:"accounts.bank_account_related_accounts"`

	ID            uuid.UUID   `bun:",pk,nullzero"`
	CreatedAt     time.Time   `bun:",nullzero"`
	BankAccountID uuid.UUID   `bun:",nullzero"`
	ConnectorID   ConnectorID `bun:",nullzero"`
	AccountID     AccountID   `bun:",nullzero"`
}

const (
	bankAccountOwnerNamespace = formanceMetadataSpecNamespace + "owner/"

	BankAccountOwnerAddressLine1MetadataKey = bankAccountOwnerNamespace + "addressLine1"
	BankAccountOwnerAddressLine2MetadataKey = bankAccountOwnerNamespace + "addressLine2"
	BankAccountOwnerCityMetadataKey         = bankAccountOwnerNamespace + "city"
	BankAccountOwnerRegionMetadataKey       = bankAccountOwnerNamespace + "region"
	BankAccountOwnerPostalCodeMetadataKey   = bankAccountOwnerNamespace + "postalCode"
)
