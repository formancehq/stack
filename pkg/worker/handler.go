package worker

import (
	"net/http"

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
