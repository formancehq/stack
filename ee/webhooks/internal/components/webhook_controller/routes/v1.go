package routes

import (
	s "github.com/formancehq/webhooks/internal/services/httpserver"
)

var V1GetHooks = s.NewRoute(s.GET, "/configs")
var V1CreateHook = s.NewRoute(s.POST, "/configs")
var V1DeleteHook = s.NewRoute(s.DELETE, "/configs/{id}")
var V1TestHook = s.NewRoute(s.GET, "/configs/{id}/test")
var V1ActiveHook = s.NewRoute(s.PUT, "/configs/{id}/activate")
var V1DeactiveHook = s.NewRoute(s.PUT, "/configs/{id}/deactivate")
var V1ChangeSecret = s.NewRoute(s.PUT, "/configs/{id}/secret")
