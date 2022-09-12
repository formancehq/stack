package retries

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/httpserver"
	"github.com/numary/webhooks/pkg/storage/mongo"
	"go.uber.org/fx"
)

func StartModule(addr string, httpClient *http.Client, retriesCron time.Duration, retriesSchedule []time.Duration) fx.Option {
	return fx.Module("webhooks worker retries",
		fx.Provide(
			func() (string, *http.Client, time.Duration, []time.Duration) {
				return addr, httpClient, retriesCron, retriesSchedule
			},
			httpserver.NewMuxServer,
			mongo.NewStore,
			NewWorkerRetries,
			newWorkerRetriesHandler,
		),
		fx.Invoke(httpserver.RegisterHandler),
		fx.Invoke(httpserver.Run),
		fx.Invoke(run),
	)
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
