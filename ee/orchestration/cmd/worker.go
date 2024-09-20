package cmd

import (
	"net/http"

	"github.com/formancehq/go-libs/aws/iam"
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/licence"
	"github.com/formancehq/go-libs/publish"

	"go.temporal.io/sdk/worker"

	"github.com/formancehq/orchestration/internal/triggers"

	sdk "github.com/formancehq/formance-sdk-go/v2"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/orchestration/internal/temporalworker"
	"github.com/spf13/cobra"
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
	temporalTaskQueue, _ := cmd.Flags().GetString(temporalTaskQueueFlag)
	temporalMaxParallelActivities, _ := cmd.Flags().GetInt(temporalMaxParallelActivitiesFlag)
	topics, _ := cmd.Flags().GetStringSlice(topicsFlag)

	return fx.Options(
		stackClientModule(cmd),
		temporalworker.NewWorkerModule(temporalTaskQueue, worker.Options{
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
	ret.Flags().Float64(temporalMaxParallelActivitiesFlag, 10, "Maximum number of parallel activities")
	ret.Flags().String(stackURLFlag, "", "Stack url")
	ret.Flags().String(stackClientIDFlag, "", "Stack client ID")
	ret.Flags().String(stackClientSecretFlag, "", "Stack client secret")
	ret.Flags().String(temporalAddressFlag, "", "Temporal server address")
	ret.Flags().String(temporalNamespaceFlag, "default", "Temporal namespace")
	ret.Flags().String(temporalSSLClientKeyFlag, "", "Temporal client key")
	ret.Flags().String(temporalSSLClientCertFlag, "", "Temporal client cert")
	ret.Flags().String(temporalTaskQueueFlag, "default", "Temporal task queue name")
	ret.Flags().Bool(temporalInitSearchAttributes, false, "Init temporal search attributes")
	ret.Flags().StringSlice(topicsFlag, []string{}, "Topics to listen")
	ret.Flags().String(stackFlag, "", "Stack")

	publish.AddFlags(ServiceName, ret.Flags())
	bunconnect.AddFlags(ret.Flags())
	iam.AddFlags(ret.Flags())
	service.AddFlags(ret.Flags())
	licence.AddFlags(ret.Flags())

	return ret
}
