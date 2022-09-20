package oidc

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func NewRouter(provider op.OpenIDProvider, storage Storage,
	relyingParty rp.RelyingParty) *mux.Router {
	router := provider.HttpHandler().(*mux.Router)
	router.Use(otelmux.Middleware("auth"))
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	router.Path("/authorize/callback").Handler(authorizeCallbackHandler(provider, storage, relyingParty))
	return router
}
