package server

import (
	"net/http"
	"os"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/spf13/cobra"
)

const (
	defaultBind = ":8080"

	healthCheckPath = "/_healthcheck"
)

var (
	Version = "dev"
)

func Start(cmd *cobra.Command, args []string) error {
	sharedlogging.Infof("app started with version: %s", Version)
	sharedlogging.Infof("env: %+v", syscall.Environ())

	sharedlogging.Infof("starting http server on address: %s", defaultBind)
	if err := http.ListenAndServe(defaultBind, newRouter()); err != nil {
		sharedlogging.Errorf("http.ListenAndServe: %s", err)
		os.Exit(1)
	}

	return nil
}
