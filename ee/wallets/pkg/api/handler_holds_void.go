package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/api"

	wallet "github.com/formancehq/wallets/pkg"
)

func (m *MainHandler) voidHoldHandler(w http.ResponseWriter, r *http.Request) {
	err := m.manager.VoidHold(r.Context(), api.IdempotencyKeyFromRequest(r), wallet.VoidHold{
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
