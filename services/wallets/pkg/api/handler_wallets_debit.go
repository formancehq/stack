package api

import (
	"errors"
	"net/http"

	"github.com/formancehq/go-libs/metadata"
	"github.com/formancehq/wallets/pkg/core"
	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type DebitWalletRequest struct {
	Amount      core.Monetary     `json:"amount"`
	Pending     bool              `json:"pending"`
	Metadata    metadata.Metadata `json:"metadata"`
	Description string            `json:"description"`
}

func (c *DebitWalletRequest) Bind(r *http.Request) error {
	return nil
}

func (m *MainHandler) DebitWalletHandler(w http.ResponseWriter, r *http.Request) {
	data := &DebitWalletRequest{}
	if err := render.Bind(r, data); err != nil {
		badRequest(w, ErrorCodeValidation, err)
		return
	}

	hold, err := m.funding.Debit(r.Context(), wallet.Debit{
		WalletID:    chi.URLParam(r, "walletID"),
		Amount:      data.Amount,
		Pending:     data.Pending,
		Description: data.Description,
		Metadata:    data.Metadata,
	})
	if err != nil {
		switch {
		case errors.Is(err, wallet.ErrInsufficientFundError):
			badRequest(w, ErrorCodeInsufficientFund, wallet.ErrInsufficientFundError)
		default:
			internalError(w, r, err)
		}
		return
	}

	if hold == nil {
		noContent(w)
		return
	}

	created(w, hold)
}
