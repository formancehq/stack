package worker

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/service"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/go-chi/chi/v5"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func NewWorkerHandler() http.Handler {
	h := chi.NewRouter()
	h.Use(service.OTLPMiddleware("webhooks"))
	h.Get(PathHealthCheck, healthCheckHandle)

	return h
}

func healthCheckHandle(_ http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).Infof("health check OK")
}
