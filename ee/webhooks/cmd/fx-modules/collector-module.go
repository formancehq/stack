package fxmodules

import (
	"context"

	"github.com/formancehq/webhooks/internal/app/cache"
	webhookcollector "github.com/formancehq/webhooks/internal/app/webhook_collector"
	"github.com/formancehq/webhooks/internal/services/httpclient"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"
	"go.uber.org/fx"
)

func InvokeCollector() fx.Option {

	return fx.Invoke(func(lc fx.Lifecycle,
		database *storage.PostgresStore,
		cacheParams *cache.CacheParams,
		client *httpclient.DefaultHttpClient,
	) {

		Collector := webhookcollector.NewCollector(*cacheParams, database, client)
		Collector.Init()

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				Collector.Run()
				return nil
			},
			OnStop: func(ctx context.Context) error {
				Collector.Stop()
				return nil
			},
		})

	})

}
