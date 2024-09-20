package api

import (
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/payments/cmd/api/internal/api/backend"
)

func healthHandler(b backend.Backend) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := b.GetService().Ping(); err != nil {
			api.InternalServerError(w, r, err)

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
