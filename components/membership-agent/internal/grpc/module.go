package grpc

import (
	"context"
	"net/url"

	membershipgrpc "github.com/formancehq/stack/components/membership-agent/internal/grpc/generated"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func NewModule(id, serverAddress string, baseUrl *url.URL, production bool, authenticator Authenticator, opts ...grpc.DialOption) fx.Option {
	return fx.Options(
		fx.Provide(func() (membershipgrpc.ServerClient, error) {
			return Connect(context.Background(), serverAddress, opts...)
		}),
		fx.Provide(func(grpcClient membershipgrpc.ServerClient, k8sClient K8SClient) *client {
			return newClient(id, grpcClient, k8sClient, baseUrl, authenticator, production)
		}),
		fx.Invoke(func(lc fx.Lifecycle, l *client) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					ctx = logging.ContextWithLogger(context.Background(), logging.FromContext(ctx))
					if err := l.Connect(ctx); err != nil {
						return err
					}
					go func() {
						if err := l.Start(context.Background()); err != nil {
							panic(err)
						}
					}()
					return nil
				},
				OnStop: l.Stop,
			})
		}),
	)
}
