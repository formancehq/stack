package app

import (
	"context"
	"io"

	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

const DebugFlag = "debug"

type App struct {
	options []fx.Option
	output  io.Writer
}

func (a *App) Run(ctx context.Context) error {
	ctx = defaultLoggingContext(ctx, a.output, viper.GetBool(DebugFlag))
	options := append(a.options,
		fx.NopLogger,
		fx.Provide(func() logging.Logger {
			return logging.FromContext(ctx)
		}),
	)
	app := fx.New(options...)
	err := app.Start(ctx)
	if err != nil {
		return err
	}

	select {
	case <-ctx.Done():
		return app.Stop(context.Background())
	case <-app.Done():
		return app.Err()
	}
}

func New(output io.Writer, options ...fx.Option) *App {
	return &App{
		options: options,
		output:  output,
	}
}
