package worker

import (
	"context"
	"net/http"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/kafka"
	"github.com/numary/webhooks/internal/storage"
	"github.com/numary/webhooks/internal/storage/mongo"
	"github.com/numary/webhooks/internal/svix"
	"github.com/numary/webhooks/pkg/mux"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	svixgo "github.com/svix/svix-webhooks/go"
	"go.uber.org/fx"
)

func Start(cmd *cobra.Command, args []string) {
	sharedlogging.Infof("env: %+v", syscall.Environ())

	app := fx.New(StartModule(cmd.Context()))
	app.Run()
}

func StartModule(ctx context.Context) fx.Option {
	return fx.Module("webhooks worker module",
		fx.Provide(
			mongo.NewConfigStore,
			svix.New,
			newKafkaWorker,
			newWorkerHandler,
			mux.NewWorker,
			func() context.Context { return ctx },
		),
		fx.Invoke(register),
	)
}

func newKafkaWorker(lc fx.Lifecycle, store storage.Store, svixClient *svixgo.Svix, svixAppId string) (*kafka.Worker, error) {
	cfg, err := kafka.NewKafkaReaderConfig()
	if err != nil {
		return nil, err
	}
	reader := kafkago.NewReader(cfg)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return reader.Close()
		},
	})

	return kafka.NewWorker(reader, store, svixClient, svixAppId), nil
}

func register(w *kafka.Worker, mux *http.ServeMux, h http.Handler, ctx context.Context) {
	go w.Start(ctx)
	mux.Handle("/", h)
}
