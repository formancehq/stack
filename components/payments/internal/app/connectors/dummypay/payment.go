package dummypay

import (
	"encoding/json"
	"math/big"

	"github.com/formancehq/payments/internal/app/models"
)

// payment represents a payment structure used in the generated files.
type payment struct {
	Reference string               `json:"reference"`
	Amount    *big.Int             `json:"amount"`
	Type      models.PaymentType   `json:"type"`
	Status    models.PaymentStatus `json:"status"`
	Scheme    models.PaymentScheme `json:"scheme"`
	Asset     models.PaymentAsset  `json:"asset"`

	RawData json.RawMessage `json:"rawData"`
}
