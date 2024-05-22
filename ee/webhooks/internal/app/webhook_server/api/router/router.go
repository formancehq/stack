package router

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/go-chi/chi/v5"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"

	"github.com/formancehq/webhooks/internal/app/webhook_server/api/handler"
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
	clientInterface "github.com/formancehq/webhooks/internal/services/httpclient/interfaces"
	storeInterface "github.com/formancehq/webhooks/internal/services/storage/interfaces"
)

func NewRouter(
	database storeInterface.IStoreProvider,
	client clientInterface.IHTTPClient,
	healthController *health.HealthController,
	a auth.Auth,
	info utils.ServiceInfo) chi.Router {

	handler.SetDatabase(database)
	handler.SetClientHTTP(client)

	mux := chi.NewRouter()

	mux.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})

	mux.Get("/_healthcheck", healthController.Check)
	mux.Get("/_info", func(w http.ResponseWriter, r *http.Request) {
		sharedapi.Ok(w, info)
	})

	mux.Group(func(r chi.Router) {
		r.Use(auth.Middleware(a))
		r.Use(service.OTLPMiddleware("webhooks"))

		r.Handle("/*", NewRouterV1())
		r.Handle("/v2*", NewRouterV2())

	})

	return mux
}

type MethodHTTP string

const (
	POST   MethodHTTP = "POST"
	GET    MethodHTTP = "GET"
	PUT    MethodHTTP = "PUT"
	DELETE MethodHTTP = "DELETE"
)

type Route struct {
	Method MethodHTTP
	Url    string
}

func NewRoute(m MethodHTTP, u string) Route {
	return Route{
		Method: m,
		Url:    u,
	}
}
