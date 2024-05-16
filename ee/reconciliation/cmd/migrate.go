package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"

	storage "github.com/formancehq/reconciliation/internal/storage/migrations"
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
	return storage.Migrate(cmd.Context(), db)
}
