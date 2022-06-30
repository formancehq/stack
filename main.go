package main

import (
	"io"
	"net/http"
	"os"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

const (
	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"
	ingestPath      = "/_ingest"
)

var Version = "v0.0"

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Infof("app started with version: %s", Version)
	logger.Infof("env: %+v", syscall.Environ())

	router := mux.NewRouter()
	router.HandleFunc(healthCheckPath, HealthCheckHandler)
	router.HandleFunc(ingestPath, IngestHandler).Methods(http.MethodPost)

	logger.Infof("starting http server on address: %s", defaultBind)
	if err := http.ListenAndServe(defaultBind, router); err != nil {
		logger.Errorf("http: %s", err)
		os.Exit(1)
	}
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	logger.Infof("health check OK")
	w.WriteHeader(http.StatusOK)
}

func IngestHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Infof("OK:%s", string(b))
	w.WriteHeader(http.StatusOK)
}
