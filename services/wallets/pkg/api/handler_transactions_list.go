package api

import (
	"net/http"

	"github.com/formancehq/wallets/pkg/wallet"
)

func (m *MainHandler) ListTransactions(w http.ResponseWriter, r *http.Request) {
	query := readPaginatedRequest[wallet.ListTransactions](r, func(r *http.Request) wallet.ListTransactions {
		return wallet.ListTransactions{
			WalletID: r.URL.Query().Get("walletID"),
		}
	})
	transactions, err := m.repository.ListTransactions(r.Context(), query)
	if err != nil {
		internalError(w, r, err)
		return
	}

	cursorFromListResponse(w, query, transactions)
}
