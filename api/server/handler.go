package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/internal/storage"
	"github.com/numary/webhooks-cloud/pkg/model"
)

const (
	HealthCheckPath = "/_healthcheck"
	ConfigsPath     = "/configs"
)

type webhooksHandler struct {
	*httprouter.Router

	store storage.Store
}

func newWebhooksHandler(store storage.Store) http.Handler {
	h := &webhooksHandler{
		Router: httprouter.New(),
		store:  store,
	}

	h.Router.GET(HealthCheckPath, h.healthCheckHandle)
	h.Router.GET(ConfigsPath, h.getAllConfigsHandle)
	h.Router.POST(ConfigsPath, h.insertConfigHandle)
	h.Router.DELETE(ConfigsPath, h.deleteAllConfigsHandle)

	return h
}

func (h *webhooksHandler) getAllConfigsHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cursor, err := h.store.FindAllConfigs(r.Context())
	if err != nil {
		sharedlogging.Errorf("mongodb.Store.FindAllConfigs: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := sharedapi.BaseResponse[model.ConfigInserted]{
		Cursor: &cursor,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		sharedlogging.Errorf("json.Encoder.Encode: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.Infof("GET /configs: %d results", len(cursor.Data))
}

func (h *webhooksHandler) insertConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	config := model.Config{}
	if err := decodeJSONBody(r, &config); err != nil {
		var errIB *errInvalidBody
		if errors.As(err, &errIB) {
			http.Error(w, errIB.Error(), errIB.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		sharedlogging.Errorf("decodeJSONBody: %s", err)
		return
	}

	if err := config.Validate(); err != nil {
		sharedlogging.Errorf("invalid config: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	var id string
	if id, err = h.store.InsertOneConfig(r.Context(), config); err != nil {
		sharedlogging.Errorf("mongodb.Store.InsertOneConfig: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.Infof("POST /configs: inserted %s", id)
}

func (h *webhooksHandler) deleteAllConfigsHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := h.store.DropConfigsCollection(r.Context()); err != nil {
		sharedlogging.Errorf("mongodb.Store.DropConfigsCollection: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.Infof("deleted all configs")
}

func (h *webhooksHandler) healthCheckHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sharedlogging.Infof("health check OK")
}
