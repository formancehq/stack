package retries

import (
	"net/http"

	"github.com/formancehq/webhooks/pkg/healthcheck"
	"github.com/julienschmidt/httprouter"
)

const (
	PathHealthCheck = "/_healthcheck"
)

func newWorkerRetriesHandler() http.Handler {
	h := httprouter.New()
	h.GET(PathHealthCheck, healthcheck.Handle)

	return h
}
