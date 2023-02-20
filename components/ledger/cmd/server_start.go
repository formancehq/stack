package cmd

import (
	"context"

	app "github.com/formancehq/stack/libs/go-libs/app"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/numary/ledger/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func NewServerStart() *cobra.Command {
	return &cobra.Command{
		Use: "start",
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := app.DefaultLoggingContext(cmd, viper.GetBool(debugFlag))

			app := NewContainer(
				viper.GetViper(),
				fx.Invoke(func(lc fx.Lifecycle, h *api.API) {
					lc.Append(httpserver.NewHook(viper.GetString(serverHttpBindAddressFlag), h))
				}),
			)
			errCh := make(chan error, 1)
			go func() {
				err := app.Start(ctx)
				if err != nil {
					errCh <- err
				}
			}()
			select {
			case err := <-errCh:
				return err
			case <-ctx.Done():
				return app.Stop(context.Background())
			case <-app.Done():
				return app.Err()
			}
		},
	}
}
