package system

import (
	"context"

	"github.com/formancehq/go-libs/platform/postgres"

	"github.com/formancehq/go-libs/migrations"
	"github.com/uptrace/bun"
)

func Migrate(ctx context.Context, db bun.IDB) error {
	migrator := migrations.NewMigrator(migrations.WithSchema(Schema, true))
	migrator.RegisterMigrations(
		migrations.Migration{
			Name: "Init schema",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {

				_, err := tx.NewCreateTable().
					Model((*Ledger)(nil)).
					Exec(ctx)
				if err != nil {
					return postgres.ResolveError(err)
				}

				_, err = tx.NewCreateTable().
					Model((*configuration)(nil)).
					Exec(ctx)
				return postgres.ResolveError(err)
			},
		},
		migrations.Migration{
			Name: "Add ledger, bucket naming constraints 63 chars",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists ledger varchar(63),
					add column if not exists bucket varchar(63);

					alter table ledgers
					alter column ledger type varchar(63),
					alter column bucket type varchar(63);
				`)
				return err
			},
		},
		migrations.Migration{
			Name: "Add ledger metadata",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists metadata jsonb;
				`)
				return err
			},
		},
		migrations.Migration{
			Name: "Fix empty ledger metadata",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					update ledgers
					set metadata = '{}'::jsonb
					where metadata is null;
				`)
				return err
			},
		},
		migrations.Migration{
			Name: "Add ledger state",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists state varchar(255) default 'initializing';

					update ledgers
					set state = 'in-use'
					where state = '';
				`)
				return err
			},
		},
		migrations.Migration{
			Name: "Add features column",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					alter table ledgers
					add column if not exists features jsonb;
				`)
				return err
			},
		},
		migrations.Migration{
			Name: "Add global helper",
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				_, err := tx.ExecContext(ctx, `
					CREATE EXTENSION IF NOT EXISTS "pgcrypto" SCHEMA public;

					CREATE OR REPLACE FUNCTION public.json_compact(p_json JSON, p_step INTEGER DEFAULT 0)
					RETURNS JSON
					AS $$
					DECLARE
					  v_type TEXT;
					  v_text TEXT := '';
					  v_indent INTEGER;
					  v_key TEXT;
					  v_object JSON;
					  v_count INTEGER;
					BEGIN
					  p_step := coalesce(p_step, 0);
					  -- Object or array?
					  v_type := json_typeof(p_json);
					
					  IF v_type = 'object' THEN
						-- Start object
						v_text := '{';
						SELECT count(*) - 1 INTO v_count
						FROM json_object_keys(p_json);
						-- go through keys, add them and recurse over value
						FOR v_key IN (SELECT json_object_keys(p_json))
						LOOP
						  v_text := v_text || to_json(v_key)::TEXT || ':' || public.json_compact(p_json->v_key, p_step + 1);
						  IF v_count > 0 THEN
							v_text := v_text || ',';
							v_count := v_count - 1;
						  END IF;
						  --v_text := v_text || E'\n';
						END LOOP;
						-- Close object
						v_text := v_text || '}';
					  ELSIF v_type = 'array' THEN
						-- Start array
						v_text := '[';
						v_count := json_array_length(p_json) - 1;
						-- go through elements and add them through recursion
						FOR v_object IN (SELECT json_array_elements(p_json))
						LOOP
						  v_text := v_text || public.json_compact(v_object, p_step + 1);
						  IF v_count > 0 THEN
							v_text := v_text || ',';
							v_count := v_count - 1;
						  END IF;
						  --v_text := v_text || E'\n';
						END LOOP;
						-- Close array
						v_text := v_text || ']';
					  ELSE -- A simple value
						v_text := v_text || p_json::TEXT;
					  END IF;
					  IF p_step > 0 THEN RETURN v_text;
					  ELSE RETURN v_text::JSON;
					  END IF;
					END;
					$$ LANGUAGE plpgsql;
				`)
				return err
			},
		},
	)
	return migrator.Up(ctx, db)
}
