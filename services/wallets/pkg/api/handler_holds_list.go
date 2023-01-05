package api

import (
	"net/http"

	"github.com/formancehq/wallets/pkg/wallet"
)

func (m *MainHandler) ListHoldsHandler(w http.ResponseWriter, r *http.Request) {
	query := readPaginatedRequest(r, func(r *http.Request) wallet.ListHolds {
		return wallet.ListHolds{
			WalletID: r.URL.Query().Get("walletID"),
			Metadata: getQueryMap(r.URL.Query(), "metadata"),
		}
	})

	holds, err := m.repository.ListHolds(r.Context(), query)
	if err != nil {
		internalError(w, r, err)
		return
	}

	cursorFromListResponse(w, query, holds)
}
