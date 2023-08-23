package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type BankAccount struct {
	bun.BaseModel `bun:"accounts.bank_account"`

	ID            uuid.UUID `bun:",pk,nullzero"`
	CreatedAt     time.Time `bun:",nullzero"`
	Provider      ConnectorProvider
	Name          string
	AccountNumber string `bun:"decrypted_account_number,scanonly"`
	IBAN          string `bun:"decrypted_iban,scanonly"`
	SwiftBicCode  string `bun:"decrypted_swift_bic_code,scanonly"`
	Country       string `bun:"country"`
}
