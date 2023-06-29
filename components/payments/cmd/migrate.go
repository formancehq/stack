package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/formancehq/payments/internal/app/storage"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"

	// Import the postgres driver.
	_ "github.com/lib/pq"

	"github.com/spf13/cobra"
)

func newMigrate() *cobra.Command {
	return &cobra.Command{
		Use:   "migrate",
		Short: "Run migrations",
		RunE:  runMigrate,
	}
}

func runMigrate(cmd *cobra.Command, args []string) error {
	postgresURI := viper.GetString(postgresURIFlag)
	if postgresURI == "" {
		postgresURI = cmd.Flag(postgresURIFlag).Value.String()
	}

	if postgresURI == "" {
		return fmt.Errorf("postgres uri is not set")
	}

	cfgEncryptionKey := viper.GetString(configEncryptionKeyFlag)
	if cfgEncryptionKey == "" {
		cfgEncryptionKey = cmd.Flag(configEncryptionKeyFlag).Value.String()
	}

	if cfgEncryptionKey != "" {
		storage.EncryptionKey = cfgEncryptionKey
	}

	// TODO: Maybe use pgx everywhere instead of pq
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return err
	}

	bunDB := bun.NewDB(db, pgdialect.New())
	if viper.GetBool(service.DebugFlag) {
		bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithWriter(cmd.OutOrStdout())))
	}

	return Migrate(cmd.Context(), bunDB)
}

func Migrate(ctx context.Context, db *bun.DB) error {
	return storage.Migrate(ctx, db)
}
