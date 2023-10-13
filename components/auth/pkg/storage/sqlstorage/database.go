package sqlstorage

import (
	"context"
	"time"

	auth "github.com/formancehq/auth/pkg"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	KindPostgres = "postgres"
)

type gormLogger struct {
	underlying logging.Logger
}

func (g gormLogger) LogMode(level logger.LogLevel) logger.Interface {
	return g
}

func (g gormLogger) Info(ctx context.Context, s string, i ...interface{}) {
	g.underlying.WithContext(ctx).Infof(s, i...)
}

func (g gormLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	g.underlying.WithContext(ctx).Errorf(s, i...)
}

func (g gormLogger) Error(ctx context.Context, s string, i ...interface{}) {
	g.underlying.WithContext(ctx).Errorf(s, i...)
}

func (g gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// TODO(gfyrag): Actually don't log traces
}

var _ logger.Interface = (*gormLogger)(nil)

func NewLogger(l logging.Logger) logger.Interface {
	return &gormLogger{
		underlying: l,
	}
}

func OpenPostgresDatabase(uri string) gorm.Dialector {
	return postgres.Open(uri)
}

func LoadGorm(d gorm.Dialector, gormConfig *gorm.Config) (*gorm.DB, error) {
	db, err := gorm.Open(d, gormConfig)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func MigrateTables(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).AutoMigrate(
		&auth.Client{},
		&auth.AuthRequest{},
		&auth.AccessToken{},
		&auth.AuthRequest{},
		&auth.RefreshToken{},
		&auth.User{},
		&auth.Scope{},
	)
}

var drivers = map[string]func(string) gorm.Dialector{}

func registerDriverConstructor(kind string, constructor func(string) gorm.Dialector) {
	drivers[kind] = constructor
}

func init() {
	registerDriverConstructor(KindPostgres, OpenPostgresDatabase)
}

func gormModule(kind, uri string) fx.Option {
	return fx.Options(
		fx.Provide(LoadGorm),
		fx.Provide(func() *gorm.Config {
			return &gorm.Config{}
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logging.FromContext(ctx).Info("Migrate tables")
					return MigrateTables(ctx, db)
				},
				OnStop: func(ctx context.Context) error {
					logging.FromContext(ctx).Info("Closing database...")
					defer func() {
						logging.FromContext(ctx).Info("Database closed.")
					}()
					sqlDB, err := db.DB()
					if err != nil {
						return err
					}

					return sqlDB.Close()
				},
			})
		}),
		fx.Provide(func() gorm.Dialector {
			return drivers[kind](uri)
		}),
	)
}
