package k8s

import (
	"github.com/formancehq/operator/pkg/client/v1beta3"
	"github.com/formancehq/stack/components/agent/internal/grpc"
	"go.uber.org/fx"
)

func NewModule(fakeModule bool) fx.Option {
	if fakeModule {
		return fx.Options(
			fx.Supply(fx.Annotate(NewK8SClientMock(), fx.As(new(grpc.K8SClient)))),
		)
	}

	return fx.Options(
		fx.Provide(newClient),
		fx.Provide(func(client *v1beta3.Client) grpc.K8SClient {
			return client.Stacks()
		}),
	)
}
