package cmd

import (
	"github.com/formancehq/stack/libs/go-libs/bun/bunmigrate"

	"github.com/formancehq/payments/internal/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"

	// Import the postgres driver.
	_ "github.com/lib/pq"
)

var (
	configEncryptionKeyFlag = "config-encryption-key"
	autoMigrateFlag         = "auto-migrate"
)

func newMigrate() *cobra.Command {
	return bunmigrate.NewDefaultCommand(Migrate, func(cmd *cobra.Command) {
		cmd.Flags().String(configEncryptionKeyFlag, "", "Config encryption key")
	})
}

func Migrate(cmd *cobra.Command, args []string, db *bun.DB) error {
	cfgEncryptionKey := viper.GetString(configEncryptionKeyFlag)
	if cfgEncryptionKey == "" {
		cfgEncryptionKey = cmd.Flag(configEncryptionKeyFlag).Value.String()
	}

	if cfgEncryptionKey != "" {
		storage.EncryptionKey = cfgEncryptionKey
	}

	return storage.Migrate(cmd.Context(), db)
}
