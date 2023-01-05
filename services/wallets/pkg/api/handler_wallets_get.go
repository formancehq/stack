package api

import (
	"errors"
	"net/http"

	"github.com/formancehq/wallets/pkg/wallet"
	"github.com/go-chi/chi/v5"
)

func (m *MainHandler) GetWalletHandler(wr http.ResponseWriter, r *http.Request) {
	w, err := m.repository.GetWallet(r.Context(), chi.URLParam(r, "walletID"))
	if err != nil {
		switch {
		case errors.Is(err, wallet.ErrWalletNotFound):
			notFound(wr)
		default:
			internalError(wr, r, err)
		}
		return
	}

	ok(wr, w)
}
