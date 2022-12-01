package server

import (
	"encoding/json"
	"net/http"

	"github.com/formancehq/go-libs/sharedapi"
	"github.com/formancehq/go-libs/sharedlogging"
	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/pkg/errors"
)

func (h *serverHandler) insertOneConfigHandle(w http.ResponseWriter, r *http.Request) {
	cfg := webhooks.ConfigUser{}
	if err := decodeJSONBody(r, &cfg, false); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("decodeJSONBody: %s", err)
		var errIB *errInvalidBody
		if errors.As(err, &errIB) {
			http.Error(w, errIB.Error(), errIB.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	if err := cfg.Validate(); err != nil {
		err := errors.Wrap(err, "invalid config")
		sharedlogging.GetLogger(r.Context()).Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	c, err := h.store.InsertOneConfig(r.Context(), cfg)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("POST %s: inserted id %s", PathConfigs, c.ID)
		resp := sharedapi.BaseResponse[webhooks.Config]{
			Data: &c,
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("POST %s: %s", PathConfigs, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
