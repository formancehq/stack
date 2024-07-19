package cmd

import (
	"net/http"

	"go.temporal.io/sdk/worker"

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
		temporalworker.NewWorkerModule(viper.GetString(temporalTaskQueueFlag), worker.Options{
			TaskQueueActivitiesPerSecond: viper.GetFloat64(temporalMaxParallelActivities),
		}),
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
		RunE: func(cmd *cobra.Command, args []string) error {
			commonOptions, err := commonOptions(cmd)
			if err != nil {
				return err
			}

			return service.New(cmd.OutOrStdout(), commonOptions, workerOptions()).Run(cmd.Context())
		},
	}
}
