package worker

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/go-chi/chi/v5"
	"github.com/riandyrn/otelchi"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func NewWorkerHandler() http.Handler {
	h := chi.NewRouter()
	h.Use(otelchi.Middleware("webhooks"))
	h.Get(PathHealthCheck, healthCheckHandle)

	return h
}

func healthCheckHandle(_ http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).Infof("health check OK")
}
