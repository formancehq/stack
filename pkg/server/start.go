package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/cmd/constants"
	"github.com/numary/webhooks-cloud/internal/storage/mongo"
	"github.com/numary/webhooks-cloud/internal/svix"
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
	return fx.Module("webhooks server module",
		fx.Provide(
			mongo.NewConfigStore,
			svix.New,
			newConfigHandler,
			newHttpServeMux,
		),
		fx.Invoke(registerConfigHandler),
	)
}

func newHttpServeMux(lc fx.Lifecycle) *http.ServeMux {
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

func registerConfigHandler(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}
