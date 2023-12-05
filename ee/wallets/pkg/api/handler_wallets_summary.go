package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (m *MainHandler) walletSummaryHandler(w http.ResponseWriter, r *http.Request) {
	summary, err := m.manager.GetWalletSummary(r.Context(), chi.URLParam(r, "walletID"))
	if err != nil {
		internalError(w, r, err)
		return
	}

	ok(w, summary)
}
