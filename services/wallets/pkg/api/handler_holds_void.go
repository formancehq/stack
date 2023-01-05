package api

import (
	"errors"
	"net/http"

	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
)

func (m *MainHandler) VoidHoldHandler(w http.ResponseWriter, r *http.Request) {
	err := m.funding.VoidHold(r.Context(), wallet.VoidHold{
		HoldID: chi.URLParam(r, "holdID"),
	})
	if err != nil {
		switch {
		case errors.Is(err, wallet.ErrClosedHold):
			badRequest(w, ErrorCodeClosedHold, err)
		default:
			internalError(w, r, err)
		}
		return
	}

	noContent(w)
}
