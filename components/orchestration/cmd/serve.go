package cmd

import (
	"context"

	"github.com/formancehq/orchestration/internal/api"
	"github.com/formancehq/orchestration/internal/storage"
	"github.com/formancehq/orchestration/internal/workflow"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

func healthCheckModule() fx.Option {
	return fx.Options(
		health.Module(),
		health.ProvideHealthCheck(func() health.NamedCheck {
			return health.NewNamedCheck("default", health.CheckFn(func(ctx context.Context) error {
				return nil
			}))
		}),
	)
}

func newServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "serve",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {

			options := []fx.Option{
				healthCheckModule(),
				api.NewModule(),
				workflow.NewModule(viper.GetString(temporalTaskQueueFlag)),
				fx.Invoke(func(lifecycle fx.Lifecycle, db *bun.DB) {
					lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return storage.Migrate(ctx, db)
						},
					})
				}),
				fx.Invoke(func(lc fx.Lifecycle, router *chi.Mux) {
					lc.Append(httpserver.NewHook(viper.GetString(listenFlag), router))
				}),
				commonOptions(cmd.OutOrStdout()),
			}
			if viper.GetBool(workerFlag) {
				options = append(options, workerOptions())
			}

			return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
		},
	}
	cmd.Flags().Bool(workerFlag, false, "Enable worker mode")
	cmd.Flags().String(listenFlag, ":8080", "Listening address")
	return cmd
}
