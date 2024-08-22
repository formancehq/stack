package api

import (
	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/aws/iam"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
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
	}

	cobra.EnableTraverseRunHooks = true

	server := newServer(version)
	addAutoMigrateCommandFunc(server)
	root.AddCommand(server)

	server.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	server.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	server.Flags().String(envFlag, "local", "Environment")
	server.Flags().String(listenFlag, ":8080", "Listen address")

	service.AddFlags(server.Flags())
	otlp.AddFlags(server.Flags())
	otlptraces.AddFlags(server.Flags())
	otlpmetrics.AddFlags(server.Flags())
	auth.AddFlags(server.Flags())
	publish.AddFlags(serviceName, server.Flags())
	bunconnect.AddFlags(server.Flags())
	iam.AddFlags(server.Flags())

	return root
}
