package cmd

import (
	"github.com/formancehq/go-libs/bun/bunmigrate"
	"github.com/formancehq/payments/internal/storage"
	"github.com/spf13/cobra"
	"github.com/uptrace/bun"

	// Import the postgres driver.
	_ "github.com/lib/pq"
)

var (
	autoMigrateFlag = "auto-migrate"
)

func newMigrate() *cobra.Command {
	cmd := bunmigrate.NewDefaultCommand(Migrate, func(cmd *cobra.Command) {
		cmd.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	})

	return cmd
}

func Migrate(cmd *cobra.Command, args []string, db *bun.DB) error {
	cfgEncryptionKey, _ := cmd.Flags().GetString(configEncryptionKeyFlag)
	if cfgEncryptionKey == "" {
		cfgEncryptionKey = cmd.Flag(configEncryptionKeyFlag).Value.String()
	}

	if cfgEncryptionKey != "" {
		storage.EncryptionKey = cfgEncryptionKey
	}

	return storage.Migrate(cmd.Context(), db)
}
