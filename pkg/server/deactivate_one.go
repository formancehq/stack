package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/formancehq/go-libs/sharedapi"
	"github.com/formancehq/go-libs/sharedlogging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/go-chi/chi/v5"
)

func (h *serverHandler) deactivateOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	cursor, err := updateOneConfigActivation(r.Context(), false, id, h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, id, PathDeactivate)
		resp := sharedapi.BaseResponse[webhooks.Config]{
			Cursor: &cursor,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if errors.Is(err, ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, id, PathDeactivate, ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else if errors.Is(err, ErrConfigNotModified) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, id, PathDeactivate, ErrConfigNotModified)
		w.WriteHeader(http.StatusNotModified)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, id, PathDeactivate, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
