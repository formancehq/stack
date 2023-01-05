package api

import (
	"errors"
	"net/http"

	"github.com/formancehq/wallets/pkg/core"
	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

type ConfirmHoldRequest struct {
	Amount int64 `json:"amount"`
	Final  bool  `json:"final"`
}

func (c ConfirmHoldRequest) Bind(r *http.Request) error {
	return nil
}

func (m *MainHandler) ConfirmHoldHandler(w http.ResponseWriter, r *http.Request) {
	data := &ConfirmHoldRequest{}
	if r.ContentLength > 0 {
		if err := render.Bind(r, data); err != nil {
			badRequest(w, ErrorCodeValidation, err)
			return
		}
	}

	err := m.funding.ConfirmHold(r.Context(), wallet.ConfirmHold{
		HoldID: chi.URLParam(r, "holdID"),
		Amount: *core.NewMonetaryInt(data.Amount),
		Final:  data.Final,
	})
	if err != nil {
		switch {
		case errors.Is(err, wallet.ErrHoldNotFound):
			notFound(w)
		case errors.Is(err, wallet.ErrInsufficientFundError):
			badRequest(w, ErrorCodeInsufficientFund, err)
		case errors.Is(err, wallet.ErrClosedHold):
			badRequest(w, ErrorCodeClosedHold, err)
		default:
			internalError(w, r, err)
		}
		return
	}

	noContent(w)
}
