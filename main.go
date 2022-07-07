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

	healthCheckPath   = "/_healthcheck"
	organizationsPath = "/organizations"
	trashPath         = "/trash"
)

var Version = "v0.0"

var logger = logrus.New()

func main() {
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.Infof("app started with version: %s", Version)
	logger.Infof("env: %+v", syscall.Environ())

	router := mux.NewRouter()
	router.HandleFunc(healthCheckPath, healthCheckHandler)
	router.HandleFunc(organizationsPath, deleteOrganizationHandler).Methods(http.MethodDelete)
	router.HandleFunc(trashPath, trashHandler).Methods(http.MethodPost)

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

func deleteOrganizationHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Infof("DELETE:%s", string(b))
	w.WriteHeader(http.StatusOK)
}

func trashHandler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		logger.Errorf("%s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	logger.Infof("TRASHED:%s", string(b))
	w.WriteHeader(http.StatusOK)
}
