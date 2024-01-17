package sqlstorage

import (
	"context"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/stack/libs/go-libs/migrations"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
)

const (
	Wallets        = "wallets"
	Orchestration  = "orchestration"
	Ledger         = "ledger"
	Payments       = "payments"
	Webhooks       = "webhooks"
	Auth           = "auth"
	Reconciliation = "reconciliation"
	Search         = "search"
)

type Services []string

var AllServices = Services{
	Wallets,
	Orchestration,
	Ledger,
	Payments,
	Webhooks,
	Auth,
	Reconciliation,
	Search,
}

func Migrate(ctx context.Context, db *bun.DB) error {
	migrator := migrations.NewMigrator()
	migrator.RegisterMigrations(
		migrations.Migration{
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				script := `
CREATE TABLE access_tokens (
    id text NOT NULL,
    application_id text,
    user_id text,
    audience text,
    expiration timestamp with time zone,
    scopes text,
    refresh_token_id text
);
CREATE TABLE auth_requests (
    id text NOT NULL,
    created_at timestamp with time zone,
    application_id text,
    callback_uri text,
    transfer_state text,
    prompt text,
    ui_locales text,
    login_hint text,
    max_auth_age bigint,
    scopes text,
    response_type text,
    nonce text,
    challenge text,
    method text,
    user_id text,
    auth_time timestamp with time zone,
    code text
);
CREATE TABLE clients (
    id text NOT NULL,
    public boolean,
    redirect_uris text,
    description text,
    name text,
    post_logout_redirect_uris text,
    metadata text,
    trusted boolean,
    scopes text,
    secrets text
);
CREATE TABLE refresh_tokens (
    id text NOT NULL,
    token text,
    auth_time timestamp with time zone,
    amr text,
    audience text,
    user_id text,
    application_id text,
    expiration timestamp with time zone,
    scopes text
);
CREATE TABLE users (
    id text NOT NULL,
    subject text,
    email text
);
ALTER TABLE ONLY public.access_tokens
    ADD CONSTRAINT access_tokens_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.auth_requests
    ADD CONSTRAINT auth_requests_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.refresh_tokens
    ADD CONSTRAINT refresh_tokens_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);
ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_subject_key UNIQUE (subject);
`
				_, err := tx.Exec(script)
				return err
			},
		},
		migrations.Migration{
			UpWithContext: func(ctx context.Context, tx bun.Tx) error {
				scopes := auth.Array[string]{"openid"}
				for _, service := range AllServices {
					scopes = append(scopes, service+":read", service+":write")
				}
				_, err := tx.Exec(
					`
						UPDATE clients
						SET scopes = ?;
					`, scopes)
				return err
			},
		},
	)
	return migrator.Up(ctx, db)
}

func bunModule(connectionOptions bunconnect.ConnectionOptions) fx.Option {
	return fx.Options(
		fx.Provide(func() (*bun.DB, error) {
			return bunconnect.OpenSQLDB(connectionOptions)
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *bun.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logging.FromContext(ctx).Info("Migrate tables")

					return Migrate(ctx, db)
				},
				OnStop: func(ctx context.Context) error {
					logging.FromContext(ctx).Info("Closing database...")
					defer func() {
						logging.FromContext(ctx).Info("Database closed.")
					}()

					return db.Close()
				},
			})
		}),
	)
}
