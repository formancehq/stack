package retries

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/formancehq/go-libs/sharedlogging"
	"github.com/formancehq/go-libs/sharedotlp/pkg/sharedotlptraces"
	"github.com/formancehq/webhooks/pkg/httpserver"
	"github.com/formancehq/webhooks/pkg/storage/mongo"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func StartModule(addr string, retriesCron time.Duration, retriesSchedule []time.Duration) fx.Option {
	var options []fx.Option

	options = append(options, sharedotlptraces.CLITracesModule(viper.GetViper()))

	options = append(options, fx.Provide(
		func() (string, time.Duration, []time.Duration) {
			return addr, retriesCron, retriesSchedule
		},
		httpserver.NewMuxServer,
		mongo.NewStore,
		NewWorkerRetries,
		newWorkerRetriesHandler,
	))
	options = append(options, fx.Invoke(httpserver.RegisterHandler))
	options = append(options, fx.Invoke(httpserver.Run))
	options = append(options, fx.Invoke(run))

	sharedlogging.Debugf("starting worker retries with env:")
	for _, e := range os.Environ() {
		sharedlogging.Debugf("%s", e)
	}

	return fx.Module("webhooks worker retries", options...)
}

func run(lc fx.Lifecycle, w *WorkerRetries) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Debugf("starting worker retries with retries cron %s and schedule %+v...", w.retriesCron, w.retriesSchedule)
			go func() {
				if err := w.Run(ctx); err != nil {
					sharedlogging.GetLogger(ctx).Errorf("kafka.WorkerRetries.Run: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Debugf("stopping worker retries...")
			w.Stop(ctx)
			if err := w.store.Close(ctx); err != nil {
				return fmt.Errorf("storage.Store.Close: %w", err)
			}
			return nil
		},
	})
}
