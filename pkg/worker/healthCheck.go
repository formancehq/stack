package worker

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedlogging"
)

func healthCheckHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}
