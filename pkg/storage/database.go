package storage

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/numary/auth/pkg"
	"github.com/numary/go-libs/sharedlogging"
	"go.uber.org/fx"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//TODO: Replace by sharedlogging
func newLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: false,       // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       // Disable color
		},
	)
}

func OpenDatabase(uri string) gorm.Dialector {
	return postgres.Open(uri)
}

func LoadGorm(d gorm.Dialector) (*gorm.DB, error) {
	return gorm.Open(d, &gorm.Config{
		Logger: newLogger(),
	})
}

func MigrateTables(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).AutoMigrate(
		&auth.Client{},
		&auth.Request{},
		&auth.Token{},
		&auth.Request{},
		&auth.RefreshToken{},
		&auth.User{},
	)
}

func gormModule(uri string) fx.Option {
	return fx.Options(
		fx.Provide(func() gorm.Dialector {
			return OpenDatabase(uri)
		}),
		fx.Provide(LoadGorm),
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
