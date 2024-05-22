package routes

import(

	s "github.com/formancehq/webhooks/internal/services/httpserver"
)


var V2GetHooks = s.NewRoute(s.GET, "/v2/hooks")
var V2CreateHook = s.NewRoute(s.POST, "/v2/hooks")
var V2DeleteHook = s.NewRoute(s.DELETE, "/v2/hooks/{id}")
var V2TestHook = s.NewRoute(s.GET, "/v2/hooks/{id}/test")
var V2ActiveHook = s.NewRoute(s.PUT,"/v2/hooks/{id}/activate")
var V2DeactiveHook = s.NewRoute(s.PUT, "/v2/hooks/{id}/deactivate")
var V2ChangeHookSecret = s.NewRoute(s.PUT, "/v2/hooks/{id}/secret")
var V2ChangeHookEndpoint = s.NewRoute(s.PUT, "/v2/hooks/{id}/endpoint")
var V2ChangeHookRetry = s.NewRoute(s.PUT, "/v2/hooks/{id}/retry")


var V2GetWaitingAttempts = s.NewRoute(s.GET, "/v2/attempts/waiting")
var V2GetAbortedAttempts = s.NewRoute(s.GET, "/v2/attempts/aborted")

var V2RetryWaitingAttempts 	= s.NewRoute(s.PUT, "/v2/attempts/waiting/flush")
var V2RetryWaitingAttempt 	= s.NewRoute(s.PUT, "/v2/attempts/waiting/{id}/flush")

var V2AbortWaitingAttempt 	= s.NewRoute(s.PUT, "/v2/attempts/waiting/{id}/abort")

