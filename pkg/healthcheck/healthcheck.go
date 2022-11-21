package healthcheck

import (
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/julienschmidt/httprouter"
)

func Handle(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}
