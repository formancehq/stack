package client

import (
	"context"
	"time"

	"github.com/formancehq/stack/components/stargate/internal/api"
	"github.com/formancehq/stack/components/stargate/internal/client/interceptors"
	"github.com/formancehq/stack/components/stargate/internal/client/opentelemetry"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.opentelemetry.io/otel/metric"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func Module(serverURL string) fx.Option {
	options := make([]fx.Option, 0)

	options = append(options,
		fx.Provide(interceptors.NewAuthInterceptor),
		fx.Provide(func(l logging.Logger, kc keepalive.ClientParameters, authInterceptor *interceptors.AuthInterceptor) (api.StargateServiceClient, error) {
			return newGrpcClient(l, serverURL, kc, authInterceptor)
		}),
		fx.Provide(fx.Annotate(metric.NewNoopMeterProvider, fx.As(new(metric.MeterProvider)))),
		fx.Provide(opentelemetry.RegisterMetricsRegistry),
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

func NewKeepAliveClientParams(
	keepAliveClientParamTime time.Duration,
	keepAliveClientParamTimeout time.Duration,
	keepAliveClientParamPermisWithoutStream bool,
) keepalive.ClientParameters {
	return keepalive.ClientParameters{
		Time:                keepAliveClientParamTime,
		Timeout:             keepAliveClientParamTimeout,
		PermitWithoutStream: keepAliveClientParamPermisWithoutStream,
	}
}

func newGrpcClient(
	logger logging.Logger,
	serverURL string,
	kc keepalive.ClientParameters,
	authInterceptors *interceptors.AuthInterceptor,
) (api.StargateServiceClient, error) {
	conn, err := grpc.Dial(
		serverURL,
		grpc.WithStreamInterceptor(authInterceptors.StreamClientInterceptor()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(kc),
	)
	if err != nil {
		logger.Errorf("failed to connect to stargate server '%s': %s", serverURL, err)
		return nil, err
	}

	return api.NewStargateServiceClient(conn), nil
}
