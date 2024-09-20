package api

import (
	"context"
	"net/http"

	"github.com/formancehq/go-libs/api"
	"github.com/pkg/errors"
)

func recoveryHandler(reporter func(ctx context.Context, e interface{})) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if e := recover(); e != nil {
					api.InternalServerError(w, r, errors.New("Internal Server Error"))
					reporter(r.Context(), e)
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
