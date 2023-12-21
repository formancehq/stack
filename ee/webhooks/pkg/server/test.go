package server

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/backoff"
	"github.com/formancehq/webhooks/pkg/server/apierrors"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func (h *serverHandler) testOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	cfgs, err := h.store.FindManyConfigs(r.Context(), map[string]any{"id": id})
	if err == nil {
		if len(cfgs) == 0 {
			logging.FromContext(r.Context()).Errorf("GET %s/%s%s: %s", PathConfigs, id, PathTest, storage.ErrConfigNotFound)
			apierrors.ResponseError(w, r, apierrors.NewNotFoundError(storage.ErrConfigNotFound.Error()))
			return
		}
		logging.FromContext(r.Context()).Debugf("GET %s/%s%s", PathConfigs, id, PathTest)
		retryPolicy := backoff.NewNoRetry()
		attempt, err := webhooks.MakeAttempt(r.Context(), h.httpClient, retryPolicy, uuid.NewString(),
			uuid.NewString(), 0, cfgs[0], []byte(`{"data":"test"}`), true)
		if err != nil {
			logging.FromContext(r.Context()).Errorf("GET %s/%s%s: %s", PathConfigs, id, PathTest, err)
			apierrors.ResponseError(w, r, err)
		} else {
			logging.FromContext(r.Context()).Debugf("GET %s/%s%s", PathConfigs, id, PathTest)
			resp := api.BaseResponse[webhooks.Attempt]{
				Data: &attempt,
			}
			if err := json.NewEncoder(w).Encode(resp); err != nil {
				logging.FromContext(r.Context()).Errorf("json.Encoder.Encode: %s", err)
				apierrors.ResponseError(w, r, err)
				return
			}
		}
	} else {
		logging.FromContext(r.Context()).Errorf("GET %s/%s%s: %s", PathConfigs, id, PathTest, err)
		apierrors.ResponseError(w, r, err)
	}
}
