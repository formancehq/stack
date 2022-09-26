package api

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func NewRouter(baseUrl *url.URL) *mux.Router {
	router := mux.NewRouter()
	router.Use(otelmux.Middleware("auth"))
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	return router.PathPrefix(baseUrl.Path).Subrouter()
}
