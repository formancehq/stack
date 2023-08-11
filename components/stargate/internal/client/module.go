package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"

	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/client/controllers"
	"github.com/formancehq/stack/components/stargate/internal/client/interceptors"
	"github.com/formancehq/stack/components/stargate/internal/client/metrics"
	"github.com/formancehq/stack/components/stargate/internal/client/routes"
	"github.com/formancehq/stack/components/stargate/internal/server/http/middlewares"
	"github.com/formancehq/stack/libs/go-libs/health"
	"github.com/formancehq/stack/libs/go-libs/httpserver"
	"github.com/formancehq/stack/libs/go-libs/logging"
	app "github.com/formancehq/stack/libs/go-libs/service"
	"github.com/go-chi/chi/v5"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

func Module(
	bind string,
	serverURL string,
	tlsEnabled bool,
	tlsCACertificate string,
	tlsInsecureSkipVerify bool,
) fx.Option {
	options := make([]fx.Option, 0)

	options = append(options,
		fx.Provide(routes.NewRouter),
		fx.Provide(controllers.NewStargateController),
		health.Module(),
		fx.Invoke(func(lc fx.Lifecycle, h chi.Router, l logging.Logger) {
			if viper.GetBool(app.DebugFlag) {
				wrappedRouter := chi.NewRouter()
				wrappedRouter.Use(middlewares.Log())
				wrappedRouter.Mount("/", h)
				h = wrappedRouter
			}

			l.Infof("HTTP server listening on %s", bind)
			lc.Append(httpserver.NewHook(bind, h))
		}),

		fx.Provide(interceptors.NewAuthInterceptor),
		fx.Provide(func(l logging.Logger, authInterceptor *interceptors.AuthInterceptor) (api.StargateServiceClient, error) {
			return newGrpcClient(l, serverURL, tlsEnabled, tlsCACertificate, tlsInsecureSkipVerify, authInterceptor)
		}),
		fx.Provide(fx.Annotate(noop.NewMeterProvider, fx.As(new(metric.MeterProvider)))),
		fx.Provide(metrics.RegisterMetricsRegistry),
		fx.Provide(NewClient),
		fx.Invoke(func(lc fx.Lifecycle, client *Client, authInterceptor *interceptors.AuthInterceptor) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					if err := authInterceptor.ScheduleRefreshToken(); err != nil {
						return err
					}

					go func() {
						err := client.Run(context.Background())
						if err != nil && err != context.Canceled {
							panic(err)
						}
					}()

					return nil
				},
				OnStop: func(ctx context.Context) error {
					authInterceptor.Close()

					return client.Close()
				},
			})
		}),
	)

	return fx.Options(options...)
}

func newGrpcClient(
	logger logging.Logger,
	serverURL string,
	tlsEnabled bool,
	tlsCACertificate string,
	tlsInsecureSkipVerify bool,
	authInterceptors *interceptors.AuthInterceptor,
) (api.StargateServiceClient, error) {
	var credential credentials.TransportCredentials
	if !tlsEnabled {
		logger.Infof("TLS not enabled")
		credential = insecure.NewCredentials()
	} else {
		var certPool *x509.CertPool
		if tlsCACertificate != "" {
			certPool := x509.NewCertPool()
			logger.Infof("Load server certificate from config")
			if !certPool.AppendCertsFromPEM([]byte(tlsCACertificate)) {
				return nil, fmt.Errorf("failed to add server CA's certificate")
			}
		} else {
			var err error
			certPool, err = x509.SystemCertPool()
			if err != nil {
				return nil, err
			}
		}

		if tlsInsecureSkipVerify {
			logger.Infof("Disable certificate checks")
		}

		credential = credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: tlsInsecureSkipVerify,
			RootCAs:            certPool,
		})
	}

	conn, err := grpc.Dial(
		serverURL,
		grpc.WithStreamInterceptor(authInterceptors.StreamClientInterceptor()),
		grpc.WithTransportCredentials(credential),
	)
	if err != nil {
		logger.Errorf("failed to connect to stargate server '%s': %s", serverURL, err)
		return nil, err
	}

	return api.NewStargateServiceClient(conn), nil
}
