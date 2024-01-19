package cmd

import (
	"net/http"

	"github.com/formancehq/orchestration/internal/triggers"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func stackClientModule() fx.Option {
	return fx.Options(
		fx.Provide(func(httpClient *http.Client) *sdk.Formance {
			return sdk.New(
				sdk.WithClient(httpClient),
				sdk.WithServerURL(viper.GetString(stackURLFlag)),
			)
		}),
	)
}

func workerOptions() fx.Option {
	return fx.Options(
		stackClientModule(),
		temporalworker.NewWorkerModule(viper.GetString(temporalTaskQueueFlag)),
		triggers.NewListenerModule(
			viper.GetString(stackFlag),
			viper.GetString(temporalTaskQueueFlag),
			viper.GetStringSlice(topicsFlag),
		),
	)
}

func newWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use: "worker",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			options := []fx.Option{
				commonOptions(cmd.OutOrStdout()),
				workerOptions(),
			}

			return service.New(cmd.OutOrStdout(), options...).Run(cmd.Context())
		},
	}
}
