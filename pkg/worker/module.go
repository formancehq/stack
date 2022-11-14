package worker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/formancehq/webhooks/pkg/httpserver"
	"github.com/formancehq/webhooks/pkg/storage/mongo"
	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/go-libs/sharedotlp/pkg/sharedotlptraces"
	"github.com/spf13/viper"
	"go.uber.org/fx"
)

func StartModule(addr string, httpClient *http.Client, retryCron time.Duration, retrySchedule []time.Duration) fx.Option {
	var options []fx.Option

	if mod := sharedotlptraces.CLITracesModule(viper.GetViper()); mod != nil {
		options = append(options, mod)
	}

	options = append(options, fx.Provide(
		func() (string, *http.Client, time.Duration, []time.Duration) {
			return addr, httpClient, retryCron, retrySchedule
		},
		httpserver.NewMuxServer,
		mongo.NewStore,
		NewWorker,
		newWorkerHandler,
	))
	options = append(options, fx.Invoke(httpserver.RegisterHandler))
	options = append(options, fx.Invoke(httpserver.Run))
	options = append(options, fx.Invoke(run))

	return fx.Module("webhooks worker", options...)
}

func run(lc fx.Lifecycle, w *Worker) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Debugf("starting worker...")
			go func() {
				if err := w.Run(ctx); err != nil {
					sharedlogging.GetLogger(ctx).Errorf("kafka.Worker.Run: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Debugf("stopping worker...")
			w.Stop(ctx)
			w.kafkaClient.Close()
			if err := w.store.Close(ctx); err != nil {
				return fmt.Errorf("storage.Store.Close: %w", err)
			}
			return nil
		},
	})
}
