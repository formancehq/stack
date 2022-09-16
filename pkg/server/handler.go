package server

import (
	"net/http"

	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/julienschmidt/httprouter"
)

const (
	PathHealthCheck  = "/_healthcheck"
	PathConfigs      = "/configs"
	PathActivate     = "/activate"
	PathDeactivate   = "/deactivate"
	PathChangeSecret = "/secret/change"
	PathId           = "/:" + PathParamId
	PathParamId      = "id"
)

type serverHandler struct {
	*httprouter.Router

	store storage.Store
}

func newServerHandler(store storage.Store) http.Handler {
	h := &serverHandler{
		Router: httprouter.New(),
		store:  store,
	}

	h.Router.GET(PathHealthCheck, h.HealthCheckHandle)
	h.Router.GET(PathConfigs, h.getManyConfigsHandle)
	h.Router.POST(PathConfigs, h.insertOneConfigHandle)
	h.Router.DELETE(PathConfigs+PathId, h.deleteOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathActivate, h.activateOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathDeactivate, h.deactivateOneConfigHandle)
	h.Router.PUT(PathConfigs+PathId+PathChangeSecret, h.changeSecretHandle)

	return h
}
