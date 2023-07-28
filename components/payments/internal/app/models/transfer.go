package models

import (
	"math/big"
	"time"

	"github.com/uptrace/bun"

	"github.com/google/uuid"
)

type Transfer struct {
	bun.BaseModel `bun:"payments.transfers"`

	ID          uuid.UUID `bun:",pk,nullzero"`
	ConnectorID uuid.UUID `bun:",nullzero"`
	PaymentID   *PaymentID
	CreatedAt   time.Time `bun:",nullzero"`

	Reference   *string
	Amount      *big.Int `bun:"type:numeric"`
	Status      TransferStatus
	Currency    string
	Source      string
	Destination string

	Error *string

	Payment   *Payment   `bun:"rel:has-one,join:payment_id=id"`
	Connector *Connector `bun:"rel:has-one,join:connector_id=id"`
}

type (
	TransferStatus string
)

const (
	TransferStatusPending   TransferStatus = "PENDING"
	TransferStatusSucceeded TransferStatus = "SUCCEEDED"
	TransferStatusFailed    TransferStatus = "FAILED"
)

func (t TransferStatus) String() string {
	return string(t)
}
