package api

import (
	"net/http"

	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

const (
	ErrorCodeInternal         = "INTERNAL"
	ErrorCodeInsufficientFund = "INSUFFICIENT_FUND"
	ErrorCodeValidation       = "VALIDATION"
	ErrorCodeClosedHold       = "HOLD_CLOSED"
)

type CreditWalletRequest struct {
	Amount   core.Monetary     `json:"amount"`
	Metadata metadata.Metadata `json:"metadata"`
}

func (c *CreditWalletRequest) Bind(r *http.Request) error {
	return nil
}

func (m *MainHandler) CreditWalletHandler(w http.ResponseWriter, r *http.Request) {
	data := &CreditWalletRequest{}
	if err := render.Bind(r, data); err != nil {
		badRequest(w, ErrorCodeValidation, err)
		return
	}

	id := chi.URLParam(r, "walletID")
	credit := wallet.Credit{
		WalletID: id,
		Amount:   data.Amount,
		Metadata: data.Metadata,
	}

	err := m.funding.Credit(r.Context(), credit)
	if err != nil {
		internalError(w, r, err)
		return
	}

	noContent(w)
}
