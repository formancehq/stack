package api

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/api"
)

type ServiceInfo struct {
	Version string `json:"version"`
}

func getInfo(info ServiceInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		api.RawOk(w, info)
	}
}
