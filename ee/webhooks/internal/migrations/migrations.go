package migrations

import (
	"context"
	"time"

	"github.com/formancehq/stack/libs/go-libs/migrations"

	"github.com/pkg/errors"
	"github.com/uptrace/bun"
)

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := migrations.NewMigrator()
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			Up: func(tx bun.Tx) error {
				_, err := tx.NewCreateTable().Model((*Config)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating 'configs' table")
				}
				_, err = tx.NewCreateIndex().Model((*Config)(nil)).
					IfNotExists().
					Index("configs_idx").
					Column("event_types").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating index on 'configs' table")
				}
				_, err = tx.NewCreateTable().Model((*Attempt)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating 'attempts' table")
				}
				_, err = tx.NewCreateIndex().Model((*Attempt)(nil)).
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
		migrations.Migration{
			Name: "Migration for V2",
			Up: func(tx bun.Tx) error {
				_, err := tx.NewCreateTable().Model((*Log)(nil)).
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating 'logs' table")
				}
				_, err = tx.NewCreateIndex().Model((*Log)(nil)).
					IfNotExists().
					Index("logs_idx").
					Column("created_at").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "creating index on 'logs' table")
				}

				_, err = tx.NewAddColumn().
					Table("configs").
					ColumnExpr("status VARCHAR(255) DEFAULT 'DISABLED'").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'status' column to Configs")
				}

				_, err = tx.NewRaw("UPDATE configs SET status = CASE WHEN active THEN 'ENABLED' ELSE 'DISABLED' END;").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "update 'status' column to Configs")
				}

				_, err = tx.NewAddColumn().
					Table("configs").
					ColumnExpr("date_status TIMESTAMP WITH TIME ZONE DEFAULT NOW()").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'date_status' column to Configs")
				}

				_, err = tx.NewRaw("UPDATE configs SET date_status = updated_at").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "update 'date_status' column to Configs")
				}

				_, err = tx.NewAddColumn().
					Table("configs").
					ColumnExpr("retry BOOLEAN DEFAULT TRUE").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'retry' column to Configs")
				}

				_, err = tx.NewAddColumn().
					Table("attempts").
					ColumnExpr("hook_name varchar DEFAULT 'Hook Name'").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'hook_name' column to Attempts")
				}

				_, err = tx.NewRaw("UPDATE attempts SET hook_name = COALESCE(config->>'name', '')").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "update 'hook_name' column to Attempts")
				}
				_, err = tx.NewAddColumn().
					Table("attempts").
					ColumnExpr("hook_endpoint varchar DEFAULT '' ").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'hook_endpoint' column to attempts")
				}

				_, err = tx.NewRaw("UPDATE attempts SET hook_endpoint = COALESCE(config->>'endpoint', '')").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "update 'hook_endpoint' column to attempts")
				}

				_, err = tx.NewAddColumn().
					Table("attempts").
					ColumnExpr("event varchar ").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'event' column to attempts")
				}

				_, err = tx.NewAddColumn().
					Table("attempts").
					ColumnExpr("date_status TIMESTAMP WITH TIME ZONE DEFAULT NOW()").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'date_status' column to attempts")
				}

				_, err = tx.NewRaw("UPDATE attempts SET date_status = updated_at").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "update 'date_status' column to Attempts")
				}

				_, err = tx.NewAddColumn().
					Table("attempts").
					ColumnExpr("date_occured TIMESTAMP WITH TIME ZONE DEFAULT NOW()").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'date_occured' column to attempts")
				}

				_, err = tx.NewRaw("UPDATE attempts SET date_occured = created_at").
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "update 'date_occured' column to Attempts")
				}

				_, err = tx.NewAddColumn().
					Table("attempts").
					ColumnExpr("comment VARCHAR").
					IfNotExists().
					Exec(ctx)
				if err != nil {
					return errors.Wrap(err, "adding 'comment' column to attempts")
				}

				return nil
			},
		},
	)

	return migrator.Up(ctx, db)
}

type Config struct {
	bun.BaseModel `bun:"table:configs"`

	ConfigUser

	ID        string    `json:"id" bun:",pk"`
	Active    bool      `json:"active"`
	Name      string    `json:"name" bun:"name,nullzero"`
	CreatedAt time.Time `json:"createdAt" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt time.Time `json:"updatedAt" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
}

type ConfigUser struct {
	Endpoint   string   `json:"endpoint"`
	Secret     string   `json:"secret"`
	EventTypes []string `json:"eventTypes" bun:"event_types,array"`
}

type Attempt struct {
	bun.BaseModel `bun:"table:attempts"`

	ID             string    `json:"id" bun:",pk"`
	WebhookID      string    `json:"webhookID" bun:"webhook_id"`
	CreatedAt      time.Time `json:"createdAt" bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt      time.Time `json:"updatedAt" bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	Config         Config    `json:"config" bun:"type:jsonb"`
	Payload        string    `json:"payload"`
	StatusCode     int       `json:"statusCode" bun:"status_code"`
	RetryAttempt   int       `json:"retryAttempt" bun:"retry_attempt"`
	Status         string    `json:"status"`
	NextRetryAfter time.Time `json:"nextRetryAfter,omitempty" bun:"next_retry_after,nullzero"`
}

type Log struct {
	bun.BaseModel `bun:"table:logs"`

	ID        string    `json:"id" bun:",pk"`
	Channel   string    `json:"channel" bun:"channel"`
	Payload   string    `json:"payload" bun:"payload"`
	CreatedAt time.Time `json:"createdAt" bun:"created_at,nullzero,notnull,default:current_timestamp"`
}
