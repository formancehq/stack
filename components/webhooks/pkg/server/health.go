package server

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/logging"
)

func (h *serverHandler) HealthCheckHandle(_ http.ResponseWriter, r *http.Request) {
	logging.FromContext(r.Context()).Infof("health check OK")
}
