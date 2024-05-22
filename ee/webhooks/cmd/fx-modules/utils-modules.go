package fxmodules

import (
	"net/http"

	"github.com/formancehq/stack/libs/go-libs/auth"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"github.com/formancehq/webhooks/cmd/flag"
	"github.com/formancehq/webhooks/internal/app/cache"
	apiutils "github.com/formancehq/webhooks/internal/app/webhook_server/api/utils"
	httpclient "github.com/formancehq/webhooks/internal/services/httpclient"
	storage "github.com/formancehq/webhooks/internal/services/storage/postgres"
	"github.com/spf13/viper"
	"github.com/uptrace/bun"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
)

var Tracer = otel.Tracer("webhook")

func ProvideHttpClient() fx.Option {

	return fx.Provide(
		func() *httpclient.DefaultHttpClient {

			client := http.Client{
				Transport: otelhttp.NewTransport(http.DefaultTransport),
			}

			defaultClient := httpclient.NewDefaultHttpClient(&client)

			return &defaultClient

		},
	)

}

func ProvideCacheParams() fx.Option {

	return fx.Provide(
		func() *cache.CacheParams {
			cacheParams := flag.LoadRunnerParams()
			return &cacheParams
		},
	)

}
func ProvideDatabase() fx.Option {

	return fx.Provide(
		func(db *bun.DB) *storage.PostgresStore {
			database := storage.NewPostgresStoreProvider(db)
			return &database
		},
	)
}

func ProvideServerParams() fx.Option {

	return fx.Provide(
		func(auth auth.Auth, logger logging.Logger, serviceInfo apiutils.ServiceInfo) *apiutils.DefaultServerParams {
			serverParams := apiutils.DefaultServerParams{}
			serverParams.Addr = viper.GetString(flag.Listen)
			serverParams.Auth = auth
			serverParams.Info = serviceInfo
			serverParams.Logger = logger

			return &serverParams
		},
	)
}

func ProvideTopics() fx.Option {
	return fx.Provide(
		func() []string {
			return viper.GetStringSlice(flag.KafkaTopics)

		},
	)
}
