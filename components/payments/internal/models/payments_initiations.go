package models

import (
	"math/big"
	"time"
)

const (
	PaymentInitiationTypeTransfer string = "TRANSFER"
	PaymentInitiationTypePayout   string = "PAYOUT"
)

type PaymentInitiation struct {
	// Unique Payment initiation ID generated from payments information
	ID PaymentInitiationID `json:"id"`
	// Related Connector ID
	ConnectorID ConnectorID `json:"connectorID"`
	// Unique reference of the payment
	Reference string `json:"reference"`

	// Time to schedule the payment
	ScheduledAt time.Time `json:"scheduledAt"`

	// Description of the payment
	Description string `json:"description"`

	// Source account of the payment
	SourceAccountID *AccountID `json:"sourceAccountID"`
	// Destination account of the payment
	DestinationAccountID AccountID `json:"destinationAccountID"`

	// Payment initial amount
	InitialAmount *big.Int `json:"initialAmount"`
	// Payment current amount (can be changed of reversed, refunded, etc...)
	Amount *big.Int `json:"amount"`
	// Asset of the payment
	Asset string `json:"asset"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`
}
