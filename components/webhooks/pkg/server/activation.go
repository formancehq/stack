package server

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/server/apierrors"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (h *serverHandler) activateOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	c, err := h.store.UpdateOneConfigActivation(r.Context(), id, true)
	if err == nil || errors.Is(err, storage.ErrConfigNotModified) {
		logging.FromContext(r.Context()).Debugf("PUT %s/%s%s", PathConfigs, id, PathActivate)
		resp := api.BaseResponse[webhooks.Config]{
			Data: &c,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logging.FromContext(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			apierrors.ResponseError(w, r, err)
			return
		}
	} else if errors.Is(err, storage.ErrConfigNotFound) {
		logging.FromContext(r.Context()).Debugf("PUT %s/%s%s: %s", PathConfigs, id, PathActivate, storage.ErrConfigNotFound)
		apierrors.ResponseError(w, r, apierrors.NewNotFoundError(storage.ErrConfigNotFound.Error()))
	} else {
		logging.FromContext(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, id, PathActivate, err)
		apierrors.ResponseError(w, r, err)
	}
}

func (h *serverHandler) deactivateOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	c, err := h.store.UpdateOneConfigActivation(r.Context(), id, false)
	if err == nil || errors.Is(err, storage.ErrConfigNotModified) {
		logging.FromContext(r.Context()).Debugf("PUT %s/%s%s", PathConfigs, id, PathDeactivate)
		resp := api.BaseResponse[webhooks.Config]{
			Data: &c,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logging.FromContext(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			apierrors.ResponseError(w, r, err)
			return
		}
	} else if errors.Is(err, storage.ErrConfigNotFound) {
		logging.FromContext(r.Context()).Debugf("PUT %s/%s%s: %s", PathConfigs, id, PathDeactivate, storage.ErrConfigNotFound)
		apierrors.ResponseError(w, r, apierrors.NewNotFoundError(storage.ErrConfigNotFound.Error()))
	} else {
		logging.FromContext(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, id, PathDeactivate, err)
		apierrors.ResponseError(w, r, err)
	}
}
