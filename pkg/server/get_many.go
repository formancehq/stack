package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
)

func (h *serverHandler) getManyConfigsHandle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	filter, err := buildQueryFilter(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cursor, err := h.store.FindManyConfigs(r.Context(), filter)
	if err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("storage.store.FindManyConfigs: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := sharedapi.BaseResponse[webhooks.Config]{
		Cursor: &cursor,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.GetLogger(r.Context()).Infof("GET /configs: %d results", len(cursor.Data))
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
			filter["_id"] = value[0]
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
