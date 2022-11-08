package server

import (
	"encoding/json"
	"errors"
	"net/http"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
)

func (h *serverHandler) activateOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cursor, err := updateOneConfigActivation(r.Context(), true, p.ByName(PathParamId), h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, p.ByName(PathParamId), PathActivate)
		resp := sharedapi.BaseResponse[webhooks.Config]{
			Cursor: &cursor,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if errors.Is(err, ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathActivate, ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else if errors.Is(err, ErrConfigNotModified) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathActivate, ErrConfigNotModified)
		w.WriteHeader(http.StatusNotModified)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathActivate, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
