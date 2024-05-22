package fxmodules

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/formancehq/webhooks/internal/app/cache"
	webhookworker "github.com/formancehq/webhooks/internal/app/webhook_worker"
	"github.com/formancehq/webhooks/internal/services/httpclient"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"
	"go.uber.org/fx"
)

func InvokeWorker() fx.Option {

	return fx.Invoke(func(
		lc fx.Lifecycle,
		database *storage.PostgresStore,
		runnerParams *cache.CacheParams,
		client *httpclient.DefaultHttpClient,
		r *message.Router,
		subscriber message.Subscriber,
		topics []string,
	) {

		Worker := webhookworker.NewWorker(*runnerParams, database, client)
		Worker.Init()
		for _, topic := range topics {
			r.AddNoPublisherHandler(fmt.Sprintf("messages-%s", topic), topic, subscriber, Worker.HandleMessage)

		}

		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {

				return nil
			},
			OnStop: func(ctx context.Context) error {

				subscriber.Close()
				r.Close()
				Worker.Stop()
				return nil
			},
		})
	})

}
