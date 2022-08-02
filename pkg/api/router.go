package api

import (
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/gorilla/mux"
	"github.com/numary/auth/pkg/delegatedauth"
	"github.com/numary/auth/pkg/storage"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"github.com/zitadel/oidc/pkg/op"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	_ "go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func NewRouter(provider op.OpenIDProvider, storage storage.Storage, healthController *sharedhealth.HealthController,
	delegatedOAuth2Config delegatedauth.OAuth2Config, delegatedOIDCProvider *oidc.Provider) *mux.Router {
	router := provider.HttpHandler().(*mux.Router)
	router.Use(otelmux.Middleware("auth"))
	router.Path("/_healthcheck").HandlerFunc(healthController.Check)
	router.Path("/delegatedoidc/callback").Handler(authorizeCallbackHandler(
		provider, storage, delegatedOAuth2Config, delegatedOIDCProvider))
	router.Path("/clients").Methods(http.MethodPost).HandlerFunc(createClient)
	router.Path("/clients").Methods(http.MethodGet).HandlerFunc(listClients)
	router.Path("/clients/{clientId}").Methods(http.MethodPut).HandlerFunc(updateClient)
	router.Path("/clients/{clientId}").Methods(http.MethodPatch).HandlerFunc(patchClient)
	router.Path("/clients/{clientId}").Methods(http.MethodGet).HandlerFunc(readClient)
	router.Path("/clients/{clientId}/secrets").Methods(http.MethodPost).HandlerFunc(createSecret)
	router.Path("/clients/{clientId}/secrets/{secretId}").Methods(http.MethodDelete).HandlerFunc(deleteSecret)
	router.Path("/clients/{clientId}/app/scopes").Methods(http.MethodPut).HandlerFunc(updateClientScopes)
	return router
}
