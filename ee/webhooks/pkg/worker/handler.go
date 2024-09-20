package worker

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/service"

	"github.com/formancehq/go-libs/logging"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func NewWorkerHandler(debug bool) http.Handler {
	h := chi.NewRouter()
	h.Use(service.OTLPMiddleware("webhooks", debug))
	h.Get(PathHealthCheck, healthCheckHandle)

	return h
}

func healthCheckHandle(_ http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).Infof("health check OK")
}
