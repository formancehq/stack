package models

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/formancehq/stack/libs/go-libs/pointer"
)

// Internal struct used by the plugins
type PSPPayment struct {
	// PSP payment/transaction reference. Should be unique.
	Reference string

	// Payment Creation date.
	CreatedAt time.Time

	// Type of payment: payin, payout, transfer etc...
	Type PaymentType

	// Payment amount.
	Amount *big.Int

	// Currency. Should be in minor currencies unit.
	// For example: USD/2
	Asset string

	// Payment scheme if existing: visa, mastercard etc...
	Scheme PaymentScheme

	// Payment status: pending, failed, succeeded etc...
	Status PaymentStatus

	// Optional, can be filled for payouts and transfers for example.
	SourceAccountReference *string
	// Optional, can be filled for payins and transfers for example.
	DestinationAccountReference *string

	// Additional metadata
	Metadata map[string]string

	// PSP response in raw
	Raw json.RawMessage
}

type Payment struct {
	// Unique Payment ID generated from payments information
	ID PaymentID `json:"id"`
	// Related Connector ID
	ConnectorID ConnectorID `json:"connectorID"`

	// PSP payment/transaction reference. Should be unique.
	Reference string `json:"reference"`

	// Payment Creation date.
	CreatedAt time.Time `json:"createdAt"`

	// Type of payment: payin, payout, transfer etc...
	Type PaymentType `json:"type"`

	// Payment Initial amount
	InitialAmount *big.Int `json:"initialAmount"`
	// Payment amount.
	Amount *big.Int `json:"amount"`

	// Currency. Should be in minor currencies unit.
	// For example: USD/2
	Asset string `json:"asset"`

	// Payment scheme if existing: visa, mastercard etc...
	Scheme PaymentScheme `json:"scheme"`

	// Payment status: pending, failed, succeeded etc...
	Status PaymentStatus `json:"status"`

	// Optional, can be filled for payouts and transfers for example.
	SourceAccountID *AccountID `json:"sourceAccountID"`
	// Optional, can be filled for payins and transfers for example.
	DestinationAccountID *AccountID `json:"destinationAccountID"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`

	// Related adjustment
	Adjustments []PaymentAdjustment `json:"adjustments"`
}

func (p Payment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID                   string              `json:"id"`
		ConnectorID          string              `json:"connectorID"`
		Reference            string              `json:"reference"`
		CreatedAt            time.Time           `json:"createdAt"`
		Type                 PaymentType         `json:"type"`
		InitialAmount        *big.Int            `json:"initialAmount"`
		Amount               *big.Int            `json:"amount"`
		Asset                string              `json:"asset"`
		Scheme               PaymentScheme       `json:"scheme"`
		Status               PaymentStatus       `json:"status"`
		SourceAccountID      *string             `json:"sourceAccountID"`
		DestinationAccountID *string             `json:"destinationAccountID"`
		Metadata             map[string]string   `json:"metadata"`
		Adjustments          []PaymentAdjustment `json:"adjustments"`
	}{
		ID:            p.ID.String(),
		ConnectorID:   p.ConnectorID.String(),
		Reference:     p.Reference,
		CreatedAt:     p.CreatedAt,
		Type:          p.Type,
		InitialAmount: p.InitialAmount,
		Amount:        p.Amount,
		Asset:         p.Asset,
		Scheme:        p.Scheme,
		Status:        p.Status,
		SourceAccountID: func() *string {
			if p.SourceAccountID == nil {
				return nil
			}
			return pointer.For(p.SourceAccountID.String())
		}(),
		DestinationAccountID: func() *string {
			if p.DestinationAccountID == nil {
				return nil
			}
			return pointer.For(p.DestinationAccountID.String())
		}(),
		Metadata:    p.Metadata,
		Adjustments: p.Adjustments,
	})
}

