package api

import (
	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/numary/auth/pkg/delegatedauth"
	"github.com/numary/auth/pkg/storage"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"github.com/zitadel/oidc/pkg/op"
)

func NewRouter(provider op.OpenIDProvider, storage storage.Storage, healthController *sharedhealth.HealthController,
	delegatedOAuth2Config delegatedauth.OAuth2Config, delegatedOIDCProvider *oidc.Provider) *mux.Router {
	router := mux.NewRouter()
	router.Path("/_healthcheck").HandlerFunc(healthController.Check)
	router.Path("/authorize/callback").Handler(authorizeCallbackHandler(
		provider, storage, delegatedOAuth2Config, delegatedOIDCProvider))
	router.PathPrefix("/").Handler(provider.HttpHandler())
	return router
}
