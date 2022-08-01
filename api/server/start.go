package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/internal/storage/mongodb"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func Start(*cobra.Command, []string) error {
	sharedlogging.Infof("env: %+v", syscall.Environ())

	app := fx.New(StartModule())
	app.Run()

	return nil
}

func StartModule() fx.Option {
	return fx.Module("webhooks-module",
		fx.Provide(
			mongodb.NewStore,
			newWebhooksHandler,
			newMux,
		),
		fx.Invoke(register),
	)
}

func newMux(lc fx.Lifecycle) *http.ServeMux {
	bindAddr := viper.GetString(constants.ServerHttpBindAddressFlag)
	if bindAddr == "" {
		bindAddr = constants.DefaultBindAddress
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    bindAddr,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			sharedlogging.Infof(fmt.Sprintf("starting HTTP server on %s", bindAddr))
			go func() {
				if err := server.ListenAndServe(); err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					sharedlogging.Errorf("ListenAndServe: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sharedlogging.Infof("stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})

	return mux
}

func register(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}
