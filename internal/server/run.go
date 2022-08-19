package server

import (
	"context"
	"net/http"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/mux"
	"github.com/numary/webhooks/internal/storage/mongo"
	"github.com/numary/webhooks/internal/svix"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func Run(cmd *cobra.Command, args []string) error {
	app := fx.New(StartModule(cmd.Context(), http.DefaultClient))

	if err := app.Start(cmd.Context()); err != nil {
		return err
	}

	<-app.Done()

	if err := app.Stop(cmd.Context()); err != nil {
		return err
	}

	return nil
}

func StartModule(ctx context.Context, httpClient *http.Client) fx.Option {
	sharedlogging.GetLogger(ctx).Debugf(
		"starting webhooks server module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	return fx.Module("webhooks server",
		fx.Provide(
			func() *http.Client { return httpClient },
			mongo.NewConfigStore,
			svix.New,
			newServerHandler,
			mux.NewServer,
		),
		fx.Invoke(registerHandler),
	)
}

func registerHandler(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}
