package cmd

import (
	"database/sql"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
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
	dsn := viper.GetString(flag.StoragePostgresConnString)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	defer sqldb.Close()
	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()
	if viper.GetBool(flag.Debug) {
		db.AddQueryHook(bundebug.NewQueryHook(bundebug.WithWriter(cmd.OutOrStdout())))
	}
	return storage.Migrate(cmd.Context(), db)
}

func handleAutoMigrate(cmd *cobra.Command, _ []string) error {
	if viper.GetBool(flag.AutoMigrate) {
		return runMigrate(cmd, nil)
	}
	return nil
}
