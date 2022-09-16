package messages

import (
	"net/http"

	"github.com/formancehq/webhooks/pkg/healthcheck"
	"github.com/julienschmidt/httprouter"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func newWorkerMessagesHandler() http.Handler {
	h := httprouter.New()
	h.GET(PathHealthCheck, healthcheck.Handle)

	return h
}
