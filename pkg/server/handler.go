package server

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/numary/webhooks/pkg/engine"
	"github.com/numary/webhooks/pkg/engine/svix"
	"github.com/numary/webhooks/pkg/storage"
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

	store  storage.Store
	engine engine.Engine
}

func newServerHandler(store storage.Store, engine svix.Engine) http.Handler {
	h := &serverHandler{
		Router: httprouter.New(),
		store:  store,
		engine: engine,
	}

	h.Router.GET(PathHealthCheck, h.healthCheckHandle)
	h.Router.GET(PathConfigs, h.getManyConfigsHandle)
	h.Router.POST(PathConfigs, h.insertOneConfigHandle)
	h.Router.DELETE(PathConfigs+PathId, h.deleteOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathActivate, h.activateOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathDeactivate, h.deactivateOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathRotateSecret, h.rotateSecretHandle)

	return h
}
