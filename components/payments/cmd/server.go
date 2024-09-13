package cmd

import (
	"github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/go-libs/auth"
	"github.com/formancehq/go-libs/aws/iam"
	"github.com/formancehq/go-libs/bun/bunconnect"
	"github.com/formancehq/go-libs/licence"
	"github.com/formancehq/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/go-libs/otlp/otlptraces"
	"github.com/formancehq/go-libs/publish"
	"github.com/formancehq/go-libs/service"
	"github.com/formancehq/go-libs/temporal"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

func newServer() *cobra.Command {
	cmd := &cobra.Command{
		Use:          "serve",
		Aliases:      []string{"server"},
		Short:        "Launch api server",
		SilenceUsage: true,
		RunE:         runServer(),
	}

	service.AddFlags(cmd.Flags())
	otlpmetrics.AddFlags(cmd.Flags())
	otlptraces.AddFlags(cmd.Flags())
	auth.AddFlags(cmd.Flags())
	publish.AddFlags(ServiceName, cmd.Flags())
	bunconnect.AddFlags(cmd.Flags())
	iam.AddFlags(cmd.Flags())
	temporal.AddFlags(cmd.Flags())
	licence.AddFlags(cmd.Flags())

	return cmd
}

func runServer() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		setLogger()

		options, err := commonOptions(cmd)
		if err != nil {
			return err
		}

		return service.New(cmd.OutOrStdout(), options).Run(cmd)
	}
}

func setLogger() {
	// Add a dedicated logger for opentelemetry in case of error
	otel.SetLogger(logrusr.New(logrus.New().WithField("component", "otlp")))
}
