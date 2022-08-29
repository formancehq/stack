package worker

import (
	"context"
	"fmt"
	"net/http"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/pkg/engine/svix"
	"github.com/numary/webhooks/pkg/httpserver"
	"github.com/numary/webhooks/pkg/storage/mongo"
	"go.uber.org/fx"
)

func StartModule(httpClient *http.Client, addr string) fx.Option {
	return fx.Module("webhooks worker",
		fx.Provide(
			func() (*http.Client, string) { return httpClient, addr },
			httpserver.NewMuxServer,
			mongo.NewConfigStore,
			svix.NewEngine,
			NewWorker,
			newWorkerHandler,
		),
		fx.Invoke(httpserver.RegisterHandler),
		fx.Invoke(httpserver.Run),
		fx.Invoke(run),
	)
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
			err1 := w.Store.Close(ctx)
			err2 := w.Reader.Close()
			if err1 != nil || err2 != nil {
				return fmt.Errorf("[closing store: %s] [closing reader: %w]", err1, err2)
			}
			return nil
		},
	})
}
