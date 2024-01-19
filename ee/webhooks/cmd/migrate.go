package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/service"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCommand() *cobra.Command {
	migrate := &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE:  runMigrate,
	}
	return migrate
}

func runMigrate(cmd *cobra.Command, args []string) error {

	connectionOptions, err := bunconnect.ConnectionOptionsFromFlags(viper.GetViper(), cmd.OutOrStdout(), viper.GetBool(service.DebugFlag))
	if err != nil {
		return err
	}

	db, err := bunconnect.OpenSQLDB(*connectionOptions)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	return storage.Migrate(cmd.Context(), db)
}

func handleAutoMigrate(cmd *cobra.Command, _ []string) error {
	if viper.GetBool(flag.AutoMigrate) {
		return runMigrate(cmd, nil)
	}
	return nil
}
