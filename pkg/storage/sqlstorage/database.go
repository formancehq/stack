package sqlstorage

import (
	"context"
	"io"
	"log"
	"os"
	"time"

	auth "github.com/formancehq/auth/pkg"
	"github.com/numary/go-libs/sharedlogging"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TODO: Replace by sharedlogging
func newLogger(debug bool) logger.Interface {
	out := io.Discard
	if debug {
		out = os.Stdout
	}
	return logger.New(
		log.New(out, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
}

func OpenPostgresDatabase(uri string) gorm.Dialector {
	return postgres.Open(uri)
}

func LoadGorm(d gorm.Dialector, debug bool) (*gorm.DB, error) {
	return gorm.Open(d, &gorm.Config{
		Logger: newLogger(debug),
	})
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

func gormModule(uri string, debug bool) fx.Option {
	return fx.Options(
		fx.Provide(func() gorm.Dialector {
			return OpenPostgresDatabase(uri)
		}),
		fx.Provide(func(d gorm.Dialector) (*gorm.DB, error) {
			return LoadGorm(d, debug)
		}),
		fx.Invoke(func(lc fx.Lifecycle, db *gorm.DB) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					sharedlogging.Info("Migrate tables")
					return MigrateTables(ctx, db)
				},
			})
		}),
	)
}
