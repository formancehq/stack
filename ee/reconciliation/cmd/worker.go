package cmd

import (
	"time"

	"github.com/formancehq/reconciliation/internal/v2/worker"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const (
	transactionBasedDelayFlag = "transaction-based-delay"
)

func workerOptions() fx.Option {
	return fx.Options(
		worker.Module(
			viper.GetStringSlice(topicsFlag),
			viper.GetDuration(transactionBasedDelayFlag),
		),
		publish.CLIPublisherModule("reconciliation"),
	)
}

func newWorkerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			commonOptions, err := commonOptions(cmd)
			if err != nil {
				return err
			}

			return service.New(cmd.OutOrStdout(), commonOptions, workerOptions()).Run(cmd.Context())
		},
	}

	cmd.Flags().Duration(transactionBasedDelayFlag, 15*24*time.Hour, "Delay after which a pending transactions/payment is marked as failed")

	return cmd
}
