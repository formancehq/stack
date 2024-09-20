package cmd

import (
	"context"

	"github.com/go-chi/chi/v5"

	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/aws/iam"
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/licence"
	"github.com/formancehq/go-libs/publish"

	"github.com/formancehq/go-libs/health"
	"github.com/formancehq/go-libs/httpserver"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/orchestration/internal/api"
	v1 "github.com/formancehq/orchestration/internal/api/v1"
	v2 "github.com/formancehq/orchestration/internal/api/v2"
	"github.com/formancehq/orchestration/internal/storage"
	"github.com/spf13/cobra"
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
		RunE: func(cmd *cobra.Command, args []string) error {
			commonOptions, err := commonOptions(cmd)
			if err != nil {
				return err
			}

			listen, _ := cmd.Flags().GetString(listenFlag)

			options := []fx.Option{
				commonOptions,
				healthCheckModule(),
				fx.Provide(func() api.ServiceInfo {
					return api.ServiceInfo{
						Version: Version,
					}
				}),
				v1.NewModule(),
				v2.NewModule(),
				fx.Invoke(func(lifecycle fx.Lifecycle, db *bun.DB) {
					lifecycle.Append(fx.Hook{
						OnStart: func(ctx context.Context) error {
							return storage.Migrate(ctx, db)
						},
					})
				}),
				api.NewModule(service.IsDebug(cmd)),
				fx.Invoke(func(lc fx.Lifecycle, router *chi.Mux) {
					lc.Append(httpserver.NewHook(router, httpserver.WithAddress(listen)))
				}),
			}
			worker, _ := cmd.Flags().GetBool(workerFlag)
			if worker {
				options = append(options, workerOptions(cmd))
			}

			return service.New(cmd.OutOrStdout(), options...).Run(cmd)
		},
	}

	cmd.Flags().Bool(workerFlag, false, "Enable worker mode")
	cmd.Flags().String(listenFlag, ":8080", "Listening address")
	cmd.Flags().Float64(temporalMaxParallelActivitiesFlag, 10, "Maximum number of parallel activities")
	cmd.Flags().String(stackURLFlag, "", "Stack url")
	cmd.Flags().String(stackClientIDFlag, "", "Stack client ID")
	cmd.Flags().String(stackClientSecretFlag, "", "Stack client secret")
	cmd.Flags().String(temporalAddressFlag, "", "Temporal server address")
	cmd.Flags().String(temporalNamespaceFlag, "default", "Temporal namespace")
	cmd.Flags().String(temporalSSLClientKeyFlag, "", "Temporal client key")
	cmd.Flags().String(temporalSSLClientCertFlag, "", "Temporal client cert")
	cmd.Flags().String(temporalTaskQueueFlag, "default", "Temporal task queue name")
	cmd.Flags().Bool(temporalInitSearchAttributes, false, "Init temporal search attributes")
	cmd.Flags().StringSlice(topicsFlag, []string{}, "Topics to listen")
	cmd.Flags().String(stackFlag, "", "Stack")

	service.AddFlags(cmd.Flags())
	publish.AddFlags(ServiceName, cmd.Flags())
	auth.AddFlags(cmd.Flags())
	bunconnect.AddFlags(cmd.Flags())
	iam.AddFlags(cmd.Flags())
	licence.AddFlags(cmd.Flags())

	return cmd
}
