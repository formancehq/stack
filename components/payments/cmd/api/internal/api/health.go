package api

import (
	"net/http"

	"github.com/formancehq/payments/cmd/api/internal/api/backend"
	"github.com/formancehq/stack/libs/go-libs/api"
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
