package messages

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/webhooks/pkg/healthcheck"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func newWorkerMessagesHandler() http.Handler {
	h := httprouter.New()
	h.GET(PathHealthCheck, healthcheck.Handle)

	return h
}
