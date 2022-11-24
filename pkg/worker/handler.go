package worker

import (
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/julienschmidt/httprouter"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func newWorkerHandler() http.Handler {
	h := httprouter.New()
	h.GET(PathHealthCheck, healthCheckHandle)

	return h
}

func healthCheckHandle(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}
