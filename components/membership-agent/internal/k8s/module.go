package k8s

import (
	"github.com/formancehq/operator/pkg/client/v1beta3"
	"github.com/formancehq/stack/components/membership-agent/internal/grpc"
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Options(
		fx.Provide(newClient),
		fx.Provide(func(client *v1beta3.Client) grpc.K8SClient {
			return client.Stacks()
		}),
	)
}
