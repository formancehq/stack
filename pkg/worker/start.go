package worker

import (
	"context"
	"fmt"
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

	app := fx.New(StartModule(cmd.Context(), http.DefaultClient))
	app.Run()
}

func StartModule(ctx context.Context, httpClient *http.Client) fx.Option {
	return fx.Module("webhooks worker module",
		fx.Provide(
			func() context.Context { return ctx },
			func() *http.Client { return httpClient },
			mongo.NewConfigStore,
			svix.New,
			newKafkaWorker,
			newWorkerHandler,
			mux.NewWorker,
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
			err1 := store.Close(ctx)
			err2 := reader.Close()
			if err1 != nil || err2 != nil {
				return fmt.Errorf("[closing store: %s] [closing reader: %s]", err1, err2)
			}
			return nil
		},
	})

	return kafka.NewWorker(reader, store, svixClient, svixAppId), nil
}

func register(w *kafka.Worker, mux *http.ServeMux, h http.Handler, ctx context.Context) (err error) {
	go func() {
		err = w.Run(ctx)
	}()
	mux.Handle("/", h)
	return
}
