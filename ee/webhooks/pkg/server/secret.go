package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/api"
	"github.com/formancehq/go-libs/logging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/server/apierrors"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/pkg/errors"
)

func (h *serverHandler) changeSecretHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	sec := webhooks.Secret{}
	if err := decodeJSONBody(r, &sec, true); err != nil {
		logging.FromContext(r.Context()).Errorf("decodeJSONBody: %s", err)
		apierrors.ResponseError(w, r, apierrors.NewValidationError(err.Error()))
		return
	}

	if err := sec.Validate(); err != nil {
		logging.FromContext(r.Context()).Errorf("invalid secret: %s", err)
		apierrors.ResponseError(w, r, apierrors.NewValidationError(err.Error()))
		return
	}

	c, err := h.store.UpdateOneConfigSecret(r.Context(), id, sec.Secret)
	if err == nil || errors.Is(err, storage.ErrConfigNotModified) {
		logging.FromContext(r.Context()).Debugf("PUT %s/%s%s", PathConfigs, id, PathChangeSecret)
		resp := api.BaseResponse[webhooks.Config]{
			Data: &c,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logging.FromContext(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			apierrors.ResponseError(w, r, err)
			return
		}
	} else if errors.Is(err, storage.ErrConfigNotFound) {
		logging.FromContext(r.Context()).Debugf("PUT %s/%s%s: %s", PathConfigs, id, PathChangeSecret, storage.ErrConfigNotFound)
		apierrors.ResponseError(w, r, apierrors.NewNotFoundError(storage.ErrConfigNotFound.Error()))
	} else {
		logging.FromContext(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, id, PathChangeSecret, err)
		apierrors.ResponseError(w, r, err)
	}
}
