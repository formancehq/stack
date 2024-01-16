package storage

import (
	"context"

	webhooks "github.com/formancehq/webhooks/pkg"
	"github.com/pkg/errors"

	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
)

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := migrations.NewMigrator()
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			Up: func(tx bun.Tx) error {
				_, err := tx.NewCreateTable().Model((*webhooks.Config)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating 'configs' table")
				}
				_, err = tx.NewCreateIndex().Model((*webhooks.Config)(nil)).
					IfNotExists().
					Index("configs_idx").
					Column("event_types").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating index on 'configs' table")
				}
				_, err = tx.NewCreateTable().Model((*webhooks.Attempt)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating 'attempts' table")
				}
				_, err = tx.NewCreateIndex().Model((*webhooks.Attempt)(nil)).
					IfNotExists().
					Index("attempts_idx").
					Column("webhook_id", "status").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating index on 'attempts' table")
				}
				return nil
			},
		},
		migrations.Migration{
			Up: func(tx bun.Tx) error {
				_, err := tx.NewAddColumn().
					Table("configs").
					ColumnExpr("name varchar(255)").
					IfNotExists().
					Exec(ctx)
				return errors.Wrap(err, "adding 'name' column")
			},
		},
	)

	return migrator.Up(ctx, db)
}
