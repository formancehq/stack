package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/logging"
	"github.com/formancehq/webhooks/pkg/server/apierrors"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/pkg/errors"
)

func (h *serverHandler) deleteOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	err := h.store.DeleteOneConfig(r.Context(), id)
	if err == nil {
		logging.FromContext(r.Context()).Debugf("DELETE %s/%s", PathConfigs, id)
	} else if errors.Is(err, storage.ErrConfigNotFound) {
		logging.FromContext(r.Context()).Debugf("DELETE %s/%s: %s", PathConfigs, id, storage.ErrConfigNotFound)
		apierrors.ResponseError(w, r, apierrors.NewNotFoundError(storage.ErrConfigNotFound.Error()))
	} else {
		logging.FromContext(r.Context()).Errorf("DELETE %s/%s: %s", PathConfigs, id, err)
		apierrors.ResponseError(w, r, err)
	}
}
