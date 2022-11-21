package server

import (
	"errors"
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/julienschmidt/httprouter"
)

func (h *serverHandler) deleteOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	err := deleteOneConfig(r.Context(), p.ByName(PathParamId), h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s", PathConfigs, p.ByName(PathParamId))
	} else if errors.Is(err, ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s: %s", PathConfigs, p.ByName(PathParamId), ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("DELETE %s/%s: %s", PathConfigs, p.ByName(PathParamId), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
