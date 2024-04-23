package server

import (
	"net/http"
)

func (h *serverHandler) HealthCheckHandle(_ http.ResponseWriter, _ *http.Request) {}
