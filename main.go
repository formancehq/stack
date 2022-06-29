package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"
)

var Version = "v0.0"

func main() {
	fmt.Println("version:", Version)

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	router := mux.NewRouter()
	router.Handle(healthCheckPath,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Infof("health check OK")
			w.WriteHeader(http.StatusOK)
		}),
	)

	logger.Infof("starting http server on address: %s", defaultBind)
	if err := http.ListenAndServe(defaultBind, router); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
