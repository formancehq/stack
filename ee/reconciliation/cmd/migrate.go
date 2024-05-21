package cmd

import (
	migrationsv1 "github.com/formancehq/reconciliation/internal/v1/storage/migrations"
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

var (
	autoMigrateFlag = "auto-migrate"
)

func newMigrate() *cobra.Command {
	return bunmigrate.NewDefaultCommand(Migrate)
}

func Migrate(cmd *cobra.Command, args []string, db *bun.DB) error {
	migrator := migrations.NewMigrator()
	migrationsv1.RegisterMigrations(migrator)
	return migrator.Up(cmd.Context(), db)

}
