package server

import (
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/julienschmidt/httprouter"
)

func (h *serverHandler) HealthCheckHandle(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}
