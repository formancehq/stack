package cmd

import (
	"context"

	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"

	storage "github.com/formancehq/reconciliation/internal/storage/migrations"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"
)

func newMigrate() *cobra.Command {
	return bunmigrate.NewDefaultCommand(func(cmd *cobra.Command, args []string, db *bun.DB) error {
		return Migrate(cmd.Context(), db)
	})
}

func Migrate(ctx context.Context, db *bun.DB) error {
	return storage.Migrate(ctx, db)
}