func (c *Payment) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID                   string              `json:"id"`
		ConnectorID          string              `json:"connectorID"`
		Reference            string              `json:"reference"`
		CreatedAt            time.Time           `json:"createdAt"`
		Type                 PaymentType         `json:"type"`
		InitialAmount        *big.Int            `json:"initialAmount"`
		Amount               *big.Int            `json:"amount"`
		Asset                string              `json:"asset"`
		Scheme               PaymentScheme       `json:"scheme"`
		Status               PaymentStatus       `json:"status"`
		SourceAccountID      *string             `json:"sourceAccountID"`
		DestinationAccountID *string             `json:"destinationAccountID"`
		Metadata             map[string]string   `json:"metadata"`
		Adjustments          []PaymentAdjustment `json:"adjustments"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	id, err := PaymentIDFromString(aux.ID)
	if err != nil {
		return err
	}

	connectorID, err := ConnectorIDFromString(aux.ConnectorID)
	if err != nil {
		return err
	}

	var sourceAccountID *AccountID
	if aux.SourceAccountID != nil {
		id, err := AccountIDFromString(*aux.SourceAccountID)
		if err != nil {
			return err
		}
		sourceAccountID = &id
	}

	var destinationAccountID *AccountID
	if aux.DestinationAccountID != nil {
		id, err := AccountIDFromString(*aux.DestinationAccountID)
		if err != nil {
			return err
		}
		destinationAccountID = &id
	}

	c.ID = id
	c.ConnectorID = connectorID
	c.Reference = aux.Reference
	c.CreatedAt = aux.CreatedAt
	c.Type = aux.Type
	c.InitialAmount = aux.InitialAmount
	c.Amount = aux.Amount
	c.Asset = aux.Asset
	c.Scheme = aux.Scheme
	c.Status = aux.Status
	c.SourceAccountID = sourceAccountID
	c.DestinationAccountID = destinationAccountID
	c.Metadata = aux.Metadata
	c.Adjustments = aux.Adjustments

	return nil
}

func FromPSPPaymentToPayment(from PSPPayment, connectorID ConnectorID) Payment {
	return Payment{
		ID: PaymentID{
			PaymentReference: PaymentReference{
				Reference: from.Reference,
				Type:      from.Type,
			},
			ConnectorID: connectorID,
		},
		ConnectorID:   connectorID,
		Reference:     from.Reference,
		CreatedAt:     from.CreatedAt,
		Type:          from.Type,
		InitialAmount: from.Amount,
		Amount:        from.Amount,
		Asset:         from.Asset,
		Scheme:        from.Scheme,
		Status:        from.Status,
		SourceAccountID: func() *AccountID {
			if from.SourceAccountReference == nil {
				return nil
			}
			return &AccountID{
				Reference:   *from.SourceAccountReference,
				ConnectorID: connectorID,
			}
		}(),
		DestinationAccountID: func() *AccountID {
			if from.DestinationAccountReference == nil {
				return nil
			}
			return &AccountID{
				Reference:   *from.DestinationAccountReference,
				ConnectorID: connectorID,
			}
		}(),
		Metadata: from.Metadata,
	}
}

func FromPSPPayments(from []PSPPayment, connectorID ConnectorID) []Payment {
	payments := make([]Payment, 0, len(from))
	for _, p := range from {
		payment := FromPSPPaymentToPayment(p, connectorID)
		payment.Adjustments = append(payment.Adjustments, FromPSPPaymentToPaymentAdjustement(p, connectorID))
		payments = append(payments, payment)
	}
	return payments
}

func FromPSPPaymentToPaymentAdjustement(from PSPPayment, connectorID ConnectorID) PaymentAdjustment {
	paymentID := PaymentID{
		PaymentReference: PaymentReference{
			Reference: from.Reference,
			Type:      from.Type,
		},
		ConnectorID: connectorID,
	}

	return PaymentAdjustment{
		ID: PaymentAdjustmentID{
			PaymentID: paymentID,
			CreatedAt: from.CreatedAt,
			Status:    from.Status,
		},
		PaymentID: paymentID,
		CreatedAt: from.CreatedAt,
		Status:    from.Status,
		Amount:    from.Amount,
		Asset:     &from.Asset,
		Metadata:  from.Metadata,
		Raw:       from.Raw,
	}
}
