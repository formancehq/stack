package oidc

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
)

func AddRoutes(router *mux.Router, provider op.OpenIDProvider, storage Storage, relyingParty rp.RelyingParty, baseUrl *url.URL) {
	router.NewRoute().Path("/authorize/callback").Queries("code", "{code}").
		Handler(authorizeCallbackHandler(provider, storage, relyingParty))
	router.PathPrefix("/").Handler(http.StripPrefix(baseUrl.Path, provider.HttpHandler()))
}
