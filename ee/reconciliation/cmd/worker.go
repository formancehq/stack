package cmd

import (
	"github.com/formancehq/reconciliation/internal/v2/worker"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func workerOptions() fx.Option {
	return fx.Options(
		worker.Module(viper.GetStringSlice(topicsFlag)),
		publish.CLIPublisherModule("reconciliation"),
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
