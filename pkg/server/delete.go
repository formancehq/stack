package server

import (
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (h *serverHandler) deleteOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	err := h.store.DeleteOneConfig(r.Context(), id)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s", PathConfigs, id)
	} else if errors.Is(err, storage.ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s: %s", PathConfigs, id, storage.ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("DELETE %s/%s: %s", PathConfigs, id, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
