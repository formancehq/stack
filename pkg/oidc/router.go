package oidc

import (
	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
)

func AddRoutes(router *mux.Router, provider op.OpenIDProvider, storage Storage, relyingParty rp.RelyingParty) {
	router.NewRoute().Path("/authorize/callback").Queries("code", "{code}").
		Handler(authorizeCallbackHandler(provider, storage, relyingParty))
	router.NewRoute().Path("/authorize/callback").Queries("error", "{error}").
		Handler(authorizeErrorHandler())
	router.PathPrefix("/").Handler(provider.HttpHandler())
}
