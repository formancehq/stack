package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/go-libs/sharedapi"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/storage"
	"github.com/numary/webhooks/pkg/model"
	"github.com/numary/webhooks/pkg/service"
	svixgo "github.com/svix/svix-webhooks/go"
)

const (
	PathHealthCheck  = "/_healthcheck"
	PathConfigs      = "/configs"
	PathActivate     = "/activate"
	PathDeactivate   = "/deactivate"
	PathRotateSecret = "/secret/rotate"
	PathId           = "/:" + PathParamId
	PathParamId      = "id"
)

type serverHandler struct {
	*httprouter.Router

	store      storage.Store
	svixClient *svixgo.Svix
	svixAppId  string
}

func newServerHandler(store storage.Store, svixClient *svixgo.Svix, svixAppId string) http.Handler {
	h := &serverHandler{
		Router:     httprouter.New(),
		store:      store,
		svixClient: svixClient,
		svixAppId:  svixAppId,
	}

	h.Router.GET(PathHealthCheck, h.healthCheckHandle)
	h.Router.GET(PathConfigs, h.getAllConfigsHandle)
	h.Router.POST(PathConfigs, h.insertOneConfigHandle)
	h.Router.DELETE(PathConfigs+PathId, h.deleteOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathActivate, h.activateOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathDeactivate, h.deactivateOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathRotateSecret, h.rotateOneConfigSecretHandle)

	return h
}

func (h *serverHandler) healthCheckHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sharedlogging.GetLogger(r.Context()).Infof("health check OK")
}

func (h *serverHandler) getAllConfigsHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
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

func (h *serverHandler) insertOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	cfg := model.Config{}
	if err := decodeJSONBody(r, &cfg, false); err != nil {
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

	if id, err := service.InsertOneConfig(cfg, r.Context(), h.store, h.svixClient, h.svixAppId); err != nil {
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

func (h *serverHandler) deleteOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch err := service.DeleteOneConfig(r.Context(), p.ByName(PathParamId), h.store, h.svixClient, h.svixAppId); err {
	case nil:
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s", PathConfigs, p.ByName(PathParamId))
	case service.ErrConfigNotFound:
		sharedlogging.GetLogger(r.Context()).Infof("DELETE %s/%s: %s", PathConfigs, p.ByName(PathParamId), service.ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	default:
		sharedlogging.GetLogger(r.Context()).Errorf("DELETE %s/%s: %s", PathConfigs, p.ByName(PathParamId), err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *serverHandler) activateOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch err := service.ActivateOneConfig(true, r.Context(), p.ByName(PathParamId), h.store, h.svixClient, h.svixAppId); err {
	case nil:
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, p.ByName(PathParamId), PathActivate)
	case service.ErrConfigNotFound:
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathActivate, service.ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	default:
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathActivate, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *serverHandler) deactivateOneConfigHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch err := service.ActivateOneConfig(false, r.Context(), p.ByName(PathParamId), h.store, h.svixClient, h.svixAppId); err {
	case nil:
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, p.ByName(PathParamId), PathDeactivate)
	case service.ErrConfigNotFound:
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathDeactivate, service.ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	default:
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathDeactivate, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (h *serverHandler) rotateOneConfigSecretHandle(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sec := model.Secret{}
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

	switch err := service.RotateOneConfigSecret(r.Context(), p.ByName(PathParamId), sec.Secret, h.store, h.svixClient, h.svixAppId); err {
	case nil:
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s", PathConfigs, p.ByName(PathParamId), PathRotateSecret)
	case service.ErrConfigNotFound:
		sharedlogging.GetLogger(r.Context()).Infof("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathRotateSecret, service.ErrConfigNotFound)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	default:
		sharedlogging.GetLogger(r.Context()).Errorf("PUT %s/%s%s: %s", PathConfigs, p.ByName(PathParamId), PathRotateSecret, err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
