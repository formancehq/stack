package main

import (
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"
)

var Version = "v0.0"

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Infof("app started with version: %s", Version)
	logger.Infof("env: %+v", syscall.Environ())

	router := mux.NewRouter()
	router.HandleFunc(healthCheckPath, healthCheckHandler)

	logger.Infof("starting http server on address: %s", defaultBind)
	if err := http.ListenAndServe(defaultBind, router); err != nil {
		logger.Errorf("http: %s", err)
		os.Exit(1)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("health check OK")
	w.WriteHeader(http.StatusOK)
}
