package worker

import (
	"context"
	"syscall"

	"github.com/numary/go-libs/sharedlogging"
	"github.com/numary/webhooks-cloud/internal/kafka"
	"github.com/numary/webhooks-cloud/internal/storage"
	"github.com/numary/webhooks-cloud/internal/storage/mongo"
	"github.com/numary/webhooks-cloud/internal/svix"
	kafkago "github.com/segmentio/kafka-go"
	"github.com/spf13/cobra"
	svixgo "github.com/svix/svix-webhooks/go"
	"go.uber.org/fx"
)

func Start(*cobra.Command, []string) error {
	sharedlogging.Infof("env: %+v", syscall.Environ())

	app := fx.New(StartModule())
	app.Run()

	return nil
}

func StartModule() fx.Option {
	return fx.Module("webhooks-module",
		fx.Provide(
			mongo.NewConfigStore,
			svix.New,
			newKafkaEngine,
		),
		fx.Invoke(runKafkaEngine),
	)
}

func newKafkaEngine(lc fx.Lifecycle, store storage.Store, svixClient *svixgo.Svix, svixAppId string) (*kafka.Engine, error) {
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

	return kafka.NewEngine(reader, store, svixClient, svixAppId), nil
}

func runKafkaEngine(e *kafka.Engine) {
	go func() {
		if _, _, err := e.Run(context.Background()); err != nil {
			sharedlogging.Errorf("Engine.Run: %s", err)
		}
	}()
}
