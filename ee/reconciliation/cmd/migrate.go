package cmd

import (
	"context"

	storage "github.com/formancehq/reconciliation/internal/storage/migrations"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
)

func newMigrate() *cobra.Command {
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

	return Migrate(cmd.Context(), db)
}

func Migrate(ctx context.Context, db *bun.DB) error {
	return storage.Migrate(ctx, db)
}
