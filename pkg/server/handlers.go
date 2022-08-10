package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/storage"
	"github.com/numary/webhooks/internal/svix"
	"github.com/numary/webhooks/pkg/model"
	svixgo "github.com/svix/svix-webhooks/go"
)

const (
	HealthCheckPath = "/_healthcheck"
	ConfigsPath     = "/configs"
	TogglePath      = "/toggle"
	IdPath          = "/:" + IdPathParam
	IdPathParam     = "id"
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
	h.Router.POST(ConfigsPath, h.insertOneConfigHandle)
	h.Router.DELETE(ConfigsPath+IdPath, h.deleteOneConfigHandle)
	h.Router.POST(ConfigsPath+TogglePath+IdPath, h.toggleOneConfigHandle)

	return h
}

func (h *webhooksHandler) healthCheckHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}

func (h *webhooksHandler) getAllConfigsHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cursor, err := h.store.FindAllConfigs(r.Context())
	if err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("storage.Store.FindAllConfigs: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := sharedapi.BaseResponse[model.ConfigInserted]{
		Cursor: &cursor,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.GetLogger(r.Context()).Infof("GET /configs: %d results", len(cursor.Data))
}

func (h *webhooksHandler) insertOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cfg := model.Config{}
	if err := decodeJSONBody(r, &cfg); err != nil {
		var errIB *errInvalidBody
		if errors.As(err, &errIB) {
			http.Error(w, errIB.Error(), errIB.status)
		} else {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		sharedlogging.GetLogger(r.Context()).Errorf("decodeJSONBody: %s", err)
		return
	}

	if err := cfg.Validate(); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("invalid config: %s", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var err error
	var id string
	if id, err = h.store.InsertOneConfig(r.Context(), cfg); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("storage.Store.InsertOneConfig: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := svix.CreateEndpoint(id, cfg, h.svixClient, h.svixAppId); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("svix.CreateEndpoint: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(id); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("json.Encoder.Encode: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.GetLogger(r.Context()).Infof("POST /configs: inserted %s", id)
}

func (h *webhooksHandler) deleteOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var deletedCount int64
	var err error
	if deletedCount, err = h.store.DeleteOneConfig(r.Context(), p.ByName(IdPathParam)); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("DELETE %s/%s: %s", ConfigsPath, p.ByName(IdPathParam), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if deletedCount != 1 {
		sharedlogging.GetLogger(r.Context()).Errorf("DELETE %s/%s: not found", ConfigsPath, p.ByName(IdPathParam))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := svix.DeleteEndpoint(p.ByName(IdPathParam), h.svixClient, h.svixAppId); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("svix.DeleteEndpoint: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s", ConfigsPath, p.ByName(IdPathParam))
}

func (h *webhooksHandler) toggleOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var updatedCfg model.ConfigInserted
	var modifiedCount int64
	var err error
	if updatedCfg, modifiedCount, err = h.store.ToggleOneConfig(r.Context(), p.ByName(IdPathParam)); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("POST %s%s/%s: %s", ConfigsPath, TogglePath, p.ByName(IdPathParam), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if modifiedCount != 1 {
		sharedlogging.GetLogger(r.Context()).Errorf("POST %s%s/%s: not found", ConfigsPath, TogglePath, p.ByName(IdPathParam))
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := svix.ToggleEndpoint(p.ByName(IdPathParam), updatedCfg, h.svixClient, h.svixAppId); err != nil {
		sharedlogging.GetLogger(r.Context()).Errorf("svix.ToggleEndpoint: %s", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	sharedlogging.GetLogger(r.Context()).Infof("POST %s%s/%s", ConfigsPath, TogglePath, p.ByName(IdPathParam))
}
