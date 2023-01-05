package wallet

import (
	"net/http"

	"github.com/formancehq/go-libs/metadata"
)

var DefaultDebitDest = NewLedgerAccountSubject("world")

type DebitRequest struct {
	Amount      Monetary          `json:"amount"`
	Pending     bool              `json:"pending"`
	Metadata    metadata.Metadata `json:"metadata"`
	Description string            `json:"description"`
	Reference   string            `json:"reference"`
	Destination *Subject          `json:"destination"`
	Balance     string            `json:"balance"`
}

func (c *DebitRequest) Bind(r *http.Request) error {
	return nil
}

type Debit struct {
	DebitRequest
	WalletID string `json:"walletID"`
}

func (d Debit) newHold(chart *Chart) DebitHold {
	md := d.Metadata
	if md == nil {
		md = metadata.Metadata{}
	}
	return NewDebitHold(
		d.WalletID,
		d.getDestination().getAccount(chart),
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

func (d Debit) sourceAccount(chart *Chart) string {
	if d.Balance == "" {
		return chart.GetMainBalanceAccount(d.WalletID)
	}
	return chart.GetBalanceAccount(d.WalletID, d.Balance)
}
