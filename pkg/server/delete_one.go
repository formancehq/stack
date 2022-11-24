package server

import (
	"errors"
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/go-chi/chi/v5"
)

func (h *serverHandler) deleteOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	err := deleteOneConfig(r.Context(), id, h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s", PathConfigs, id)
	} else if errors.Is(err, ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s: %s", PathConfigs, id, ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("DELETE %s/%s: %s", PathConfigs, id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
