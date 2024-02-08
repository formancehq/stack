package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"
	"github.com/uptrace/bun"

	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/pkg/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func newMigrateCommand() *cobra.Command {
	return bunmigrate.NewDefaultCommand(func(cmd *cobra.Command, args []string, db *bun.DB) error {
		return storage.Migrate(cmd.Context(), db)
	})
}

func handleAutoMigrate(cmd *cobra.Command, args []string) error {
	if viper.GetBool(flag.AutoMigrate) {
		return bunmigrate.Run(cmd, args, func(cmd *cobra.Command, args []string, db *bun.DB) error {
			return storage.Migrate(cmd.Context(), db)
		})
	}
	return nil
}
