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

func (h *serverHandler) changeSecretHandle(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, PathParamId)
	sec := webhooks.Secret{}
	if err := decodeJSONBody(r, &sec, true); err != nil {
		var errIB *errInvalidBody
		if errors.As(err, &errIB) {
			http.Error(w, errIB.Error(), errIB.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		sharedlogging.GetLogger(r.Context()).Errorf("decodeJSONBody: %s", err)
		return
	}

	if err := sec.Validate(); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("invalid secret: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cursor, err := changeOneConfigSecret(r.Context(), id, sec.Secret, h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, id, PathChangeSecret)
		resp := sharedapi.BaseResponse[webhooks.Config]{
			Cursor: &cursor,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else if errors.Is(err, ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, id, PathChangeSecret, ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, id, PathChangeSecret, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
