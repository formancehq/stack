package sqlstorage

import (
	"database/sql"

	storage "github.com/formancehq/ledger/pkg/storage"
	ledgerstore "github.com/formancehq/ledger/pkg/storage/sqlstorage/ledger"
	"github.com/formancehq/ledger/pkg/storage/sqlstorage/schema"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"go.uber.org/fx"
)

type PostgresConfig struct {
	ConnString string
}

type ModuleConfig struct {
	PostgresConfig *PostgresConfig
}

func OpenSQLDB(dataSourceName string) (*bun.DB, error) {
	sqldb, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		return nil, err
	}

	db := bun.NewDB(sqldb, pgdialect.New())

	// TODO(polo): Should we ping here?
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func DriverModule(cfg ModuleConfig) fx.Option {
	options := make([]fx.Option, 0)

	options = append(options, fx.Provide(func() (*bun.DB, error) {
		return OpenSQLDB(cfg.PostgresConfig.ConnString)
	}))
	options = append(options, fx.Provide(func(db *bun.DB) schema.DB {
		return schema.NewPostgresDB(db)
	}))
	options = append(options, fx.Provide(func(db schema.DB) (*Driver, error) {
		return NewDriver("postgres", db), nil
	}))
	options = append(options, health.ProvideHealthCheck(func(db *sql.DB) health.NamedCheck {
		return health.NewNamedCheck("postgres", health.CheckFn(db.PingContext))
	}))

	options = append(options, fx.Provide(func(driver *Driver) storage.Driver[*ledgerstore.Store] {
		return driver
	}))
	options = append(options, fx.Invoke(func(driver storage.Driver[*ledgerstore.Store], lifecycle fx.Lifecycle) error {
		lifecycle.Append(fx.Hook{
			OnStart: driver.Initialize,
			OnStop:  driver.Close,
		})
		return nil
	}))
	options = append(options, fx.Provide(
		NewLedgerStorageDriverFromRawDriver,
		NewDefaultStorageDriverFromRawDriver,
		NewIdempotencyStorageDriverFromRawDriver,
	))
	return fx.Options(options...)
}
