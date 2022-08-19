package mux

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/constants"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewServer(lc fx.Lifecycle) *http.ServeMux {
	return newHttpServeMux(lc, viper.GetString(constants.HttpBindAddressServerFlag))
}

func NewWorker(lc fx.Lifecycle) *http.ServeMux {
	return newHttpServeMux(lc, viper.GetString(constants.HttpBindAddressWorkerFlag))
}

func newHttpServeMux(lc fx.Lifecycle, addr string) *http.ServeMux {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sharedlogging.Infof(fmt.Sprintf("starting HTTP listening on %s", addr))
			go func() {
				if err := server.ListenAndServe(); err != nil &&
					!errors.Is(err, http.ErrServerClosed) {
					sharedlogging.GetLogger(ctx).Errorf("http.Server.ListenAndServe: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Infof("stopping HTTP listening")
			return server.Shutdown(ctx)
		},
	})

	return mux
}
