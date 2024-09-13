package cmd

import (
	"net/http"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/orchestration/internal/triggers"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/licence"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/formancehq/stack/libs/go-libs/temporal"
	"github.com/spf13/cobra"
	"go.temporal.io/sdk/worker"
	"go.uber.org/fx"
)

func stackClientModule(cmd *cobra.Command) fx.Option {
	stackURL, _ := cmd.Flags().GetString(stackURLFlag)

	return fx.Options(
		fx.Provide(func(httpClient *http.Client) *sdk.Formance {
			return sdk.New(
				sdk.WithClient(httpClient),
				sdk.WithServerURL(stackURL),
			)
		}),
	)
}

func workerOptions(cmd *cobra.Command) fx.Option {
	stack, _ := cmd.Flags().GetString(stackFlag)
	temporalTaskQueue, _ := cmd.Flags().GetString(temporal.TemporalTaskQueueFlag)
	temporalMaxParallelActivities, _ := cmd.Flags().GetInt(temporal.TemporalMaxParallelActivitiesFlag)
	topics, _ := cmd.Flags().GetStringSlice(topicsFlag)

	return fx.Options(
		stackClientModule(cmd),
		temporal.NewWorkerModule(temporalTaskQueue, worker.Options{
			TaskQueueActivitiesPerSecond: float64(temporalMaxParallelActivities),
		}),
		triggers.NewListenerModule(
			stack,
			temporalTaskQueue,
			topics,
		),
	)
}

func newWorkerCommand() *cobra.Command {
	ret := &cobra.Command{
		Use: "worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			commonOptions, err := commonOptions(cmd)
			if err != nil {
				return err
			}

			return service.New(cmd.OutOrStdout(), commonOptions, workerOptions(cmd)).Run(cmd)
		},
	}
	ret.Flags().String(stackURLFlag, "", "Stack url")
	ret.Flags().String(stackClientIDFlag, "", "Stack client ID")
	ret.Flags().String(stackClientSecretFlag, "", "Stack client secret")
	ret.Flags().StringSlice(topicsFlag, []string{}, "Topics to listen")
	ret.Flags().String(stackFlag, "", "Stack")

	temporal.AddFlags(ret.Flags())
	publish.AddFlags(ServiceName, ret.Flags())
	bunconnect.AddFlags(ret.Flags())
	iam.AddFlags(ret.Flags())
	service.AddFlags(ret.Flags())
	licence.AddFlags(ret.Flags())

	return ret
}
