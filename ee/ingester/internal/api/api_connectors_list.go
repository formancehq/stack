package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

func (a *API) listConnectors(w http.ResponseWriter, r *http.Request) {
	connectors, err := a.backend.ListConnectors(r.Context())
	if err != nil {
		api.InternalServerError(w, r, err)
		return
	}

	api.RenderCursor(w, *connectors)
}
