package api

import (
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
)

func NewAPI(
	version string,
	addAutoMigrateCommandFunc func(cmd *cobra.Command),
) *cobra.Command {

	root := &cobra.Command{
		Use:               "api",
		Short:             "api",
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	server := newServer(version)
	addAutoMigrateCommandFunc(server)
	root.AddCommand(server)

	root.PersistentFlags().Bool(service.DebugFlag, false, "Debug mode")

	server.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	server.Flags().String(postgresURIFlag, "postgres://localhost/payments", "PostgreSQL DB address")
	server.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	server.Flags().String(envFlag, "local", "Environment")
	server.Flags().String(listenFlag, ":8080", "Listen address")

	otlp.InitOTLPFlags(server.Flags())
	otlptraces.InitOTLPTracesFlags(server.Flags())
	otlpmetrics.InitOTLPMetricsFlags(server.Flags())
	publish.InitCLIFlags(server)

	return root
}
