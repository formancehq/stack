package httpserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/numary/go-libs/sharedlogging"
	"go.uber.org/fx"
)

func NewMuxServer(addr string) (*http.ServeMux, *http.Server) {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	return mux, server
}

func RegisterHandler(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}

func Run(lc fx.Lifecycle, server *http.Server, addr string) {
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
}
