package driver

import (
	"context"
	"io"

	"github.com/formancehq/stack/libs/go-libs/bun/bunconnect"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

type PostgresConfig struct {
	ConnString string
}

func CLIModule(v *viper.Viper, output io.Writer, debug bool) fx.Option {

	options := make([]fx.Option, 0)
	options = append(options, fx.Provide(func(logger logging.Logger) bunconnect.ConnectionOptions {
		connectionOptions := bunconnect.ConnectionOptionsFromFlags(v, output, debug)
		logger.WithField("config", connectionOptions).Infof("Opening connection to database...")
		return connectionOptions
	}))
	options = append(options, fx.Provide(func(connectionOptions bunconnect.ConnectionOptions) (*Driver, error) {
		return New(connectionOptions), nil
	}))

	options = append(options, fx.Invoke(func(driver *Driver, lifecycle fx.Lifecycle, logger logging.Logger) error {
		lifecycle.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				logger.Infof("Initializing database...")
				return driver.Initialize(ctx)
			},
			OnStop: func(ctx context.Context) error {
				logger.Infof("Closing driver...")
				return driver.Close()
			},
		})
		return nil
	}))
	return fx.Options(options...)
}
