package server

import (
	"errors"
	"net/http"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedlogging"
)

func (h *serverHandler) changeSecretHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

	err := changeOneConfigSecret(r.Context(), p.ByName(PathParamId), sec.Secret, h.store)
	if err == nil {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, p.ByName(PathParamId), PathChangeSecret)
	} else if errors.Is(err, ErrConfigNotFound) {
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathChangeSecret, ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	} else {
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathChangeSecret, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
