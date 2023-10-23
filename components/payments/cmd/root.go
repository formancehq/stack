//nolint:gochecknoglobals,golint,revive // allow for cobra & logrus init
package cmd

import (
	"fmt"
	"os"

	_ "github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/payments/cmd/api"
	"github.com/formancehq/payments/cmd/connectors"
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
		PreRunE: func(cmd *cobra.Command, args []string) error {
			if viper.GetBool(autoMigrateFlag) {
				return runMigrate(cmd, []string{"up"})
			}
			return nil
		},
	}

	version := newVersion()
	root.AddCommand(version)

	migrate := newMigrate()
	root.AddCommand(migrate)

	api := api.NewAPI(Version)
	root.AddCommand(api)

	connectors := connectors.NewConnectors(Version)
	root.AddCommand(connectors)

	migrate.Flags().String(postgresURIFlag, "postgres://localhost/payments", "PostgreSQL DB address")
	migrate.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	root.Flags().Bool(autoMigrateFlag, false, "Auto migrate database")

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
