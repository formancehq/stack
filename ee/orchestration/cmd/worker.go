package cmd

import (
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.temporal.io/sdk/worker"
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
		temporal.NewWorkerModule(viper.GetString(temporal.TemporalTaskQueueFlag), worker.Options{
			TaskQueueActivitiesPerSecond: viper.GetFloat64(temporal.TemporalMaxParallelActivities),
		}),
		triggers.NewListenerModule(
			viper.GetString(stackFlag),
			viper.GetString(temporal.TemporalTaskQueueFlag),
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
