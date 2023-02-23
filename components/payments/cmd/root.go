//nolint:gochecknoglobals,golint,revive // allow for cobra & logrus init
package cmd

import (
	"fmt"
	"os"

	_ "github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/stack/libs/go-libs/otlp/otlptraces"
	"github.com/formancehq/stack/libs/go-libs/publish"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	viper.SetDefault("version", Version)

	root := &cobra.Command{
		Use:               "payments",
		Short:             "payments",
		DisableAutoGenTag: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return bindFlagsToViper(cmd)
		},
	}

	version := newVersion()
	root.AddCommand(version)

	server := newServer()
	root.AddCommand(server)

	migrate := newMigrate()
	root.AddCommand(migrate)

	root.PersistentFlags().Bool(service.DebugFlag, false, "Debug mode")

	migrate.Flags().String(postgresURIFlag, "postgres://localhost/payments", "PostgreSQL DB address")
	migrate.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")

	server.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	server.Flags().String(postgresURIFlag, "postgres://localhost/payments", "PostgreSQL DB address")
	server.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	server.Flags().String(envFlag, "local", "Environment")
	server.Flags().Bool(authBasicEnabledFlag, false, "Enable basic auth")
	server.Flags().StringSlice(authBasicCredentialsFlag, []string{},
		"HTTP basic auth credentials (<username>:<password>)")
	server.Flags().Bool(authBearerEnabledFlag, false, "Enable bearer auth")
	server.Flags().String(authBearerIntrospectURLFlag, "", "OAuth2 introspect URL")
	server.Flags().StringSlice(authBearerAudienceFlag, []string{}, "Allowed audiences")
	server.Flags().Bool(authBearerAudiencesWildcardFlag, false, "Don't check audience")
	server.Flags().Bool(authBearerUseScopesFlag,
		false, "Use scopes as defined by rfc https://datatracker.ietf.org/doc/html/rfc8693")
	server.Flags().String(listenFlag, ":8080", "Listen address")
	server.Flags().Bool(autoMigrateFlag, false, "Auto migrate database")

	otlptraces.InitOTLPTracesFlags(server.Flags())
	publish.InitCLIFlags(server)

	return root
}

func Execute() {
	if err := NewRootCommand().Execute(); err != nil {
		if _, err = fmt.Fprintln(os.Stderr, err); err != nil {
			panic(err)
		}

		os.Exit(1)
	}
}
