package api

import (
	"net/http"

	"github.com/formancehq/payments/cmd/connectors/internal/api/backend"
)

func healthHandler(b backend.ServiceBackend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := b.GetService().Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func liveHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}
