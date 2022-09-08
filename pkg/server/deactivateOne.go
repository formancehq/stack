package server

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/service"
)

func (h *serverHandler) deactivateOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := service.UpdateOneConfigActivation(r.Context(), false, p.ByName(PathParamId), h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, p.ByName(PathParamId), PathDeactivate)
	} else if errors.Is(err, service.ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathDeactivate, service.ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else if errors.Is(err, service.ErrConfigNotModified) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathDeactivate, service.ErrConfigNotModified)
		w.WriteHeader(http.StatusNotModified)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathDeactivate, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
