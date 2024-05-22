package router

import (
	"github.com/formancehq/webhooks/internal/app/webhook_server/api/handler"
	"github.com/go-chi/chi/v5"
)

var V2GetHooks = NewRoute(GET, "/v2/hooks")
var V2CreateHook = NewRoute(POST, "/v2/hooks")
var V2GetHook = NewRoute(GET, "/v2/hooks/{id}")
var V2DeleteHook = NewRoute(DELETE, "/v2/hooks/{id}")
var V2TestHook = NewRoute(POST, "/v2/hooks/{id}/test")
var V2ActiveHook = NewRoute(PUT, "/v2/hooks/{id}/activate")
var V2DeactiveHook = NewRoute(PUT, "/v2/hooks/{id}/deactivate")
var V2ChangeHookSecret = NewRoute(PUT, "/v2/hooks/{id}/secret")
var V2ChangeHookEndpoint = NewRoute(PUT, "/v2/hooks/{id}/endpoint")
var V2ChangeHookRetry = NewRoute(PUT, "/v2/hooks/{id}/retry")

var V2GetWaitingAttempts = NewRoute(GET, "/v2/attempts/waiting")
var V2GetAbortedAttempts = NewRoute(GET, "/v2/attempts/aborted")

var V2RetryWaitingAttempts = NewRoute(PUT, "/v2/attempts/waiting/flush")
var V2RetryWaitingAttempt = NewRoute(PUT, "/v2/attempts/waiting/{id}/flush")

var V2AbortWaitingAttempt = NewRoute(PUT, "/v2/attempts/waiting/{id}/abort")

func NewRouterV2() chi.Router {
	mux := chi.NewRouter()

	mux.Get(V2GetHooks.Url, handler.V2GetHooks)
	mux.Get(V2GetHook.Url, handler.V2GetHook)
	mux.Post(V2CreateHook.Url, handler.V2CreateHook)
	mux.Delete(V2DeleteHook.Url, handler.V2DeleteHook)
	mux.Post(V2TestHook.Url, handler.V2TestHook)
	mux.Put(V2ActiveHook.Url, handler.V2ActivateHook)
	mux.Put(V2DeactiveHook.Url, handler.V2DeactivateHook)
	mux.Put(V2ChangeHookSecret.Url, handler.V2ChangeHookSecret)
	mux.Put(V2ChangeHookEndpoint.Url, handler.V2ChangeHookEndpoint)
	mux.Put(V2ChangeHookRetry.Url, handler.V2ChangeHookRetry)

	mux.Get(V2GetWaitingAttempts.Url, handler.V2GetWaitingAttempts)
	mux.Get(V2GetAbortedAttempts.Url, handler.V2GetAbortedAttempts)
	mux.Put(V2RetryWaitingAttempts.Url, handler.V2RetryWaitingAttempts)
	mux.Put(V2RetryWaitingAttempt.Url, handler.V2RetryWaitingAttempt)
	mux.Put(V2AbortWaitingAttempt.Url, handler.V2AbortWaitingAttempt)

	return mux
}
