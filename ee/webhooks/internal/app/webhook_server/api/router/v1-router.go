package router

import (
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/handler"
	"github.com/go-chi/chi/v5"
)

var V1GetHooks = NewRoute(GET, "/configs")
var V1CreateHook = NewRoute(POST, "/configs")
var V1DeleteHook = NewRoute(DELETE, "/configs/{id}")
var V1TestHook = NewRoute(GET, "/configs/{id}/test")
var V1ActiveHook = NewRoute(PUT, "/configs/{id}/activate")
var V1DeactiveHook = NewRoute(PUT, "/configs/{id}/deactivate")
var V1ChangeSecret = NewRoute(PUT, "/configs/{id}/secret/change")

func NewRouterV1() chi.Router {

	mux := chi.NewRouter()

	mux.Get(V1GetHooks.Url, handler.V1GetHooks)
	mux.Post(V1CreateHook.Url, handler.V1CreateHook)
	mux.Delete(V1DeleteHook.Url, handler.DeleteHook)
	mux.Get(V1TestHook.Url, handler.V1TestHook)
	mux.Put(V1ActiveHook.Url, handler.V1ActivateHook)
	mux.Put(V1DeactiveHook.Url, handler.V1DeactivateHook)
	mux.Put(V1ChangeSecret.Url, handler.V1ChangeHookSecret)

	return mux
}
