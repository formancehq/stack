package healthcheck

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedlogging"
)

func Handle(_ http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}
