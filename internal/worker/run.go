package worker

import (
	"context"
	"fmt"
	"net/http"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks/internal/kafka"
	"github.com/numary/webhooks/internal/mux"
	"github.com/numary/webhooks/internal/storage"
	"github.com/numary/webhooks/internal/storage/mongo"
	"github.com/numary/webhooks/internal/svix"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	svixgo "github.com/svix/svix-webhooks/go"
	"go.uber.org/fx"
)

func Run(cmd *cobra.Command, args []string) error {
	app := fx.New(StartModule(cmd.Context(), http.DefaultClient))

	if err := app.Start(cmd.Context()); err != nil {
		return err
	}

	<-app.Done()

	if err := app.Stop(cmd.Context()); err != nil {
		return err
	}

	return nil
}

func StartModule(ctx context.Context, httpClient *http.Client) fx.Option {
	sharedlogging.GetLogger(ctx).Debugf(
		"starting webhooks worker module: env variables: %+v viper keys: %+v",
		syscall.Environ(), viper.AllKeys())

	return fx.Module("webhooks worker",
		fx.Provide(
			func() *http.Client { return httpClient },
			mongo.NewConfigStore,
			svix.New,
			newKafkaReaderWorker,
			newWorkerHandler,
			mux.NewWorker,
		),
		fx.Invoke(registerHandler),
		fx.Invoke(runWorker),
	)
}

func newKafkaReaderWorker(lc fx.Lifecycle, store storage.Store, svixClient *svixgo.Svix, svixAppId string) (reader *kafkago.Reader, worker *kafka.Worker, err error) {
	cfg, err := kafka.NewKafkaReaderConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("kafka.NewKafkaReaderConfig: %w", err)
	}

	reader = kafkago.NewReader(cfg)
	worker = kafka.NewWorker(reader, store, svixClient, svixAppId)
	return reader, worker, nil
}

func registerHandler(mux *http.ServeMux, h http.Handler) {
	mux.Handle("/", h)
}

func runWorker(lc fx.Lifecycle, worker *kafka.Worker, store storage.Store, reader *kafkago.Reader) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Debugf("starting worker...")
			go func() {
				if err := worker.Run(ctx); err != nil {
					sharedlogging.GetLogger(ctx).Errorf("kafka.Worker.Run: %s", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			sharedlogging.GetLogger(ctx).Debugf("stopping worker...")
			worker.Stop(ctx)
			err1 := store.Close(ctx)
			err2 := reader.Close()
			if err1 != nil || err2 != nil {
				return fmt.Errorf("[closing store: %s] [closing reader: %s]", err1, err2)
			}
			return nil
		},
	})
}
