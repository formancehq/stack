package api

import (
	"net/http"

	sharedapi "github.com/formancehq/stack/libs/go-libs/api"
)

func listTriggers(backend Backend) func(writer http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		triggers, err := backend.ListTriggers(r.Context())
		if err != nil {
			sharedapi.InternalServerError(w, r, err)
			return
		}

		sharedapi.Ok(w, triggers)
	}
}
