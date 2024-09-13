package cmd

import (
	"github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
)

func newServer() *cobra.Command {
	return &cobra.Command{
		Use:          "serve",
		Aliases:      []string{"server"},
		Short:        "Launch api server",
		SilenceUsage: true,
		RunE:         runServer(),
	}
}

func runServer() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		setLogger()

		options, err := commonOptions(cmd)
		if err != nil {
			return err
		}

		return service.New(cmd.OutOrStdout(), options).Run(cmd.Context())
	}
}

func setLogger() {
	// Add a dedicated logger for opentelemetry in case of error
	otel.SetLogger(logrusr.New(logrus.New().WithField("component", "otlp")))
}
