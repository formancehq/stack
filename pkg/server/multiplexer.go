package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/cmd/constants"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

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
