package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/numary/auth/pkg/storage"
	sharedhealth "github.com/numary/go-libs/sharedhealth/pkg"
	"github.com/zitadel/oidc/pkg/client/rp"
	"github.com/zitadel/oidc/pkg/op"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
	"gorm.io/gorm"
)

func NewRouter(provider op.OpenIDProvider, storage storage.Storage, healthController *sharedhealth.HealthController,
	relyingParty rp.RelyingParty, db *gorm.DB) *mux.Router {
	router := provider.HttpHandler().(*mux.Router)
	router.Use(otelmux.Middleware("auth"))
	router.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	router.Path("/_healthcheck").HandlerFunc(healthController.Check)
	router.Path("/delegatedoidc/callback").Handler(authorizeCallbackHandler(
		provider, storage, relyingParty))
	addClientRoutes(db, router)
	addScopeRoutes(db, router)
	return router
}
