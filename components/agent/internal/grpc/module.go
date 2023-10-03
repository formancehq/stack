package grpc

import (
	"context"

	membershipgrpc "github.com/formancehq/stack/components/agent/internal/grpc/generated"
	"github.com/formancehq/stack/libs/go-libs/logging"
	"go.uber.org/fx"
	"google.golang.org/grpc"
)

func NewModule(serverAddress string, authenticator Authenticator, clientInfo ClientInfo, opts ...grpc.DialOption) fx.Option {
	return fx.Options(
		fx.Provide(func() (membershipgrpc.ServerClient, error) {
			return Connect(context.Background(), serverAddress, opts...)
		}),
		fx.Provide(func(grpcClient membershipgrpc.ServerClient, k8sClient K8SClient) *client {
			return newClient(grpcClient, k8sClient, authenticator, clientInfo)
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
