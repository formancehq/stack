package wallet

import (
	"math/big"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/time"

	"github.com/formancehq/ledger/pkg/assets"

	"github.com/formancehq/stack/libs/go-libs/metadata"
)

var DefaultDebitDest = NewLedgerAccountSubject("world")

type DebitRequest struct {
	Amount      Monetary          `json:"amount"`
	Pending     bool              `json:"pending"`
	Metadata    metadata.Metadata `json:"metadata"`
	Description string            `json:"description"`
	Reference   string            `json:"reference"`
	Destination *Subject          `json:"destination"`
	Balances    []string          `json:"balances"`
	Timestamp   *time.Time        `json:"timestamp"`
}

func (c *DebitRequest) Bind(r *http.Request) error {
	return nil
}

type Debit struct {
	DebitRequest
	WalletID string `json:"walletID"`
}

func (d Debit) newHold() DebitHold {
	md := d.Metadata
	if md == nil {
		md = metadata.Metadata{}
	}
	return NewDebitHold(
		d.WalletID,
		d.getDestination(),
		d.Amount.Asset,
		d.Description,
		md,
	)
}

func (d Debit) getDestination() Subject {
	dest := DefaultDebitDest
	if d.Destination != nil {
		dest = *d.Destination
	}
	return dest
}

func (d Debit) Validate() error {
	if d.Destination != nil {
		if err := d.Destination.Validate(); err != nil {
			return err
		}
	}
	if d.Amount.Amount.Cmp(big.NewInt(0)) < 0 {
		return ErrNegativeAmount
	}
	if !assets.IsValid(d.Amount.Asset) {
		return newErrInvalidAsset(d.Amount.Asset)
	}
	return nil
}
