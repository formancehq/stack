package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/formancehq/go-libs/sharedlogging"
	webhooks "github.com/formancehq/webhooks/pkg"
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
		err := fmt.Errorf("invalid config: %w", err)
		sharedlogging.GetLogger(r.Context()).Errorf(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if id, err := insertOneConfig(r.Context(), cfg, h.store); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("POST %s: %s", PathConfigs, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else if err := json.NewEncoder(w).Encode(id); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		sharedlogging.GetLogger(r.Context()).Infof("POST %s: inserted id %s", PathConfigs, id)
	}
}
