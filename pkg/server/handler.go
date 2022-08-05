package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/internal/storage"
	"github.com/numary/webhooks-cloud/internal/svix"
	"github.com/numary/webhooks-cloud/pkg/model"
	svixgo "github.com/svix/svix-webhooks/go"
)

const (
	HealthCheckPath = "/_healthcheck"
	ConfigsPath     = "/configs"
)

type webhooksHandler struct {
	*httprouter.Router

	store      storage.Store
	svixClient *svixgo.Svix
	svixAppId  string
}

func newConfigHandler(store storage.Store, svixClient *svixgo.Svix, svixAppId string) http.Handler {
	h := &webhooksHandler{
		Router:     httprouter.New(),
		store:      store,
		svixClient: svixClient,
		svixAppId:  svixAppId,
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
		sharedlogging.Errorf("storage.Store.FindAllConfigs: %s", err)
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
	cfg := model.Config{}
	if err := decodeJSONBody(r, &cfg); err != nil {
		var errIB *errInvalidBody
		if errors.As(err, &errIB) {
			http.Error(w, errIB.Error(), errIB.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		sharedlogging.Errorf("decodeJSONBody: %s", err)
		return
	}

	if err := cfg.Validate(); err != nil {
		sharedlogging.Errorf("invalid config: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	var id string
	if id, err = h.store.InsertOneConfig(r.Context(), cfg); err != nil {
		sharedlogging.Errorf("storage.Store.InsertOneConfig: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := svix.DeleteAllEndpoints(h.svixClient, h.svixAppId); err != nil {
		sharedlogging.Errorf("svix.DeleteAllEndpoints: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	for _, endpoint := range cfg.Endpoints {
		if err := svix.CreateEndpoint(h.svixClient, h.svixAppId, endpoint); err != nil {
			sharedlogging.Errorf(
				"svix.CreateEndpoint: %s", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	sharedlogging.Infof("POST /configs: inserted %s", id)
}

func (h *webhooksHandler) deleteAllConfigsHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	if err := h.store.DropConfigsCollection(r.Context()); err != nil {
		sharedlogging.Errorf("storage.Store.DropConfigsCollection: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := svix.DeleteAllEndpoints(h.svixClient, h.svixAppId); err != nil {
		sharedlogging.Errorf("svix.DeleteAllEndpoints: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.Infof("deleted all configs")
}

func (h *webhooksHandler) healthCheckHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sharedlogging.Infof("health check OK")
}
