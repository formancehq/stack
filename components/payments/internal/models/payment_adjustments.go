package models

import (
	"encoding/json"
	"math/big"
	"time"
)

type PaymentAdjustment struct {
	// Unique ID of the payment adjustment
	ID PaymentAdjustmentID `json:"id"`
	// Related Payment ID
	PaymentID PaymentID `json:"paymentID"`

	// Creation date of the adjustment
	CreatedAt time.Time `json:"createdAt"`

	// Status of the payment adjustement
	Status PaymentStatus `json:"status"`

	// Optional
	// Amount moved
	Amount *big.Int `json:"amount"`
	// Optional
	// Asset related to amount
	Asset *string `json:"asset"`

	// Additional metadata
	Metadata map[string]string `json:"metadata"`
	// PSP response in raw
	Raw json.RawMessage `json:"raw"`
}

func (c PaymentAdjustment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		ID        string            `json:"id"`
		PaymentID string            `json:"paymentID"`
		CreatedAt time.Time         `json:"createdAt"`
		Status    PaymentStatus     `json:"status"`
		Amount    *big.Int          `json:"amount"`
		Asset     *string           `json:"asset"`
		Metadata  map[string]string `json:"metadata"`
		Raw       json.RawMessage   `json:"raw"`
	}{
		ID:        c.ID.String(),
		PaymentID: c.PaymentID.String(),
		CreatedAt: c.CreatedAt,
		Status:    c.Status,
		Amount:    c.Amount,
		Asset:     c.Asset,
		Metadata:  c.Metadata,
		Raw:       c.Raw,
	})
}

func (c *PaymentAdjustment) UnmarshalJSON(data []byte) error {
	var aux struct {
		ID        string            `json:"id"`
		PaymentID string            `json:"paymentID"`
		CreatedAt time.Time         `json:"createdAt"`
		Status    PaymentStatus     `json:"status"`
		Amount    *big.Int          `json:"amount"`
		Asset     *string           `json:"asset"`
		Metadata  map[string]string `json:"metadata"`
		Raw       json.RawMessage   `json:"raw"`
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	paymentID, err := PaymentIDFromString(aux.PaymentID)
	if err != nil {
		return err
	}

	adjustmentID, err := PaymentAdjustmentIDFromString(aux.ID)
	if err != nil {
		return err
	}

	c.ID = *adjustmentID
	c.PaymentID = paymentID
	c.CreatedAt = aux.CreatedAt
	c.Status = aux.Status
	c.Amount = aux.Amount
	c.Asset = aux.Asset
	c.Metadata = aux.Metadata
	c.Raw = aux.Raw

	return nil
}
