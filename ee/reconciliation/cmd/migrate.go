package cmd

import (
	"context"
	"database/sql"
	"fmt"

	storage "github.com/formancehq/reconciliation/internal/storage/migrations"
	"github.com/formancehq/stack/libs/go-libs/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/extra/bundebug"
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
	postgresURI := viper.GetString(postgresURIFlag)
	if postgresURI == "" {
		postgresURI = cmd.Flag(postgresURIFlag).Value.String()
	}

	if postgresURI == "" {
		return fmt.Errorf("postgres uri is not set")
	}

	// TODO: Maybe use pgx everywhere instead of pq
	db, err := sql.Open("postgres", postgresURI)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	bunDB := bun.NewDB(db, pgdialect.New())
	if viper.GetBool(service.DebugFlag) {
		bunDB.AddQueryHook(bundebug.NewQueryHook(bundebug.WithWriter(cmd.OutOrStdout())))
	}

	return Migrate(cmd.Context(), bunDB)
}

func Migrate(ctx context.Context, db *bun.DB) error {
	return storage.Migrate(ctx, db)
}
