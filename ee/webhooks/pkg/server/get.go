package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	"github.com/formancehq/stack/libs/go-libs/bun/bunpaginate"

	"github.com/formancehq/stack/libs/go-libs/api"
	"github.com/formancehq/stack/libs/go-libs/logging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/formancehq/webhooks/pkg/server/apierrors"
)

func (h *serverHandler) getManyConfigsHandle(w http.ResponseWriter, r *http.Request) {
	filter, err := buildQueryFilter(r.URL.Query())
	if err != nil {
		apierrors.ResponseError(w, r, apierrors.NewValidationError(err.Error()))
		return
	}

	cfgs, err := h.store.FindManyConfigs(r.Context(), filter)
	if err != nil {
		logging.FromContext(r.Context()).Errorf("storage.store.FindManyConfigs: %s", err)
		apierrors.ResponseError(w, r, err)
		return
	}

	resp := api.BaseResponse[webhooks.Config]{
		Cursor: &bunpaginate.Cursor[webhooks.Config]{
			Data: cfgs,
		},
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logging.FromContext(r.Context()).Errorf("json.Encoder.Encode: %s", err)
		apierrors.ResponseError(w, r, err)
		return
	}

	logging.FromContext(r.Context()).Debugf("GET /configs: %d results", len(resp.Cursor.Data))
}

var ErrInvalidParams = errors.New("invalid params: only 'id' and 'endpoint' with a valid URL are accepted")

func buildQueryFilter(values url.Values) (map[string]any, error) {
	filter := map[string]any{}

	for key, value := range values {
		if len(value) != 1 {
			return nil, ErrInvalidParams
		}
		switch key {
		case "id":
			filter["id"] = value[0]
		case "endpoint":
			if u, err := url.Parse(value[0]); err != nil {
				return nil, ErrInvalidParams
			} else {
				filter["endpoint"] = u.String()
			}
		default:
			return nil, ErrInvalidParams
		}
	}

	return filter, nil
}
