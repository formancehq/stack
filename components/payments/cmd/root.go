//nolint:gochecknoglobals,golint,revive // allow for cobra & logrus init
package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/service"

	_ "github.com/bombsimon/logrusr/v3"
	"github.com/formancehq/payments/cmd/api"
	"github.com/formancehq/payments/cmd/connectors"
	"github.com/spf13/cobra"
)

var (
	Version   = "develop"
	BuildDate = "-"
	Commit    = "-"
)

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:               "payments",
		Short:             "payments",
		DisableAutoGenTag: true,
		Version:           Version,
	}

	version := newVersion()
	root.AddCommand(version)

	migrate := newMigrate()
	root.AddCommand(migrate)

	api := api.NewAPI(Version, addAutoMigrateCommand)
	root.AddCommand(api)

	connectors := connectors.NewConnectors(Version, addAutoMigrateCommand)
	root.AddCommand(connectors)

	return root
}

func Execute() {
	service.Execute(NewRootCommand())
}

func addAutoMigrateCommand(cmd *cobra.Command) {
	cmd.Flags().Bool(autoMigrateFlag, false, "Auto migrate database")
	cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
		autoMigrate, _ := cmd.Flags().GetBool(autoMigrateFlag)
		if autoMigrate {
			return bunmigrate.Run(cmd, args, Migrate)
		}
		return nil
	}
}
