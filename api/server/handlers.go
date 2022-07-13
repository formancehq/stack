package server

import (
	"net/http"

	"github.com/numary/go-libs/sharedlogging"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	sharedlogging.Infof("health check OK")
	_, err := w.Write([]byte("ok"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
