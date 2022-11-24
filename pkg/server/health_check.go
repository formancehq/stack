package server

import (
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
)

func (h *serverHandler) HealthCheckHandle(_ http.ResponseWriter, r *http.Request) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}
