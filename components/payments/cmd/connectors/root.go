package connectors

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/otlp"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlpmetrics"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
)

func NewConnectors(
	version string,
	addAutoMigrateCommandFunc func(cmd *cobra.Command),
) *cobra.Command {

	root := &cobra.Command{
		Use:               "connectors",
		Short:             "connectors",
		DisableAutoGenTag: true,
	}

	cobra.EnableTraverseRunHooks = true

	server := newServer(version)
	addAutoMigrateCommandFunc(server)
	root.AddCommand(server)

	server.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	server.Flags().String(postgresURIFlag, "postgres://localhost/payments", "PostgreSQL DB address")
	server.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	server.Flags().String(envFlag, "local", "Environment")
	server.Flags().String(listenFlag, ":8080", "Listen address")

	service.AddFlags(server.Flags())
	otlp.AddFlags(server.Flags())
	otlptraces.AddFlags(server.Flags())
	otlpmetrics.AddFlags(server.Flags())
	publish.AddFlags(serviceName, server.Flags())
	iam.AddFlags(server.Flags())
	auth.AddFlags(server.Flags())

	return root
}
